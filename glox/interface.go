package glox

type ArrayGrowable[T any] interface {
	Init()
	Write(value T)
	Free()
}
