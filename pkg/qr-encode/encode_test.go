package qr_encode

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_getVersion(t *testing.T) {
	testCases := []struct {
		byteLen         int
		level           CodeLevel
		expectedVersion int
	}{
		{
			byteLen:         100,
			level:           L,
			expectedVersion: 4,
		},
		{
			byteLen:         1250,
			level:           M,
			expectedVersion: 28,
		},
		{
			byteLen:         140,
			level:           H,
			expectedVersion: 11,
		},
		{
			byteLen:         242,
			level:           Q,
			expectedVersion: 13,
		},
		{
			byteLen:         241,
			level:           Q,
			expectedVersion: 12,
		},
	}

	for _, test := range testCases {
		e := Encoder{level: test.level}
		actual, err := e.getVersion(test.byteLen)
		require.NoError(t, err)
		require.Equal(t, test.expectedVersion, actual)
	}

	e := Encoder{level: L}
	actual, err := e.getVersion(2954)
	require.Equal(t, -1, actual)
	require.NotNil(t, err)

}

func Test_fillBuffer(t *testing.T) {
	buff := bytes.NewBuffer(make([]byte, 0))
	data := []byte{13, 14, 28, 42, 56, 88, 123, 233, 255}
	e := Encoder{level: L}
	version, _ := e.getVersion(len(data))
	e.version = version

	e.fillBuffer(buff, data)

	require.Equal(t, versionSize[e.level][version], buff.Len()*8)

	b := buff.Bytes()
	header := int(b[0] >> 4)
	require.Equal(t, headerNibble, header)

	actualLen := int(byte(int(b[0])&Nible) | ((b[1] >> 4) << 4))
	require.Equal(t, len(data), actualLen)

	lastByte := b[10] | byte(int(b[9])&Nible)
	require.Equal(t, data[len(data)-1], lastByte)

}

func Test_divideIntoBlocks(t *testing.T) {
	buf := bytes.NewBufferString("0123456789ABCDEF")
	expected := [][]byte{
		[]byte("0123456789ABCDEF"),
	}
	e := Encoder{level: M, version: 0}
	result := e.divideIntoBlocks(buf)
	require.Equal(t, expected, result)

	buf = bytes.NewBufferString("0123456789ABCDEFGH")
	expected = [][]byte{
		[]byte("0123"),
		[]byte("4567"),
		[]byte("89ABC"),
		[]byte("DEFGH"),
	}
	e = Encoder{level: H, version: 3}
	result = e.divideIntoBlocks(buf)
	require.Equal(t, expected, result)
}

func Test_generateCorrectionBlocks(t *testing.T) {
	input := [][]byte{
		{64, 196, 132, 84, 196, 196, 242, 194, 4, 132, 20, 37, 34, 16, 236, 17},
	}
	expected := [][]byte{
		{16, 85, 12, 231, 54, 54, 140, 70, 118, 84, 10, 174, 235, 197, 99, 218, 12, 254, 246, 4, 190, 56, 39, 217, 115, 189, 193, 24},
	}

	e := Encoder{level: H, version: 1}
	result := e.generateCorrectionBlocks(input)

	require.Equal(t, expected, result)
}

func Test_mergeBlocks(t *testing.T) {
	blocks1 := [][]byte{{0x01, 0x02, 0x03}, {0x04, 0x05, 0x06}}
	correctionBlocks1 := [][]byte{{0x07, 0x08, 0x09}, {0x0A, 0x0B, 0x0C}}
	expected1 := []byte{0x01, 0x04, 0x02, 0x05, 0x03, 0x06, 0x07, 0x0A, 0x08, 0x0B, 0x09, 0x0C}

	e := Encoder{}
	result := e.mergeBlocks(blocks1, correctionBlocks1)
	require.Equal(t, expected1, result)
}
