package utility

func Ter[T any](cond bool, a, b T) T {
	if cond {
		return a
	}

	return b
}

func PtrTo[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}
	return &v
}

func Pointer[T any](v T) *T {
	return &v
}
