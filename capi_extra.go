package llvm

/*
#include <llvm-c/Core.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// This file wraps LLVM-C entry points that callers used to reach through
// cgo directly. Two rules apply to every helper here:
//
//  1. No cgo type ever escapes the package.
//  2. Opcode checks that changed across C API versions are structural, not
//     enum-based. LLVM 20 removed LLVMBr from llvm-c/Core.h and reports
//     unconditional branches through LLVMUncondBr (= 70), while older
//     releases still return LLVMBr (= 2). Headers and libraries can also be
//     different releases in one build, so numeric or symbolic comparisons
//     are both unsafe.

// IsNil reports whether m is a zero Module handle.
func (m Module) IsNil() bool {
	return m.C == nil
}

// IsNil reports whether t is a zero Type handle.
func (t Type) IsNil() bool {
	return t.C == nil
}

// IsAInstruction reports whether v is an instruction.
func (v Value) IsAInstruction() bool {
	return v.C != nil && C.LLVMIsAInstruction(v.C) != nil
}

// IsATerminatorInst reports whether v is a terminator instruction.
func (v Value) IsATerminatorInst() bool {
	return v.C != nil && C.LLVMIsATerminatorInst(v.C) != nil
}

// InstructionParent returns the basic block that owns v. Values that are not
// instructions return a nil BasicBlock.
func (v Value) InstructionParent() BasicBlock {
	if !v.IsAInstruction() {
		return BasicBlock{}
	}
	return BasicBlock{C: C.LLVMGetInstructionParent(v.C)}
}

// IsCallInst reports whether v is a call instruction.
func (v Value) IsCallInst() bool {
	if !v.IsAInstruction() {
		return false
	}
	return C.LLVMGetInstructionOpcode(v.C) == C.LLVMCall
}

// IsUnconditionalBr reports whether v is an unconditional branch ("br label
// %bb"). It is the only terminator with exactly one operand, and that operand
// is a basic block, so the check works for both the pre-LLVM-20 LLVMBr and the
// later LLVMUncondBr C API opcodes.
func (v Value) IsUnconditionalBr() bool {
	if !v.IsATerminatorInst() {
		return false
	}
	if v.NumOperands() != 1 {
		return false
	}
	operand := v.Operand(0)
	return operand.C != nil && C.LLVMIsABasicBlock(operand.C) != nil
}

// CalledValue returns the callee operand of a call instruction.
func (v Value) CalledValue() Value {
	if v.C == nil {
		return Value{}
	}
	return Value{C: C.LLVMGetCalledValue(v.C)}
}

// HasUses reports whether any value references v.
func (v Value) HasUses() bool {
	return v.C != nil && C.LLVMGetFirstUse(v.C) != nil
}

// DeleteFunction deletes the function v from its module.
func (v Value) DeleteFunction() {
	if v.C == nil {
		return
	}
	C.LLVMDeleteFunction(v.C)
}

// IsNullPointer reports whether v is a null pointer constant.
func (v Value) IsNullPointer() bool {
	return v.C != nil && C.LLVMIsNull(v.C) != 0
}

// SameAs reports whether v and other are the same LLVM value.
func (v Value) SameAs(other Value) bool {
	return v.C == other.C
}

// AsBasicBlock views v as a basic block. It returns a nil BasicBlock when v is
// not one.
func (v Value) AsBasicBlock() BasicBlock {
	if v.C == nil {
		return BasicBlock{}
	}
	return BasicBlock{C: C.LLVMValueAsBasicBlock(v.C)}
}

// Terminator returns the terminator instruction of b, or a nil Value when b
// has no terminator.
func (b BasicBlock) Terminator() Value {
	if b.C == nil {
		return Value{}
	}
	return Value{C: C.LLVMGetBasicBlockTerminator(b.C)}
}

// AsValue views b as a value so it can be compared with instruction operands.
func (b BasicBlock) AsValue() Value {
	if b.C == nil {
		return Value{}
	}
	return Value{C: C.LLVMBasicBlockAsValue(b.C)}
}

// SameAs reports whether b and other are the same basic block.
func (b BasicBlock) SameAs(other BasicBlock) bool {
	return b.C == other.C
}

// InsertBlock returns the block that the builder currently inserts into.
func (b Builder) InsertBlock() BasicBlock {
	if b.C == nil {
		return BasicBlock{}
	}
	return BasicBlock{C: C.LLVMGetInsertBlock(b.C)}
}

// CreateAlloca builds a stack allocation of one element of t.
func (b Builder) CreateAlloca(t Type, name string) Value {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Value{C: C.LLVMBuildAlloca(b.C, t.C, cname)}
}

// CreateArrayAlloca builds a stack allocation of count elements of elemType.
func (b Builder) CreateArrayAlloca(elemType Type, count Value, name string) Value {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return Value{C: C.LLVMBuildArrayAlloca(b.C, elemType.C, count.C, cname)}
}

// CreateUnreachable emits an unreachable terminator at the insert point.
func (b Builder) CreateUnreachable() Value {
	return Value{C: C.LLVMBuildUnreachable(b.C)}
}
