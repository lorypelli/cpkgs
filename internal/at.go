package internal

func At[T any](arr []T, n int) T {
	l := len(arr)
	return arr[(n+l)%l]
}
