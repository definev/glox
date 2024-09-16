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

type VM struct {
	chunk    *chunk
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

func (vm *VM) Interpret(chunk *chunk) InterpretResult {
	vm.chunk = chunk
	vm.ip = 0
	return vm.Run()
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
		case OP_SUBTRACT:
		case OP_MULTIPLY:
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
