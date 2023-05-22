package algorithms

import "golang.org/x/exp/constraints"

func MaxInt(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func Floor[T constraints.Float](a T) T {
	return T(int(a))
}

func Abs[T constraints.Signed | constraints.Float](a T) T {
	if a < 0 {
		return -a
	}
	return a
}
