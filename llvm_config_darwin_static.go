//go:build !llvm_dynamic && darwin
// +build !llvm_dynamic,darwin

package llvm

/*
#cgo darwin,amd64 CFLAGS: -I${SRCDIR}/llvm/18.1.3/darwin-amd64/include
#cgo darwin,amd64 CXXFLAGS: -I${SRCDIR}/llvm/18.1.3/darwin-amd64/include
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/llvm/18.1.3/darwin-amd64/lib -llldMachO -llldCommon -lc++ -lz -lcurses -lm -lxml2 -lcompression
#cgo darwin,arm64 CFLAGS: -I${SRCDIR}/llvm/18.1.3/darwin-arm64/include
#cgo darwin,arm64 CXXFLAGS: -I${SRCDIR}/llvm/18.1.3/darwin-arm64/include
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/llvm/18.1.3/darwin-arm64/lib -llldMachO -llldCommon -lc++ -lz -lcurses -lm -lxml2 -lcompression
*/
import "C"
