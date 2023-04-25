package common

func SliceToMap[T any, K comparable](slice []T, keyFunc func(T) K) map[K]T {
	result := make(map[K]T)
	for _, val := range slice {
		result[keyFunc(val)] = val
	}
	return result
}
