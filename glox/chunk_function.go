package glox

func (c *chunk) AddConstant(value Value) int {
	c.Constants.Write(value)
	return c.Constants.Count - 1
}

func (c *chunk) WriteConstant(value Value, Line int) {
	c.Write(OP_CONSTANT_LONG, Line)
	constant := c.AddConstant(value)

	constantBytes := make([]byte, 3)
	constantBytes[0] = byte(constant >> 16)
	constantBytes[1] = byte(constant >> 8)
	constantBytes[2] = byte(constant)

	c.Write(constantBytes[0], Line)
	c.Write(constantBytes[1], Line)
	c.Write(constantBytes[2], Line)
}

func (c *chunk) ReadConstantLong(offset int) int {
	constant0 := (*c.Code)[offset]
	constant1 := (*c.Code)[offset+1]
	constant2 := (*c.Code)[offset+2]

	constant := int(constant0)<<16 | int(constant1)<<8 | int(constant2)
	return constant
}

func (c *chunk) GetLine(offset int) int {
	line := 0
	for i := 0; i < len(*c.Lines); i++ {
		line += (*c.Lines)[i].Count
		if line > offset {
			return (*c.Lines)[i].Value
		}
	}
	return -1
}

func (c *chunk) WriteLine(Line int) {
	linesLen := len(*c.Lines)
	if linesLen == 0 {
		(*c.Lines)[c.Count] = line{
			Value: Line,
			Count: 1,
		}
	} else if (*c.Lines)[c.Count].Value == Line {
		(*c.Lines)[c.Count].Count += 1
	} else {
		(*c.Lines)[c.Count] = line{
			Value: Line,
			Count: 1,
		}
	}
}
