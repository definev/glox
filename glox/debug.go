package glox

import "fmt"

const DEBUG_TRACE_EXECUTION = true
const DEBUG_PRINT_CODE = false

func (c *Chunk) DisassembleChunk(name string) {
	fmt.Printf("== %s ==\n", name)

	for offset := 0; offset < c.Count; {
		offset = c.DisassembleInstruction(offset)
	}
}

func (c *Chunk) DisassembleInstruction(offset int) int {
	fmt.Printf("%04d ", offset)
	if offset > 0 &&
		c.GetLine(offset) == c.GetLine(offset-1) {
		fmt.Printf("   | ")
	} else {
		fmt.Printf("%4d ", c.GetLine(offset))
	}

	instruction := (*c.Code)[offset]

	switch instruction {
	case OP_RETURN:
		return simpleInstruction("OP_RETURN", offset)
	case OP_CONSTANT:
		return constantInstruction("OP_CONSTANT", c, offset)
	case OP_CONSTANT_LONG:
		return constantLongInstruction("OP_CONSTANT_LONG", c, offset)
	case OP_NEGATE:
		return simpleInstruction("OP_NEGATE", offset)
	case OP_ADD:
		return simpleInstruction("OP_ADD", offset)
	case OP_SUBTRACT:
		return simpleInstruction("OP_SUBTRACT", offset)
	case OP_MULTIPLY:
		return simpleInstruction("OP_MULTIPLY", offset)
	case OP_DIVIDE:
		return simpleInstruction("OP_DIVIDE", offset)
	default:
		fmt.Println("Unknown opcode")
		return offset + 1
	}
}

func simpleInstruction(name string, offset int) int {
	fmt.Println(name)
	return offset + 1
}

func constantInstruction(name string, c *Chunk, offset int) int {
	constant := (*c.Code)[offset+1]
	fmt.Printf("%-16s %4d '", name, constant)
	c.Constants.Print(int(constant))
	fmt.Printf("'\n")
	return offset + 2
}

func constantLongInstruction(name string, c *Chunk, offset int) int {
	constant := c.ReadConstantLong(offset + 1)

	fmt.Printf("%-16s %4d '", name, constant)
	c.Constants.Print(int(constant))
	fmt.Printf("'\n")
	return offset + 4
}
