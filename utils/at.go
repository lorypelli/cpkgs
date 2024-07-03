package utils

func At[T any](arr []T, n int) T {
	len := len(arr)
	return arr[(n+len)%len]
}