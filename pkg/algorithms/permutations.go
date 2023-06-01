package algorithms

func GeneratePermutations(elems []int) [][2]int {
	perms := make([][2]int, 0)

	for i := 0; i < len(elems); i++ {
		for j := 0; j < len(elems); j++ {
			perms = append(perms, [2]int{elems[i], elems[j]})
		}
	}

	return perms
}

func ToBoolArray(in byte) [8]bool {
	var out [8]bool
	for i := 7; i >= 0; i-- {
		if in>>i&1 > 0 {
			out[7-i] = true
		}
	}
	return out
}

// nolint:gomnd
func ToByte(in [8]bool) byte {
	var out byte
	for i, v := range in {
		if v {
			out |= 1 << (7 - i)
		}
	}
	return out
}
