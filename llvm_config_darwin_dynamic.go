//go:build llvm_dynamic && darwin
// +build llvm_dynamic,darwin

package llvm

/*
#cgo darwin,amd64 CFLAGS: -I${SRCDIR}/llvm/18.1.3/amd64/include
#cgo darwin,arm64 CFLAGS: -I${SRCDIR}/llvm/18.1.3/arm64/include
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/llvm/18.1.3/amd64/lib -lLLVM
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/llvm/18.1.3/arm64/lib -lLLVM
*/
import "C"
