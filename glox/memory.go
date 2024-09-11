package glox

func GROW_CAPACITY(capacity int) int {
	if capacity < 8 {
		return 8
	}
	return capacity * 2
}

func GROW_ARRAY[T any](array *[]T, capacity int) *[]T {
	var newCode []T
	if array != nil {
		newCode = make([]T, 0)
		newCode = append(newCode, *array...)
		newCode = append(newCode, make([]T, capacity-len(*array))...)
	} else {
		newCode = make([]T, capacity)
	}
	return &newCode
}
