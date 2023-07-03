package algorithms

/*
ToBoolArray converts single byte into array of booleans.
Significance of bits in resulting array is in a descending order.
*/
func ToBoolArray(in byte) [8]bool {
	var out [8]bool
	for i := 7; i >= 0; i-- {
		if in>>i&1 > 0 {
			out[7-i] = true
		}
	}
	return out
}
