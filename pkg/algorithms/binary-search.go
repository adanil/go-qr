package algorithms

import (
	"errors"

	"golang.org/x/exp/constraints"
)

/*
LowerBound is a generic function that performs a lower bound search on a sorted array.
It searches for the first element in the array that is greater than or equal to the target value.
The function takes in a sorted array and the target value as parameters and returns the index of the found element.
If the target value is greater than all elements in the array, an error is returned.
*/
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
