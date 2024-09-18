package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/definev/glox/glox"
)

func Repl() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("> ")
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			break
		}
		glox.Interpret(text)
		fmt.Printf("\n> ")
	}
}

func RunFile(file string) {
	source := ReadFile(file)
	result := glox.Interpret(source)

	if result == glox.INTERPRET_COMPILE_ERROR {
		os.Exit(65)
		return
	}

	if result == glox.INTERPRET_RUNTIME_ERROR {
		os.Exit(70)
		return
	}
}

func ReadFile(file string) string {
	dat, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Could not open file \"%s\".\n", err)
		os.Exit(74)
		return ""
	}
	return string(dat)
}

func main() {
	glox.Compile("1000 + 231 * 33113")

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
