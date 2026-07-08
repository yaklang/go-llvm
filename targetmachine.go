package llvm

/*
#include <llvm-c/Core.h>
#include <llvm-c/TargetMachine.h>
#include <llvm-c/Target.h>
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"unsafe"
)

// Target wraps an LLVMTargetRef (a registered target, e.g. x86_64).
type Target struct {
	C C.LLVMTargetRef
}

// TargetMachine wraps an LLVMTargetMachineRef. It drives AOT code generation
// (assembly or object emission) for a specific target triple / cpu / features,
// replacing the external `llc` tool.
type TargetMachine struct {
	C C.LLVMTargetMachineRef
}

// TargetData wraps an LLVMTargetDataRef, the data layout of a TargetMachine.
type TargetData struct {
	C C.LLVMTargetDataRef
}

// CodeGenOptLevel selects the optimization level for code generation.
type CodeGenOptLevel C.LLVMCodeGenOptLevel

const (
	CodeGenLevelNone       CodeGenOptLevel = C.LLVMCodeGenLevelNone
	CodeGenLevelLess       CodeGenOptLevel = C.LLVMCodeGenLevelLess
	CodeGenLevelDefault    CodeGenOptLevel = C.LLVMCodeGenLevelDefault
	CodeGenLevelAggressive CodeGenOptLevel = C.LLVMCodeGenLevelAggressive
)

// RelocMode selects the relocation model.
type RelocMode C.LLVMRelocMode

const (
	RelocDefault      RelocMode = C.LLVMRelocDefault
	RelocStatic       RelocMode = C.LLVMRelocStatic
	RelocPIC          RelocMode = C.LLVMRelocPIC
	RelocDynamicNoPic RelocMode = C.LLVMRelocDynamicNoPic
	RelocROPI         RelocMode = C.LLVMRelocROPI
	RelocRWPI         RelocMode = C.LLVMRelocRWPI
	RelocROPIRWPI     RelocMode = C.LLVMRelocROPI_RWPI
)

// CodeModel selects the code model.
type CodeModel C.LLVMCodeModel

const (
	CodeModelDefault    CodeModel = C.LLVMCodeModelDefault
	CodeModelJITDefault CodeModel = C.LLVMCodeModelJITDefault
	CodeModelTiny       CodeModel = C.LLVMCodeModelTiny
	CodeModelSmall      CodeModel = C.LLVMCodeModelSmall
	CodeModelKernel     CodeModel = C.LLVMCodeModelKernel
	CodeModelMedium     CodeModel = C.LLVMCodeModelMedium
	CodeModelLarge      CodeModel = C.LLVMCodeModelLarge
)

// CodeGenFileType selects assembly vs object output.
type CodeGenFileType C.LLVMCodeGenFileType

const (
	AssemblyFile CodeGenFileType = C.LLVMAssemblyFile
	ObjectFile   CodeGenFileType = C.LLVMObjectFile
)

// GetTargetFromTriple resolves a registered Target by triple (e.g.
// "x86_64-unknown-linux-gnu"). The native target must have been initialized
// first (see InitializeNativeTarget / InitializeNativeAsmPrinter).
func GetTargetFromTriple(triple string) (Target, error) {
	ct := C.CString(triple)
	defer C.free(unsafe.Pointer(ct))
	var t C.LLVMTargetRef
	var msg *C.char
	if C.LLVMGetTargetFromTriple(ct, &t, &msg) != 0 {
		err := errors.New(C.GoString(msg))
		C.LLVMDisposeMessage(msg)
		return Target{}, err
	}
	return Target{C: t}, nil
}

// NewTargetMachine creates a TargetMachine for the given triple/cpu/features.
// For portable AOT executables use RelocStatic (or RelocPIC for PIE) and
// CodeModelSmall with CodeGenLevelDefault.
func NewTargetMachine(triple, cpu, features string, opt CodeGenOptLevel, reloc RelocMode, model CodeModel) (TargetMachine, error) {
	t, err := GetTargetFromTriple(triple)
	if err != nil {
		return TargetMachine{}, err
	}
	ct := C.CString(triple)
	cpuC := C.CString(cpu)
	featC := C.CString(features)
	defer C.free(unsafe.Pointer(ct))
	defer C.free(unsafe.Pointer(cpuC))
	defer C.free(unsafe.Pointer(featC))
	tm := C.LLVMCreateTargetMachine(t.C, ct, cpuC, featC,
		C.LLVMCodeGenOptLevel(opt), C.LLVMRelocMode(reloc), C.LLVMCodeModel(model))
	if tm == nil {
		return TargetMachine{}, errors.New("LLVMCreateTargetMachine returned nil")
	}
	return TargetMachine{C: tm}, nil
}

func (tm TargetMachine) Dispose() {
	C.LLVMDisposeTargetMachine(tm.C)
}

// Triple returns the triple this TargetMachine was created with.
func (tm TargetMachine) Triple() string {
	cs := C.LLVMGetTargetMachineTriple(tm.C)
	defer C.LLVMDisposeMessage(cs)
	return C.GoString(cs)
}

// CPU returns the cpu this TargetMachine was created with.
func (tm TargetMachine) CPU() string {
	cs := C.LLVMGetTargetMachineCPU(tm.C)
	defer C.LLVMDisposeMessage(cs)
	return C.GoString(cs)
}

// CreateTargetDataLayout returns the TargetData for this TargetMachine. The
// returned TargetData must be Dispose()d by the caller.
func (tm TargetMachine) CreateTargetDataLayout() TargetData {
	return TargetData{C: C.LLVMCreateTargetDataLayout(tm.C)}
}

func (d TargetData) Dispose() {
	C.LLVMDisposeTargetData(d.C)
}

// String returns the textual data layout representation (e.g.
// "e-m:e-p270:32:..."). Suitable for Module.SetDataLayout.
func (d TargetData) String() string {
	cs := C.LLVMCopyStringRepOfTargetData(d.C)
	defer C.LLVMDisposeMessage(cs)
	return C.GoString(cs)
}

// SetTarget sets the module's target triple.
func (m Module) SetTarget(triple string) {
	ct := C.CString(triple)
	defer C.free(unsafe.Pointer(ct))
	C.LLVMSetTarget(m.C, ct)
}

// SetDataLayout sets the module's data layout from its textual representation.
func (m Module) SetDataLayout(layout string) {
	cl := C.CString(layout)
	defer C.free(unsafe.Pointer(cl))
	C.LLVMSetDataLayout(m.C, cl)
}

// ApplyTargetMachine stamps the module with the TargetMachine's triple and
// data layout so emitted code matches the target ABI.
func (m Module) ApplyTargetMachine(tm TargetMachine) {
	m.SetTarget(tm.Triple())
	dl := tm.CreateTargetDataLayout()
	defer dl.Dispose()
	m.SetDataLayout(dl.String())
}

// EmitToFile compiles the module to an assembly or object file at filename.
// This is the in-process replacement for `llc`.
func (tm TargetMachine) EmitToFile(m Module, filename string, fileType CodeGenFileType) error {
	cf := C.CString(filename)
	defer C.free(unsafe.Pointer(cf))
	var msg *C.char
	if C.LLVMTargetMachineEmitToFile(tm.C, m.C, cf, C.LLVMCodeGenFileType(fileType), &msg) != 0 {
		err := errors.New(C.GoString(msg))
		C.LLVMDisposeMessage(msg)
		return err
	}
	return nil
}

// EmitToMemoryBuffer compiles the module in-memory and returns the resulting
// MemoryBuffer (assembly or object). The caller must Dispose the buffer.
func (tm TargetMachine) EmitToMemoryBuffer(m Module, fileType CodeGenFileType) (MemoryBuffer, error) {
	var mb C.LLVMMemoryBufferRef
	var msg *C.char
	if C.LLVMTargetMachineEmitToMemoryBuffer(tm.C, m.C, C.LLVMCodeGenFileType(fileType), &msg, &mb) != 0 {
		err := errors.New(C.GoString(msg))
		C.LLVMDisposeMessage(msg)
		return MemoryBuffer{}, err
	}
	return MemoryBuffer{C: mb}, nil
}

// DefaultTargetTriple returns the host's default target triple.
func DefaultTargetTriple() string {
	cs := C.LLVMGetDefaultTargetTriple()
	defer C.LLVMDisposeMessage(cs)
	return C.GoString(cs)
}

// NormalizeTargetTriple normalizes a target triple string.
func NormalizeTargetTriple(triple string) string {
	ct := C.CString(triple)
	defer C.free(unsafe.Pointer(ct))
	cs := C.LLVMNormalizeTargetTriple(ct)
	defer C.LLVMDisposeMessage(cs)
	return C.GoString(cs)
}

// HostCPUName returns the host cpu name (e.g. "skylake"). Use "" for
// NewTargetMachine to let LLVM pick a default.
func HostCPUName() string {
	cs := C.LLVMGetHostCPUName()
	defer C.LLVMDisposeMessage(cs)
	return C.GoString(cs)
}

// HostCPUFeatures returns the host cpu feature string (comma-separated).
func HostCPUFeatures() string {
	cs := C.LLVMGetHostCPUFeatures()
	defer C.LLVMDisposeMessage(cs)
	return C.GoString(cs)
}