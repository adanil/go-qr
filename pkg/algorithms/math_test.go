package algorithms

import "testing"

func TestFloor(t *testing.T) {
	var f32tests = []struct {
		input float32
		want  float32
	}{
		{0.0, 0},
		{0.5, 0},
		{1.5, 1},
		{3.14, 3},
		{10.0, 10},
		{1.2345, 1},
	}

	var f64tests = []struct {
		input float64
		want  float64
	}{
		{0, 0},
		{0.5, 0},
		{1.5, 1},
		{3.14, 3},
		{10.0, 10},
		{1.2345, 1},
	}

	for _, test := range f32tests {
		if got := Floor(test.input); got != test.want {
			t.Errorf("Floor(%g) = %g", test.input, got)
		}
	}

	for _, test := range f64tests {
		if got := Floor(test.input); got != test.want {
			t.Errorf("Floor(%g) = %g", test.input, got)
		}
	}

}

func TestAbs(t *testing.T) {
	var integerTests = []struct {
		input int
		want  int
	}{
		{0, 0},
		{1, 1},
		{-2, 2},
	}
	var floatTests = []struct {
		input float64
		want  float64
	}{
		{0.0, 0.0},
		{1.1, 1.1},
		{-2.23, 2.23},
		{-3.491892348091809348, 3.491892348091809348},
	}

	for _, test := range integerTests {
		if got := Abs(test.input); got != test.want {
			t.Errorf("Abs(%d) = %d", test.input, got)
		}
	}

	for _, test := range floatTests {
		if got := Abs(test.input); got != test.want {
			t.Errorf("Abs(%g) = %g", test.input, got)
		}
	}
}
