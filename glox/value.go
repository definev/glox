package glox

import "fmt"

type Value float64

type valueArray struct {
	Count    int
	Capacity int
	Values   *[]Value
}

func NewValueArray() *valueArray {
	return &valueArray{
		Count:    0,
		Capacity: 0,
		Values:   nil,
	}
}

func (v *valueArray) Init() {
	v.Count = 0
	v.Capacity = 0
	v.Values = nil
}

func (v *valueArray) Write(value Value) {
	if v.Capacity < v.Count+1 {
		v.Capacity = GROW_CAPACITY(v.Capacity)
		v.Values = GROW_ARRAY(v.Values, v.Capacity)
	}

	(*v.Values)[v.Count] = value
	v.Count += 1
}

func (v *valueArray) Free() {
	v.Init()
}

func (v *valueArray) Print(index int) {
	PrintValue((*v.Values)[index])
}

func PrintValue(value Value) {
	fmt.Printf("%g", value)
}
