//go:build llvm_dynamic && windows && amd64

package llvm

/*
#cgo CFLAGS: -I${SRCDIR}/llvm/18.1.3/include
#cgo CXXFLAGS: -I${SRCDIR}/llvm/18.1.3/include
#cgo LDFLAGS: -lLLVM -lws2_32 -lole32 -luuid -lkernel32 -luser32 -lgdi32 -lwinspool -lshell32 -ladvapi32 -lversion
*/
import "C"
