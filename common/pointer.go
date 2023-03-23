package common

import "reflect"

func PointerVal[T any](val T) *T {
	return &val
}

func ExtractPointer[T any](val *T) T {
	var result T
	if val == nil {
		return result
	}
	return *val
}

func EmptyAsPointerNil[T any](val T) *T {
	if reflect.ValueOf(&val).Elem().IsZero() {
		return nil
	}
	return &val
}
