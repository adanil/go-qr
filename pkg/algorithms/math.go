package algorithms

import "golang.org/x/exp/constraints"

// Max is a generic function that returns the maximum value between two ordered elements.
func Max[T constraints.Ordered](a, b T) T {
	if a < b {
		return b
	}
	return a
}

// Floor is a generic function that rounds down a floating-point value to the nearest integer.
func Floor[T constraints.Float](a T) T {
	return T(int(a))
}

// Abs is a generic function that returns the absolute value of a signed or floating-point number.
func Abs[T constraints.Signed | constraints.Float](a T) T {
	if a < 0 {
		return -a
	}
	return a
}
