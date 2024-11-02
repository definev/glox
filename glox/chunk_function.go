package glox

func SplitConstant(value int) (byte, byte, byte) {
	constant0 := byte(value >> 16)
	constant1 := byte(value >> 8)
	constant2 := byte(value)

	return constant0, constant1, constant2
}

func (c *Chunk) AddConstant(value Value) int {
	c.Constants.Write(value)
	return c.Constants.Count - 1
}

func (c *Chunk) WriteConstant(value Value, Line int) {
	c.Write(OP_CONSTANT_LONG, Line)
	constant := c.AddConstant(value)
	constant0, constant1, constant2 := SplitConstant(constant)
	c.Write(byte(constant0), Line)
	c.Write(byte(constant1), Line)
	c.Write(byte(constant2), Line)
}

func (c *Chunk) ReadConstant(offset int) int {
	constant0 := (*c.Code)[offset]
	constant1 := (*c.Code)[offset+1]
	constant2 := (*c.Code)[offset+2]

	constant := int(constant0)<<16 | int(constant1)<<8 | int(constant2)
	return constant
}

func (c *Chunk) GetLine(offset int) int {
	line := 0
	for i := 0; i < len(*c.Lines); i++ {
		line += (*c.Lines)[i].Count
		if line > offset {
			return (*c.Lines)[i].Value
		}
	}
	return -1
}

func (c *Chunk) WriteLine(Line int) {
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
