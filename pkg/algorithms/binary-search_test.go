package algorithms

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/constraints"
	"testing"
)

type testCase[T constraints.Ordered] struct {
	Array         []T
	TargetElement T
	Expected      int
}

func Test_LowerBound(t *testing.T) {
	testInts := [...]testCase[int]{
		{
			Array:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			TargetElement: 5,
			Expected:      4,
		},
		{
			Array:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			TargetElement: 1,
			Expected:      0,
		},
		{
			Array:         []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			TargetElement: 23,
			Expected:      -1,
		},
		{
			Array:         []int{1, 2, 3, 3, 3, 62, 72, 72, 83, 91},
			TargetElement: 3,
			Expected:      2,
		},
		{
			Array:         []int{1, 2, 3, 3, 3, 62, 72, 72, 83, 91},
			TargetElement: 72,
			Expected:      6,
		},
	}

	testFloats := [...]testCase[float64]{
		{
			Array:         []float64{1.13, 1.15, 3.14, 4.28, 5.56, 6.66, 6.66, 6.66, 9.33, 10.0},
			TargetElement: 5,
			Expected:      4,
		},
		{
			Array:         []float64{1.13, 1.15, 3.14, 4.28, 5.56, 6.66, 6.66, 6.66, 9.33, 10.0},
			TargetElement: 1,
			Expected:      0,
		},
		{
			Array:         []float64{1.13, 1.15, 3.14, 4.28, 5.56, 6.66, 6.66, 6.66, 9.33, 10.0},
			TargetElement: 6.66,
			Expected:      5,
		},
		{
			Array:         []float64{1.13, 1.15, 3.14, 4.28, 5.56, 6.66, 6.66, 6.66, 9.33, 10.0},
			TargetElement: 33.47,
			Expected:      -1,
		},
		{
			Array:         []float64{1, 2, 3, 3, 3, 62, 72, 72, 83, 91},
			TargetElement: 72,
			Expected:      6,
		},
	}

	for _, test := range testInts {
		actual, _ := LowerBound(test.Array, test.TargetElement)
		require.Equal(t, test.Expected, actual)
	}

	for _, test := range testFloats {
		actual, _ := LowerBound(test.Array, test.TargetElement)
		require.Equal(t, test.Expected, actual)
	}

}
