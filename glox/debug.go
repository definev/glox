package glox

import "fmt"

func (c *chunk) DisassembleChunk(name string) {
	fmt.Printf("== %s ==\n", name)

	for offset := 0; offset < c.Count; {
		offset = c.DisassembleInstruction(offset)
	}
}

func (c *chunk) DisassembleInstruction(offset int) int {
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
	default:
		fmt.Println("Unknown opcode")
		return offset + 1
	}
}

func simpleInstruction(name string, offset int) int {
	fmt.Println(name)
	return offset + 1
}

func constantInstruction(name string, c *chunk, offset int) int {
	constant := (*c.Code)[offset+1]
	fmt.Printf("%-16s %4d '", name, constant)
	c.Constants.Print(int(constant))
	fmt.Printf("'\n")
	return offset + 2
}

func constantLongInstruction(name string, c *chunk, offset int) int {
	constant0 := (*c.Code)[offset+1]
	constant1 := (*c.Code)[offset+2]
	constant2 := (*c.Code)[offset+3]

	constant := int(constant0)<<16 | int(constant1)<<8 | int(constant2)

	fmt.Printf("%-16s %4d '", name, constant)
	c.Constants.Print(int(constant))
	fmt.Printf("'\n")
	return offset + 4
}
