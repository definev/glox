package glox

import "fmt"

type ValueType uint8

const (
	VAL_BOOL ValueType = iota
	VAL_NIL
	VAL_NUMBER
)

type Value struct {
	Type ValueType
	As   ValueData
}

type ValueData struct {
	Bool   *bool
	Number *float64
}

func NewBoolVal(value bool) Value {
	return Value{
		Type: VAL_BOOL,
		As: ValueData{
			Bool: &value,
		},
	}
}

func NewNilVal() Value {
	return Value{
		Type: VAL_NIL,
		As:   ValueData{},
	}
}

func NewNumberVal(value float64) Value {
	return Value{
		Type: VAL_NUMBER,
		As: ValueData{
			Number: &value,
		},
	}
}

func (val Value) AsBool() *bool {
	return val.As.Bool
}

func (val Value) AsNumber() *float64 {
	return val.As.Number
}

func (val Value) IsBool() bool {
	return val.Type == VAL_BOOL
}

func (val Value) IsNumber() bool {
	return val.Type == VAL_NUMBER
}

func (val Value) IsNil() bool {
	return val.Type == VAL_NIL
}

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
	switch value.Type {
	case VAL_NIL:
		fmt.Printf("nil")
	case VAL_NUMBER:
		fmt.Printf("%g", *value.AsNumber())
	case VAL_BOOL:
		fmt.Printf("%v", *value.AsBool())
	}
}
