package llvm

/*
#include <llvm-c/Core.h>
#include <llvm-c/ExecutionEngine.h>
#include <llvm-c/Target.h>
#include <llvm-c/Analysis.h>
#include <llvm-c/BitWriter.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

type Context struct {
	C C.LLVMContextRef
}

type Module struct {
	C C.LLVMModuleRef
}

type Builder struct {
	C C.LLVMBuilderRef
}

type Type struct {
	C C.LLVMTypeRef
}

type Value struct {
	C C.LLVMValueRef
}

type BasicBlock struct {
	C C.LLVMBasicBlockRef
}

type GenericValue struct {
	C C.LLVMGenericValueRef
}

type ExecutionEngine struct {
	C C.LLVMExecutionEngineRef
}

type PassManager struct {
	C C.LLVMPassManagerRef
}

type MemoryBuffer struct {
	C C.LLVMMemoryBufferRef
}

// Context implementation
func NewContext() Context {
	return Context{C: C.LLVMContextCreate()}
}

func (c Context) Dispose() {
	C.LLVMContextDispose(c.C)
}

func GlobalContext() Context {
	return Context{C: C.LLVMGetGlobalContext()}
}

func (c Context) NewModule(name string) Module {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Module{C: C.LLVMModuleCreateWithNameInContext(cname, c.C)}
}

// Module implementation
func NewModule(name string) Module {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Module{C: C.LLVMModuleCreateWithName(cname)}
}

func (m Module) Dispose() {
	C.LLVMDisposeModule(m.C)
}

func (m Module) Dump() {
	C.LLVMDumpModule(m.C)
}

func (m Module) String() string {
	cstr := C.LLVMPrintModuleToString(m.C)
	defer C.LLVMDisposeMessage(cstr)
	return C.GoString(cstr)
}

// Basic Types
func (c Context) Int64Type() Type {
	return Type{C: C.LLVMInt64TypeInContext(c.C)}
}

func (c Context) Int32Type() Type {
	return Type{C: C.LLVMInt32TypeInContext(c.C)}
}

func (c Context) Int8Type() Type {
	return Type{C: C.LLVMInt8TypeInContext(c.C)}
}

func (c Context) Int1Type() Type {
	return Type{C: C.LLVMInt1TypeInContext(c.C)}
}

func (c Context) VoidType() Type {
	return Type{C: C.LLVMVoidTypeInContext(c.C)}
}

func (c Context) FloatType() Type {
	return Type{C: C.LLVMFloatTypeInContext(c.C)}
}

func (c Context) DoubleType() Type {
	return Type{C: C.LLVMDoubleTypeInContext(c.C)}
}

func Int64Type() Type {
	return Type{C: C.LLVMInt64Type()}
}

func Int32Type() Type {
	return Type{C: C.LLVMInt32Type()}
}

func Int8Type() Type {
	return Type{C: C.LLVMInt8Type()}
}

func Int1Type() Type {
	return Type{C: C.LLVMInt1Type()}
}

func VoidType() Type {
	return Type{C: C.LLVMVoidType()}
}

func FloatType() Type {
	return Type{C: C.LLVMFloatType()}
}

func DoubleType() Type {
	return Type{C: C.LLVMDoubleType()}
}

func PointerType(elementType Type, addressSpace uint) Type {
	return Type{C: C.LLVMPointerType(elementType.C, C.uint(addressSpace))}
}

func FunctionType(returnType Type, paramTypes []Type, isVarArg bool) Type {
	var cParamTypes *C.LLVMTypeRef
	if len(paramTypes) > 0 {
		cParamTypes = &paramTypes[0].C
	}
	isVarArgInt := C.int(0)
	if isVarArg {
		isVarArgInt = 1
	}
	return Type{C: C.LLVMFunctionType(returnType.C, cParamTypes, C.uint(len(paramTypes)), isVarArgInt)}
}

func StructType(elementTypes []Type, packed bool) Type {
	var cElementTypes *C.LLVMTypeRef
	if len(elementTypes) > 0 {
		cElementTypes = &elementTypes[0].C
	}
	packedInt := C.int(0)
	if packed {
		packedInt = 1
	}
	return Type{C: C.LLVMStructType(cElementTypes, C.uint(len(elementTypes)), packedInt)}
}

func (t Type) IntTypeWidth() int {
	return int(C.LLVMGetIntTypeWidth(t.C))
}

// Value helpers
func ConstInt(t Type, n uint64, signExtend bool) Value {
	signExtendInt := C.int(0)
	if signExtend {
		signExtendInt = 1
	}
	return Value{C: C.LLVMConstInt(t.C, C.ulonglong(n), signExtendInt)}
}

func (v Value) Type() Type {
	return Type{C: C.LLVMTypeOf(v.C)}
}

func (v Value) SetName(name string) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.LLVMSetValueName2(v.C, cname, C.size_t(len(name)))
}

func (v Value) Name() string {
	var len C.size_t
	cname := C.LLVMGetValueName2(v.C, &len)
	return C.GoStringN(cname, C.int(len))
}

func (v Value) IsNil() bool {
	return v.C == nil
}

func (v Value) GlobalValueType() Type {
	return Type{C: C.LLVMGlobalGetValueType(v.C)}
}

// Module functions
func AddFunction(m Module, name string, ftype Type) Value {
	return m.AddFunction(name, ftype)
}

func (m Module) AddFunction(name string, ftype Type) Value {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Value{C: C.LLVMAddFunction(m.C, cname, ftype.C)}
}

func (m Module) NamedFunction(name string) Value {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Value{C: C.LLVMGetNamedFunction(m.C, cname)}
}

func (v Value) ParamsCount() int {
	return int(C.LLVMCountParams(v.C))
}

func (v Value) Params() []Value {
	count := v.ParamsCount()
	if count == 0 {
		return nil
	}
	params := make([]Value, count)
	var cParams *C.LLVMValueRef = &params[0].C
	C.LLVMGetParams(v.C, cParams)
	return params
}

func (v Value) Param(i int) Value {
	return Value{C: C.LLVMGetParam(v.C, C.uint(i))}
}

// Basic Block
func AddBasicBlock(f Value, name string) BasicBlock {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return BasicBlock{C: C.LLVMAppendBasicBlock(f.C, cname)}
}

func (c Context) AddBasicBlock(f Value, name string) BasicBlock {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return BasicBlock{C: C.LLVMAppendBasicBlockInContext(c.C, f.C, cname)}
}

func (b BasicBlock) Parent() Value {
	return Value{C: C.LLVMGetBasicBlockParent(b.C)}
}

// Builder implementation
func NewBuilder() Builder {
	return Builder{C: C.LLVMCreateBuilder()}
}

func (c Context) NewBuilder() Builder {
	return Builder{C: C.LLVMCreateBuilderInContext(c.C)}
}

func (b Builder) Dispose() {
	C.LLVMDisposeBuilder(b.C)
}

func (b Builder) SetInsertPointAtEnd(block BasicBlock) {
	C.LLVMPositionBuilderAtEnd(b.C, block.C)
}

func (b Builder) CreateRet(v Value) Value {
	return Value{C: C.LLVMBuildRet(b.C, v.C)}
}

func (b Builder) CreateRetVoid() Value {
	return Value{C: C.LLVMBuildRetVoid(b.C)}
}

func (b Builder) CreateBr(block BasicBlock) Value {
	return Value{C: C.LLVMBuildBr(b.C, block.C)}
}

func (b Builder) CreateCondBr(ifVal Value, thenBlock, elseBlock BasicBlock) Value {
	return Value{C: C.LLVMBuildCondBr(b.C, ifVal.C, thenBlock.C, elseBlock.C)}
}

func (b Builder) CreateAdd(lhs, rhs Value, name string) Value {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Value{C: C.LLVMBuildAdd(b.C, lhs.C, rhs.C, cname)}
}

func (b Builder) CreateSub(lhs, rhs Value, name string) Value {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Value{C: C.LLVMBuildSub(b.C, lhs.C, rhs.C, cname)}
}

func (b Builder) CreateMul(lhs, rhs Value, name string) Value {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Value{C: C.LLVMBuildMul(b.C, lhs.C, rhs.C, cname)}
}

func (b Builder) CreateSDiv(lhs, rhs Value, name string) Value {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Value{C: C.LLVMBuildSDiv(b.C, lhs.C, rhs.C, cname)}
}

func (b Builder) CreateSRem(lhs, rhs Value, name string) Value {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Value{C: C.LLVMBuildSRem(b.C, lhs.C, rhs.C, cname)}
}

type IntPredicate C.LLVMIntPredicate

const (
	IntEQ  IntPredicate = C.LLVMIntEQ
	IntNE  IntPredicate = C.LLVMIntNE
	IntUGT IntPredicate = C.LLVMIntUGT
	IntUGE IntPredicate = C.LLVMIntUGE
	IntULT IntPredicate = C.LLVMIntULT
	IntULE IntPredicate = C.LLVMIntULE
	IntSGT IntPredicate = C.LLVMIntSGT
	IntSGE IntPredicate = C.LLVMIntSGE
	IntSLT IntPredicate = C.LLVMIntSLT
	IntSLE IntPredicate = C.LLVMIntSLE
)

func (b Builder) CreateICmp(pred IntPredicate, lhs, rhs Value, name string) Value {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Value{C: C.LLVMBuildICmp(b.C, C.LLVMIntPredicate(pred), lhs.C, rhs.C, cname)}
}

func (b Builder) CreateCall(fn Type, callee Value, args []Value, name string) Value {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	var cArgs *C.LLVMValueRef
	if len(args) > 0 {
		cArgs = &args[0].C
	}
	return Value{C: C.LLVMBuildCall2(b.C, fn.C, callee.C, cArgs, C.uint(len(args)), cname)}
}

func (b Builder) CreatePHI(t Type, name string) Value {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Value{C: C.LLVMBuildPhi(b.C, t.C, cname)}
}

func (v Value) AddIncoming(vals []Value, blocks []BasicBlock) {
	if len(vals) != len(blocks) {
		panic("values and blocks must have same length")
	}
	if len(vals) == 0 {
		return
	}
	C.LLVMAddIncoming(v.C, &vals[0].C, &blocks[0].C, C.uint(len(vals)))
}

// Analysis
func VerifyModule(m Module, action C.LLVMVerifierFailureAction) error {
	var outMessage *C.char
	if C.LLVMVerifyModule(m.C, action, &outMessage) != 0 {
		defer C.LLVMDisposeMessage(outMessage)
		return errors.New(C.GoString(outMessage))
	}
	return nil
}

const (
	AbortProcessAction = C.LLVMAbortProcessAction
	PrintMessageAction = C.LLVMPrintMessageAction
	ReturnStatusAction = C.LLVMReturnStatusAction
)

// Target
func InitializeNativeTarget() error {
	if C.LLVMInitializeNativeTarget() != 0 {
		return errors.New("failed to initialize native target")
	}
	return nil
}

func InitializeNativeAsmPrinter() error {
	if C.LLVMInitializeNativeAsmPrinter() != 0 {
		return errors.New("failed to initialize native asm printer")
	}
	return nil
}

// ExecutionEngine
func LinkInInterpreter() {
	C.LLVMLinkInInterpreter()
}

func LinkInMCJIT() {
	C.LLVMLinkInMCJIT()
}

func NewExecutionEngine(m Module) (ExecutionEngine, error) {
	var engine C.LLVMExecutionEngineRef
	var outError *C.char
	// Force Interpreter for RunFunction compatibility in LLVM 18
	if C.LLVMCreateInterpreterForModule(&engine, m.C, &outError) != 0 {
		defer C.LLVMDisposeMessage(outError)
		return ExecutionEngine{}, errors.New(C.GoString(outError))
	}
	return ExecutionEngine{C: engine}, nil
}

func (ee ExecutionEngine) Dispose() {
	C.LLVMDisposeExecutionEngine(ee.C)
}

func (ee ExecutionEngine) RunFunction(f Value, args []GenericValue) GenericValue {
	var cArgs *C.LLVMGenericValueRef
	if len(args) > 0 {
		cArgs = &args[0].C
	}
	return GenericValue{C: C.LLVMRunFunction(ee.C, f.C, C.uint(len(args)), cArgs)}
}

func (ee ExecutionEngine) AddGlobalMapping(global Value, addr unsafe.Pointer) {
	C.LLVMAddGlobalMapping(ee.C, global.C, addr)
}

// GenericValue
func NewGenericValueFromInt(t Type, n uint64, isSigned bool) GenericValue {
	isSignedInt := C.int(0)
	if isSigned {
		isSignedInt = 1
	}
	return GenericValue{C: C.LLVMCreateGenericValueOfInt(t.C, C.ulonglong(n), isSignedInt)}
}

func (g GenericValue) Int(isSigned bool) uint64 {
	isSignedInt := C.int(0)
	if isSigned {
		isSignedInt = 1
	}
	return uint64(C.LLVMGenericValueToInt(g.C, isSignedInt))
}

func (g GenericValue) Dispose() {
	C.LLVMDisposeGenericValue(g.C)
}
