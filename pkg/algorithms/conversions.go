package algorithms

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
