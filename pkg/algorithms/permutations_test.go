package algorithms

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestPermutations(t *testing.T) {
	var tests = []struct {
		input []int
		want  [][2]int
	}{
		{[]int{}, [][2]int{}},
		{[]int{1}, [][2]int{{1, 1}}},
		{[]int{1, 2}, [][2]int{{1, 1}, {1, 2}, {2, 1}, {2, 2}}},
		{[]int{1, 2, 3}, [][2]int{{1, 1}, {1, 2}, {1, 3}, {2, 1}, {2, 2}, {2, 3}, {3, 1}, {3, 2}, {3, 3}}},
	}

	for _, test := range tests {
		if got := GeneratePermutations(test.input); !slices.Equal(got, test.want) {
			t.Errorf("GeneratePermutations(%v) = %v; expected %v", test.input, got, test.want)
		}
	}
}

func TestToBoolArray(t *testing.T) {
	var tests = []struct {
		input byte
		want  [8]bool
	}{
		{0x00, [8]bool{false, false, false, false, false, false, false, false}},
		{0x01, [8]bool{false, false, false, false, false, false, false, true}},
		{0x10, [8]bool{false, false, false, true, false, false, false, false}},
		{0xAB, [8]bool{true, false, true, false, true, false, true, true}},
		{0x44, [8]bool{false, true, false, false, false, true, false, false}},
		{0x0F, [8]bool{false, false, false, false, true, true, true, true}},
		{0xF0, [8]bool{true, true, true, true, false, false, false, false}},
		{0xFF, [8]bool{true, true, true, true, true, true, true, true}},
	}

	for _, test := range tests {
		if got := ToBoolArray(test.input); got != test.want {
			t.Errorf("ToBoolArray(%v) = %v; expected %v", test.input, got, test.want)
		}
	}
}

func TestToByte(t *testing.T) {
	for _, test := range []struct {
		input [8]bool
		want  byte
	}{
		{[8]bool{false, false, false, false, false, false, false, false}, 0x00},
		{[8]bool{false, false, false, false, false, false, false, true}, 0x01},
		{[8]bool{false, false, false, true, false, false, false, false}, 0x10},
		{[8]bool{true, false, true, false, true, false, true, true}, 0xAB},
		{[8]bool{false, true, false, false, false, true, false, false}, 0x44},
		{[8]bool{false, false, false, false, true, true, true, true}, 0x0F},
		{[8]bool{true, true, true, true, false, false, false, false}, 0xF0},
		{[8]bool{true, true, true, true, true, true, true, true}, 0xFF},
	} {
		if got := ToByte(test.input); got != test.want {
			t.Errorf("ToByte(%v) = %x", test.input, got)
		}
	}
}
