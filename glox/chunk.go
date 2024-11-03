package glox

const (
	OP_RETURN byte = iota
	OP_CONSTANT_LONG
	OP_NIL
	OP_TRUE
	OP_FALSE
	OP_EQUAL
	OP_GREATER
	OP_LESS
	OP_ADD
	OP_SUBTRACT
	OP_MULTIPLY
	OP_DIVIDE
	OP_NEGATE
	OP_NOT
	OP_PRINT
	OP_POP
	OP_DEFINE_GLOBAL
	OP_GET_GLOBAL
)

type line struct {
	Value int
	Count int
}

type Chunk struct {
	Count    int
	Capacity int
	Code     *[]byte
	// This is a incredibly bad idea, but I'm going to do it anyway
	Lines     *[]line
	Constants valueArray
}

func NewChunk() *Chunk {
	return &Chunk{
		Count:    0,
		Capacity: 0,
		Code:     nil,
	}
}

func (c *Chunk) Init() {
	c.Count = 0
	c.Capacity = 0
	c.Code = nil
	c.Lines = nil
	c.Constants.Init()
}

func (c *Chunk) Write(Byte byte, Line int) {
	if c.Capacity < c.Count+1 {
		c.Capacity = GROW_CAPACITY(c.Capacity)
		c.Code = GROW_ARRAY(c.Code, c.Capacity)
		c.Lines = GROW_ARRAY(c.Lines, c.Capacity)
	}

	(*c.Code)[c.Count] = Byte
	c.WriteLine(Line)

	c.Count += 1
}

func (c *Chunk) Free() {
	c.Init()
	c.Constants.Free()
}
