//go:build llvm_dynamic
// +build llvm_dynamic

package llvm

/*
#cgo linux,amd64 CFLAGS: -I${SRCDIR}/llvm/18.1.3/linux-amd64/include
#cgo linux,amd64 CXXFLAGS: -I${SRCDIR}/llvm/18.1.3/linux-amd64/include
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/llvm/18.1.3/linux-amd64/lib -lLLVM
*/
import "C"
