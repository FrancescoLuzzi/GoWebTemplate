package utils

func Ternary[T any](cond bool, vTrue, vFalse T) T {
	if cond {
		return vTrue
	} else {
		return vFalse
	}
}
