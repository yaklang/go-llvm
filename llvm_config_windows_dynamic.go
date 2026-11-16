//go:build llvm_dynamic && windows
// +build llvm_dynamic,windows

package llvm

/*
#cgo windows,amd64 CFLAGS: -I${SRCDIR}/llvm/18.1.3/windows-amd64/include
#cgo windows,amd64 CXXFLAGS: -I${SRCDIR}/llvm/18.1.3/windows-amd64/include
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/llvm/18.1.3/windows-amd64/lib -lLLVM -lws2_32 -lole32 -luuid -lkernel32 -luser32 -lgdi32 -lwinspool -lshell32 -ladvapi32 -lversion
*/
import "C"
