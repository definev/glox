package main

import "github.com/definev/glox/glox"

func main() {
	chunk := glox.NewChunk()
	chunk.Init()

	for i := 0; i < 1000; i++ {
		chunk.WriteConstant(glox.Value(i), i/10)
	}

	chunk.Write(glox.OP_RETURN, 100)
	chunk.DisassembleChunk("test_chunk")
	chunk.Free()
}
