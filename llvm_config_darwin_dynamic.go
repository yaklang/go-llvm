//go:build llvm_dynamic && darwin && arm64

package llvm

/*
#cgo CFLAGS: -I${SRCDIR}/llvm/18.1.3/include
#cgo CXXFLAGS: -I${SRCDIR}/llvm/18.1.3/include
#cgo LDFLAGS: -lLLVM
*/
import "C"
