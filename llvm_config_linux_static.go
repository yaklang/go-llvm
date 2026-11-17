//go:build linux && amd64 && !llvm_dynamic

package llvm

/*
#cgo CFLAGS: -I${SRCDIR}/llvm/18.1.3/linux-amd64/include
#cgo CXXFLAGS: -I${SRCDIR}/llvm/18.1.3/linux-amd64/include
*/
import "C"

// The platform archives and their link flags live in a per-platform module.
import _ "github.com/yaklang/go-llvm/native/linux-amd64"
