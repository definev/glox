package glox

import (
	"fmt"
)

type InterpretResult int

const (
	INTERPRET_OK InterpretResult = iota
	INTERPRET_COMPILE_ERROR
	INTERPRET_RUNTIME_ERROR
)

const (
	STACK_MAX = 256
)

func BinaryOp(vm *VM, op byte) {
	b := vm.Pop()
	a := vm.Pop()

	switch op {
	case OP_ADD:
		vm.Push(a + b)
	case OP_SUBTRACT:
		vm.Push(a - b)
	case OP_MULTIPLY:
		vm.Push(a * b)
	case OP_DIVIDE:
		vm.Push(a / b)
	}
}

func Interpret(line string) InterpretResult {
	return INTERPRET_OK
}

type VM struct {
	chunk    *Chunk
	ip       int
	stack    []Value
	stackTop int
}

func NewVM() *VM {
	return &VM{}
}

func (vm *VM) Init() {
	vm.stack = make([]Value, STACK_MAX)
	vm.stackTop = 0
}

func (vm *VM) Free() {}

func (vm *VM) Interpret(source string) InterpretResult {
	chunk := NewChunk()
	chunk.Init()

	if !vm.Compile(source, chunk) {
		chunk.Free()
		return INTERPRET_COMPILE_ERROR
	}

	vm.chunk = chunk
	vm.ip = 0

	result := vm.Run()

	chunk.Free()
	return result
}

func (vm *VM) Compile(source string, chunk *Chunk) bool {
	scanner.initScanner(source)

	parser.complierChunk = chunk

	parser.hadError = false
	parser.panicMode = false

	parser.advance()
	parser.expression()

	parser.consume(TOKEN_EOF, "Expect end of expression.")
	parser.endCompiler()

	return !parser.hadError
}

func (vm *VM) ReadByte() (byte, error) {
	if vm.ip >= len(*vm.chunk.Code) {
		return 0, fmt.Errorf("%b", INTERPRET_COMPILE_ERROR)
	}
	value := (*vm.chunk.Code)[vm.ip]
	vm.ip++
	return value, nil

}

func (vm *VM) Run() InterpretResult {
	for {
		if DEBUG_TRACE_EXECUTION {
			fmt.Printf("          ")
			for i := 0; i < vm.stackTop; i++ {
				slot := vm.stack[i]
				fmt.Printf("[ ")
				PrintValue(slot)
				fmt.Printf(" ]")
			}
			fmt.Printf("\n")
			vm.chunk.DisassembleInstruction(vm.ip)
		}

		instruction, err := vm.ReadByte()
		if err != nil {
			return INTERPRET_COMPILE_ERROR
		}

		switch instruction {
		case OP_RETURN:
			PrintValue(vm.Pop())
			fmt.Printf("\n")
			return INTERPRET_OK
		case OP_CONSTANT:
			constant, _ := vm.ReadByte()
			constantValue := (*vm.chunk.Constants.Values)[constant]
			vm.Push(constantValue)
		case OP_CONSTANT_LONG:
			constant := vm.chunk.ReadConstantLong(vm.ip)
			vm.ip = vm.ip + 3
			constantValue := (*vm.chunk.Constants.Values)[constant]
			vm.Push(constantValue)
		case OP_NEGATE:
			vm.Push(-vm.Pop())
		case OP_ADD:
			BinaryOp(vm, instruction)
		case OP_SUBTRACT:
			BinaryOp(vm, instruction)
		case OP_MULTIPLY:
			BinaryOp(vm, instruction)
		case OP_DIVIDE:
			BinaryOp(vm, instruction)
		}
	}
}

func (vm *VM) ResetStack() {
	vm.stackTop = 0
}

func (vm *VM) Push(value Value) {
	vm.stack[vm.stackTop] = value
	vm.stackTop++
}

func (vm *VM) Pop() Value {
	vm.stackTop--
	return vm.stack[vm.stackTop]
}

func (vm *VM) Peek(distance int) Value {
	return vm.stack[vm.stackTop-distance-1]
}
