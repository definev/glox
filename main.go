package main

import "github.com/definev/glox/glox"

func main() {
	vm := glox.NewVM()

	vm.Init()

	chunk := glox.NewChunk()
	chunk.Init()

	chunk.WriteConstant(glox.Value(10), 100)
	chunk.WriteConstant(glox.Value(5), 100)

	chunk.Write(glox.OP_DIVIDE, 100)
	chunk.Write(glox.OP_NEGATE, 100)

	chunk.Write(glox.OP_RETURN, 100)
	// chunk.DisassembleChunk("test_chunk")
	vm.Interpret(chunk)
	vm.Free()
	chunk.Free()
}
