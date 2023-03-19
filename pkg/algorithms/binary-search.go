package algorithms

import (
	"errors"
	"golang.org/x/exp/constraints"
)

func LowerBound[T constraints.Ordered](array []T, target T) (int, error) {
	left := 0
	right := len(array)
	mid := len(array) / 2 // nolint:gomnd
	for left < right {
		if target <= array[mid] {
			right = mid
		} else {
			left = mid + 1
		}
		mid = (left + right) / 2 // nolint:gomnd
	}
	if left == len(array) {
		return -1, errors.New("couldn't find element not less than target")
	}
	return left, nil

}
