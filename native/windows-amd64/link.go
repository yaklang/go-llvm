//go:build windows && amd64

// Package native carries the vendored LLVM 18 static archives for windows-amd64.
//
// The archives live next to this file so that the whole platform tree ships as
// its own Go module, well below the module proxy size limit. The main module
// keeps the headers and blank-imports this package, so callers keep importing
// github.com/yaklang/go-llvm only.
package native

/*
#cgo LDFLAGS: -L${SRCDIR}/lib -static -Wl,--start-group -llldCOFF -llldCommon -llldMinGW -lLLVMAggressiveInstCombine -lLLVMAnalysis -lLLVMAsmParser -lLLVMAsmPrinter -lLLVMBinaryFormat -lLLVMBitReader -lLLVMBitWriter -lLLVMBitstreamReader -lLLVMCFGuard -lLLVMCFIVerify -lLLVMCodeGen -lLLVMCodeGenTypes -lLLVMCore -lLLVMCoroutines -lLLVMCoverage -lLLVMDWARFLinker -lLLVMDWARFLinkerClassic -lLLVMDWARFLinkerParallel -lLLVMDWP -lLLVMDebugInfoBTF -lLLVMDebugInfoCodeView -lLLVMDebugInfoDWARF -lLLVMDebugInfoGSYM -lLLVMDebugInfoLogicalView -lLLVMDebugInfoMSF -lLLVMDebugInfoPDB -lLLVMDebuginfod -lLLVMDemangle -lLLVMDiff -lLLVMDlltoolDriver -lLLVMExecutionEngine -lLLVMExegesis -lLLVMExegesisX86 -lLLVMExtensions -lLLVMFileCheck -lLLVMFrontendDriver -lLLVMFrontendHLSL -lLLVMFrontendOffloading -lLLVMFrontendOpenACC -lLLVMFrontendOpenMP -lLLVMFuzzMutate -lLLVMFuzzerCLI -lLLVMGlobalISel -lLLVMHipStdPar -lLLVMIRPrinter -lLLVMIRReader -lLLVMInstCombine -lLLVMInstrumentation -lLLVMInterfaceStub -lLLVMInterpreter -lLLVMJITLink -lLLVMLTO -lLLVMLibDriver -lLLVMLineEditor -lLLVMLinker -lLLVMMC -lLLVMMCA -lLLVMMCDisassembler -lLLVMMCJIT -lLLVMMCParser -lLLVMMIRParser -lLLVMObjCARCOpts -lLLVMObjCopy -lLLVMObject -lLLVMObjectYAML -lLLVMOption -lLLVMOrcDebugging -lLLVMOrcJIT -lLLVMOrcShared -lLLVMOrcTargetProcess -lLLVMPasses -lLLVMProfileData -lLLVMRemarks -lLLVMRuntimeDyld -lLLVMScalarOpts -lLLVMSelectionDAG -lLLVMSupport -lLLVMSymbolize -lLLVMTableGen -lLLVMTableGenCommon -lLLVMTableGenGlobalISel -lLLVMTarget -lLLVMTargetParser -lLLVMTextAPI -lLLVMTextAPIBinaryReader -lLLVMTransformUtils -lLLVMVectorize -lLLVMWindowsDriver -lLLVMWindowsManifest -lLLVMX86AsmParser -lLLVMX86CodeGen -lLLVMX86Desc -lLLVMX86Disassembler -lLLVMX86Info -lLLVMX86TargetMCA -lLLVMXRay -lLLVMipo -Wl,--end-group -lstdc++ -lwinpthread -lpsapi -lshell32 -lole32 -luuid -ladvapi32 -lws2_32 -loleaut32 -luxtheme -lversion
*/
import "C"
