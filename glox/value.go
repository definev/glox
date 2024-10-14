package glox

import "fmt"

type ValueType uint8

const (
	VAL_BOOL ValueType = iota
	VAL_NIL
	VAL_NUMBER
	VAL_OBJ
)

type Obj interface {
	GetObjType() ObjType
	Print()
}

type Value struct {
	Type ValueType
	As   ValueData
}

type ValueData struct {
	Bool   *bool
	Number *float64
	Obj    *Obj
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

func NewObjVal(value Obj) Value {
	return Value{
		Type: VAL_OBJ,
		As: ValueData{
			Obj: &value,
		},
	}
}

func (val Value) AsBool() *bool {
	return val.As.Bool
}

func (val Value) AsNumber() *float64 {
	return val.As.Number
}

func (val Value) AsObj() *Obj {
	return val.As.Obj
}

func (val Value) AsString() *ObjString {
	if val.Type != VAL_OBJ {
		return nil
	}
	objString := *val.As.Obj
	if v, ok := objString.(ObjString); ok {
		return &v
	}
	return nil
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

func (val Value) IsObj() bool {
	return val.Type == VAL_OBJ
}

func (val Value) IsString() bool {
	return val.Type == VAL_OBJ && (*val.AsObj()).GetObjType() == OBJ_STRING
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

func (a Value) IsEqual(b Value) bool {
	if a.Type != b.Type {
		return false
	}

	switch a.Type {
	case VAL_NIL:
		return true
	case VAL_BOOL:
		return *a.AsBool() == *b.AsBool()
	case VAL_NUMBER:
		return *a.AsNumber() == *b.AsNumber()
	case VAL_OBJ:
		aObj := *a.AsObj()
		bObj := *b.AsObj()

		if aObj == nil && bObj == nil {
			return true
		}
		if aObj.GetObjType() != bObj.GetObjType() {
			return false
		}

		switch aObj.GetObjType() {
		case OBJ_STRING:
			aObjStr := AsObjString(aObj)
			bObjStr := AsObjString(bObj)

			return string(aObjStr.Chars) == string(bObjStr.Chars)
		}

		return false
	default:
		return false
	}
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
	case VAL_OBJ:
		(*value.AsObj()).Print()
	}
}
