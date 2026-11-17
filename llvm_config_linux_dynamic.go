//go:build llvm_dynamic && linux && amd64

package llvm

/*
#cgo CFLAGS: -I${SRCDIR}/llvm/18.1.3/linux-amd64/include
#cgo CXXFLAGS: -I${SRCDIR}/llvm/18.1.3/linux-amd64/include
#cgo LDFLAGS: -lLLVM
*/
import "C"
