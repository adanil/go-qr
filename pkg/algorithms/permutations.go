package algorithms

func GeneratePermutations(elems []uint) [][2]uint {
	perms := make([][2]uint, 0)

	for i := 0; i < len(elems); i++ {
		for j := 0; j < len(elems); j++ {
			perms = append(perms, [2]uint{elems[i], elems[j]})
		}
	}

	return perms
}

// TODO: rewrite for variadic binary
func ToBoolArray(in byte) [6]bool {
	var out [6]bool
	for i := 5; i >= 0; i-- {
		if in>>i&1 > 0 {
			out[5-i] = true
		}
	}
	return out
}
