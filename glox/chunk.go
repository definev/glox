package glox

const (
	OP_RETURN byte = iota
	OP_CONSTANT
	OP_CONSTANT_LONG
)

type line struct {
	Value int
	Count int
}

type chunk struct {
	Count    int
	Capacity int
	Code     *[]byte
	// This is a incredibly bad idea, but I'm going to do it anyway
	Lines     *[]line
	Constants valueArray
}

func NewChunk() *chunk {
	return &chunk{
		Count:    0,
		Capacity: 0,
		Code:     nil,
	}
}

func (c *chunk) Init() {
	c.Count = 0
	c.Capacity = 0
	c.Code = nil
	c.Lines = nil
	c.Constants.Init()
}

func (c *chunk) Write(Byte byte, Line int) {
	if c.Capacity < c.Count+1 {
		c.Capacity = GROW_CAPACITY(c.Capacity)
		c.Code = GROW_ARRAY(c.Code, c.Capacity)
		c.Lines = GROW_ARRAY(c.Lines, c.Capacity)
	}

	(*c.Code)[c.Count] = Byte
	c.WriteLine(Line)

	c.Count += 1
}

func (c *chunk) Free() {
	c.Init()
	c.Constants.Free()
}
