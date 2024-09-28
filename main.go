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
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Could not open file \"%s\".\n", err)
		os.Exit(74)
		return ""
	}
	return string(data)
}

func main() {
	vm := glox.NewVM()

	vm.Init()
	vm.Interpret("(5)")
	vm.Free()
}
