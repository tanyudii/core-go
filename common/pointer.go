package common

import "reflect"

func PointerVal[T any](val T) *T {
	return &val
}

func SafetyEmptyAsNil[T any](val T) *T {
	if reflect.ValueOf(&val).Elem().IsZero() {
		return nil
	}
	return &val
}

func SafetyPointerVal[T any](val *T) T {
	var result T
	if val == nil {
		return result
	}
	return *val
}
