package llvm

// StaticLinkInput describes a fully self-contained static ELF link performed
// in-process by lld (LinkExecutableStatic). No system library search is done:
// -nostdlib -static plus explicit crt/libc/archive paths are used, so the link
// has zero dependence on the host LLVM/toolchain environment.
type StaticLinkInput struct {
	// ObjectPath is the relocatable object to link (e.g. produced by
	// TargetMachine.EmitToFile with ObjectFile).
	ObjectPath string
	// Archives are static runtime archives in link order (e.g. libyak.a,
	// libgc.a, obfuscation runtime archives). They are linked inside a
	// --start-group/--end-group together with SystemLibs so mutual references
	// (runtime -> libc/libgcc, libgc -> libc) resolve regardless of order.
	Archives []string
	// SystemLibs are the static system libraries (libc.a, libgcc.a,
	// libgcc_eh.a) linked inside the group. glibc's libc.a references the
	// libgcc/libgcc_eh unwinding helpers (_Unwind_Resume, __gcc_personality_v0,
	// soft-float), and they may reference libc back, so they must be grouped.
	SystemLibs []string
	// CRTBegin are crt objects placed at the start of the link, in order:
	// crt1.o (provides _start), crti.o, crtbegin.o (gcc EH frame init).
	CRTBegin []string
	// CRTEnd are crt objects placed at the end of the link, in order:
	// crtend.o (gcc EH frame fini), crtn.o.
	CRTEnd []string
	// OutputPath is the resulting executable path.
	OutputPath string
	// ExtraArgs are extra lld arguments (e.g. "--gc-sections", "-e", "_start").
	ExtraArgs []string
}

// LinkExecutableStatic links a portable, fully-static ELF executable in-process
// via lld. The argument vector follows the canonical static-link order:
//
//	crt1.o crti.o crtbegin.o  <object>  --start-group <archives> <system libs> --end-group  crtend.o crtn.o  [extra]  -o out
//
// with -nostdlib -static --no-undefined prepended. Archives and system libs are
// grouped so their mutual references resolve. Because every input is an
// explicit path, lld performs no system library search.
func LinkExecutableStatic(in StaticLinkInput) error {
	return LinkELF(StaticLinkArgs(in))
}

// StaticLinkArgs builds the lld argument vector for a static executable link
// without invoking lld. Exposed so callers can trace the equivalent command
// line when driving lld in-process (no subprocess is spawned).
func StaticLinkArgs(in StaticLinkInput) []string {
	args := make([]string, 0, 12+len(in.CRTBegin)+len(in.Archives)+len(in.SystemLibs)+len(in.CRTEnd)+len(in.ExtraArgs))
	args = append(args, "-nostdlib", "-static", "--no-undefined")
	args = append(args, in.CRTBegin...)
	args = append(args, in.ObjectPath)
	args = append(args, "--start-group")
	args = append(args, in.Archives...)
	args = append(args, in.SystemLibs...)
	args = append(args, "--end-group")
	args = append(args, in.CRTEnd...)
	args = append(args, in.ExtraArgs...)
	args = append(args, "-o", in.OutputPath)
	return args
}