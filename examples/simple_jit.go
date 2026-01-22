package main

import (
	"fmt"

	"github.com/yaklang/go-llvm"
)

func main() {
	llvm.LinkInInterpreter()
	llvm.InitializeNativeTarget()
	llvm.InitializeNativeAsmPrinter()

	ctx := llvm.NewContext()
	defer ctx.Dispose()

	mod := llvm.NewModule("test_module")
	// Note: ExecutionEngine takes ownership of module, so we don't dispose it if engine creation succeeds
	// defer mod.Dispose()

	builder := ctx.NewBuilder()
	defer builder.Dispose()

	// Define function: int sum(int a, int b)
	int64Type := llvm.Int64Type()
	fnType := llvm.FunctionType(int64Type, []llvm.Type{int64Type, int64Type}, false)
	fn := mod.AddFunction("sum", fnType)

	// Basic block
	block := ctx.AddBasicBlock(fn, "entry")
	builder.SetInsertPointAtEnd(block)

	// params
	params := fn.Params()
	a := params[0]
	b := params[1]

	// add
	sum := builder.CreateAdd(a, b, "sum_res")
	builder.CreateRet(sum)

	// Verify
	if err := llvm.VerifyModule(mod, llvm.PrintMessageAction); err != nil {
		panic(err)
	}

	// JIT
	engine, err := llvm.NewExecutionEngine(mod)
	if err != nil {
		panic(err)
	}
	defer engine.Dispose()

	// Verify NamedFunction works
	fnFound := mod.NamedFunction("sum")
	if fnFound.IsNil() {
		panic("NamedFunction('sum') returned nil!")
	}
	fmt.Printf("NamedFunction found: %s\n", fnFound.Name())

	// Run
	args := []llvm.GenericValue{
		llvm.NewGenericValueFromInt(int64Type, 10, true),
		llvm.NewGenericValueFromInt(int64Type, 32, true),
	}
	res := engine.RunFunction(fn, args)
	fmt.Printf("10 + 32 = %d\n", res.Int(true))
}
