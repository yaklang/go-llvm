package llvm

/*
#cgo CXXFLAGS: -std=c++17 -I${SRCDIR}/llvm/18.1.3/amd64/include
#include <stdlib.h>

extern int go_llvm_lld_link(int argc, char** argv,
                            char** stdout_buf, int* stdout_len,
                            char** stderr_buf, int* stderr_len);
extern void go_llvm_lld_free(char* buf);
*/
import "C"

import (
	"errors"
	"strings"
	"unsafe"
)

// LinkELF drives lld's ELF linker in-process with the given argument vector.
// args excludes the leading program name; a dummy "ld.lld" argv[0] is added.
// exitEarly=false is forced so lld returns instead of calling exit() after a
// completed link. On failure the returned error carries lld's stderr output.
//
// lld must be linked into the binary as a cgo/C++ library (see lld_link.cpp);
// it is not invoked as a subprocess. Callers must pass -nostdlib and explicit
// crt/libc/archive paths to avoid any system library search.
func LinkELF(args []string) error {
	if len(args) == 0 {
		return errors.New("lld.LinkELF: empty args")
	}
	argv := make([]*C.char, 0, len(args)+1)
	argv = append(argv, C.CString("ld.lld"))
	for _, a := range args {
		argv = append(argv, C.CString(a))
	}
	defer func() {
		for _, c := range argv {
			C.free(unsafe.Pointer(c))
		}
	}()

	var stdoutBuf, stderrBuf *C.char
	var stdoutLen, stderrLen C.int
	rc := C.go_llvm_lld_link(C.int(len(argv)), &argv[0],
		&stdoutBuf, &stdoutLen, &stderrBuf, &stderrLen)

	if stdoutBuf != nil {
		C.go_llvm_lld_free(stdoutBuf)
	}
	var stderrStr string
	if stderrBuf != nil && stderrLen > 0 {
		stderrStr = C.GoStringN(stderrBuf, stderrLen)
		C.go_llvm_lld_free(stderrBuf)
	}

	if rc != 0 {
		msg := strings.TrimSpace(stderrStr)
		if msg == "" {
			msg = "lld elf link failed (no diagnostics captured)"
		}
		return errors.New("lld: " + msg)
	}
	return nil
}