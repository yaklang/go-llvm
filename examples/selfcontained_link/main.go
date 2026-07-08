// selfcontained_link demonstrates the in-process AOT pipeline: build a tiny
// LLVM module, emit it to a relocatable object via go-llvm TargetMachine
// (replacing llc), then link object + crt + static libc/libgcc/libgcc_eh into a
// fully-static executable via go-llvm's in-process lld (replacing clang). No
// llc/clang/ld subprocess is invoked. The output executable is run and its exit
// code is checked.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/yaklang/go-llvm"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// gccFile resolves a build-time gcc file (crtbegin.o, libgcc.a, ...) via
// `gcc -print-file-name=...`. This is a go-llvm unit check, so using the
// build-time gcc is fine; ssa2llvm embeds these artifacts instead and never
// invokes gcc at runtime.
func gccFile(name string) (string, error) {
	out, err := exec.Command("gcc", "-print-file-name="+name).Output()
	if err != nil {
		return "", fmt.Errorf("locate %s: %w", name, err)
	}
	p := strings.TrimSpace(string(out))
	if p == "" || p == name {
		return "", fmt.Errorf("could not resolve %s via gcc", name)
	}
	return p, nil
}

func run() error {
	if err := llvm.InitializeNativeTarget(); err != nil {
		return fmt.Errorf("init native target: %w", err)
	}
	if err := llvm.InitializeNativeAsmPrinter(); err != nil {
		return fmt.Errorf("init asm printer: %w", err)
	}

	ctx := llvm.NewContext()
	defer ctx.Dispose()
	mod := ctx.NewModule("hello")
	builder := ctx.NewBuilder()
	defer builder.Dispose()

	// int main() { return 42; }
	i32 := ctx.Int32Type()
	mainType := llvm.FunctionType(i32, nil, false)
	mainFn := mod.AddFunction("main", mainType)
	entry := ctx.AddBasicBlock(mainFn, "entry")
	builder.SetInsertPointAtEnd(entry)
	builder.CreateRet(llvm.ConstInt(i32, 42, false))

	if err := llvm.VerifyModule(mod, llvm.PrintMessageAction); err != nil {
		return fmt.Errorf("verify module: %w", err)
	}

	workDir, err := os.MkdirTemp("", "sc_link_")
	if err != nil {
		return err
	}
	defer os.RemoveAll(workDir)

	triple := llvm.DefaultTargetTriple()
	tm, err := llvm.NewTargetMachine(triple, "", "", llvm.CodeGenLevelDefault, llvm.RelocStatic, llvm.CodeModelSmall)
	if err != nil {
		return fmt.Errorf("create target machine: %w", err)
	}
	defer tm.Dispose()
	mod.ApplyTargetMachine(tm)

	objPath := filepath.Join(workDir, "hello.o")
	if err := tm.EmitToFile(mod, objPath, llvm.ObjectFile); err != nil {
		return fmt.Errorf("emit object: %w", err)
	}
	fmt.Println("emitted object:", objPath)

	// Locate the static-link crt + system libs via the build-time gcc.
	crtBegin, err := gccFile("crtbegin.o")
	if err != nil {
		return err
	}
	crtEnd, err := gccFile("crtend.o")
	if err != nil {
		return err
	}
	libgcc, err := gccFile("libgcc.a")
	if err != nil {
		return err
	}
	libgccEh, err := gccFile("libgcc_eh.a")
	if err != nil {
		return err
	}
	crtDir := "/usr/lib/x86_64-linux-gnu"

	outPath := filepath.Join(workDir, "hello")
	if err := llvm.LinkExecutableStatic(llvm.StaticLinkInput{
		ObjectPath: objPath,
		CRTBegin:   []string{filepath.Join(crtDir, "crt1.o"), filepath.Join(crtDir, "crti.o"), crtBegin},
		CRTEnd:     []string{crtEnd, filepath.Join(crtDir, "crtn.o")},
		SystemLibs: []string{filepath.Join(crtDir, "libc.a"), libgcc, libgccEh},
		OutputPath: outPath,
		ExtraArgs:  []string{"-e", "_start"},
	}); err != nil {
		return fmt.Errorf("in-process lld link: %w", err)
	}
	fmt.Println("linked static exe:", outPath)

	// Verify it runs with exit code 42.
	cmd := exec.Command(outPath)
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code := exitErr.ExitCode()
			fmt.Printf("output exe exited with code %d\n", code)
			if code == 42 {
				fmt.Println("PASS: in-process compile (TargetMachine) + link (lld) works end-to-end")
				return nil
			}
			return fmt.Errorf("unexpected exit code %d (want 42)", code)
		}
		return fmt.Errorf("run output exe: %w", err)
	}
	return fmt.Errorf("expected exit code 42, got 0")
}