package glox

import "fmt"

type ObjType uint8

const (
	OBJ_STRING ObjType = iota
)

type ObjString struct {
	Length int
	Chars  []byte
	Hash   uint32
}

func (o ObjString) IsEqual(other ObjString) bool {
	if o.Length != other.Length {
		return false
	}
	if o.Hash != other.Hash {
		return false
	}
	return string(o.Chars) == string(other.Chars)
}

func (o ObjString) GetObjType() ObjType {
	return OBJ_STRING
}

func (o ObjString) Print() {
	fmt.Printf("%s", string(o.Chars))
}

func (value *Value) IsObjValue(objType ObjType) bool {
	return value.IsObj() && (*value.AsObj()).GetObjType() == objType
}

func AsObjString(value Obj) ObjString {
	if v, ok := value.(*ObjString); ok {
		return *v
	}
	return ObjString{}
}

func hashString(key string) uint32 {
	var hash uint32 = 2166136261
	for i := 0; i < len(key); i++ {
		hash ^= uint32(key[i])
		hash *= 16777619
	}
	return hash
}

func NewObjString(value string) ObjString {
	hash := hashString(value)
	return ObjString{
		Length: len(value),
		Chars:  []byte(value),
		Hash:   hash,
	}
}
