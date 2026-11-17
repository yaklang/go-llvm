//go:build windows && amd64 && !llvm_dynamic

package llvm

/*
#cgo CFLAGS: -I${SRCDIR}/llvm/18.1.3/include
#cgo CXXFLAGS: -I${SRCDIR}/llvm/18.1.3/include
*/
import "C"

// The platform archives and their link flags live in a per-platform module.
import _ "github.com/yaklang/go-llvm/native/windows-amd64"
