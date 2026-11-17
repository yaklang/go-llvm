//go:build !llvm_dynamic && darwin
// +build !llvm_dynamic,darwin

package llvm

/*
#cgo darwin,amd64 CFLAGS: -I${SRCDIR}/llvm/18.1.3/darwin-amd64/include
#cgo darwin,amd64 CXXFLAGS: -I${SRCDIR}/llvm/18.1.3/darwin-amd64/include
#cgo darwin,arm64 CFLAGS: -I${SRCDIR}/llvm/18.1.3/darwin-arm64/include
#cgo darwin,arm64 CFLAGS: -I${SRCDIR}/llvm/18.1.3/include
#cgo darwin,arm64 CXXFLAGS: -I${SRCDIR}/llvm/18.1.3/darwin-arm64/include
#cgo darwin,arm64 CXXFLAGS: -I${SRCDIR}/llvm/18.1.3/include
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/llvm/18.1.3/darwin-amd64/lib
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/llvm/18.1.3/darwin-arm64/lib
// Homebrew llvm@18 is keg-only. /opt/homebrew/lib and /usr/local/lib supply
// the zstd archive its static LLVM builds reference.
#cgo darwin LDFLAGS: -L/opt/homebrew/lib -L/usr/local/lib -llldMachO -llldCommon -lLLVMWindowsManifest -lLLVMXRay -lLLVMLibDriver -lLLVMDlltoolDriver -lLLVMTextAPIBinaryReader -lLLVMCoverage -lLLVMLineEditor -lLLVMAArch64Disassembler -lLLVMAArch64AsmParser -lLLVMAArch64CodeGen -lLLVMAArch64Desc -lLLVMAArch64Utils -lLLVMAArch64Info -lLLVMOrcDebugging -lLLVMOrcJIT -lLLVMWindowsDriver -lLLVMMCJIT -lLLVMJITLink -lLLVMInterpreter -lLLVMExecutionEngine -lLLVMRuntimeDyld -lLLVMOrcTargetProcess -lLLVMOrcShared -lLLVMDWP -lLLVMDebugInfoLogicalView -lLLVMDebugInfoGSYM -lLLVMOption -lLLVMObjectYAML -lLLVMObjCopy -lLLVMMCA -lLLVMMCDisassembler -lLLVMLTO -lLLVMFrontendOpenACC -lLLVMFrontendHLSL -lLLVMFrontendDriver -lLLVMExtensions -lLLVMPasses -lLLVMHipStdPar -lLLVMCoroutines -lLLVMCFGuard -lLLVMipo -lLLVMInstrumentation -lLLVMVectorize -lLLVMLinker -lLLVMFrontendOpenMP -lLLVMFrontendOffloading -lLLVMDWARFLinkerParallel -lLLVMDWARFLinkerClassic -lLLVMDWARFLinker -lLLVMGlobalISel -lLLVMMIRParser -lLLVMAsmPrinter -lLLVMSelectionDAG -lLLVMCodeGen -lLLVMTarget -lLLVMObjCARCOpts -lLLVMCodeGenTypes -lLLVMIRPrinter -lLLVMInterfaceStub -lLLVMFileCheck -lLLVMFuzzMutate -lLLVMScalarOpts -lLLVMInstCombine -lLLVMAggressiveInstCombine -lLLVMTransformUtils -lLLVMBitWriter -lLLVMAnalysis -lLLVMProfileData -lLLVMSymbolize -lLLVMDebugInfoBTF -lLLVMDebugInfoPDB -lLLVMDebugInfoMSF -lLLVMDebugInfoDWARF -lLLVMObject -lLLVMTextAPI -lLLVMMCParser -lLLVMIRReader -lLLVMAsmParser -lLLVMMC -lLLVMDebugInfoCodeView -lLLVMBitReader -lLLVMFuzzerCLI -lLLVMCore -lLLVMRemarks -lLLVMBitstreamReader -lLLVMBinaryFormat -lLLVMTargetParser -lLLVMTableGen -lLLVMSupport -lLLVMDemangle -lc++ -lz -lzstd -lcurses -lm -lxml2 -lcompression -lffi
*/
import "C"
