package common

import "reflect"

func PointerVal[T any](val T) *T {
	return &val
}

func SafetyEmptyAsNil[T any](val T) *T {
	var zero T
	if reflect.ValueOf(&val).Elem().IsZero() {
		return nil
	}
	return &zero
}

func SafetyPointerVal[T any](val *T) T {
	var result T
	if val == nil {
		return result
	}
	return *val
}
