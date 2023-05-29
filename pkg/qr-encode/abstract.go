package qr_encode

type CodeLevel int

const (
	L CodeLevel = iota
	M
	Q
	H
)

const (
	Nible        byte = 0b1111
	Byte         byte = 0xFF
	headerNibble byte = 0b0100
)

var FillerBytes = [2]byte{0b11101100, 0b00010001}
