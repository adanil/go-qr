package qr_encode

type CodeLevel int

const (
	L CodeLevel = iota
	M
	Q
	H
)

const (
	Nible        int = 0b1111
	Byte         int = 0xFF
	headerNibble int = 0b0100
)

var FillerBytes = [2]byte{0b11101100, 0b00010001}

type Encoder struct {
	level   CodeLevel
	version int
}
