package common

func PointerVal[T any](val T) *T {
	return &val
}
