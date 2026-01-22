//go:build llvm_dynamic && windows
// +build llvm_dynamic,windows

package llvm

/*
#cgo windows,amd64 CFLAGS: -I${SRCDIR}/llvm/18.1.3/amd64/include
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/llvm/18.1.3/amd64/lib -lLLVM
*/
import "C"
