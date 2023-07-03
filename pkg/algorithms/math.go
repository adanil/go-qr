package algorithms

import "golang.org/x/exp/constraints"

// GeneratePermutations is a function that generates permutations of length 2 for the input slice
func GeneratePermutations(elems []int) [][2]int {
	lenPerms := len(elems) * len(elems)
	perms := make([][2]int, 0, lenPerms)

	for i := 0; i < len(elems); i++ {
		for j := 0; j < len(elems); j++ {
			perms = append(perms, [2]int{elems[i], elems[j]})
		}
	}

	return perms
}

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
