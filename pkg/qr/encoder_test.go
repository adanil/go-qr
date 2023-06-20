package qr

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_getVersion(t *testing.T) {
	testCases := []struct {
		byteLen         int
		level           Correction
		expectedVersion int
		minVersion      int
		maxVersion      int
	}{
		{
			byteLen:         10,
			level:           L,
			expectedVersion: 0,
			minVersion:      0,
			maxVersion:      40,
		},
		{
			byteLen:         100,
			level:           L,
			expectedVersion: 4,
			minVersion:      4,
			maxVersion:      5,
		},
		{
			byteLen:         1250,
			level:           M,
			expectedVersion: 28,
			minVersion:      20,
			maxVersion:      40,
		},
		{
			byteLen:         140,
			level:           H,
			expectedVersion: 11,
			minVersion:      0,
			maxVersion:      40,
		},
		{
			byteLen:         242,
			level:           Q,
			expectedVersion: 13,
			minVersion:      0,
			maxVersion:      40,
		},
		{
			byteLen:         241,
			level:           Q,
			expectedVersion: 12,
			minVersion:      0,
			maxVersion:      40,
		},
	}

	for _, test := range testCases {
		e := NewEncoder(WithCorrectionLevel(test.level), WithVersionRange(test.minVersion, test.maxVersion))
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
	e := NewEncoder(WithCorrectionLevel(L))
	version, _ := e.getVersion(len(data))
	e.version = version

	e.fillBuffer(buff, data)

	require.Equal(t, versionSize[e.level][version], buff.Len()*8)

	b := buff.Bytes()
	header := b[0] >> 4
	require.Equal(t, headerNibble, header)

	actualLen := int(b[0]&nibble | b[1]>>4)
	require.Equal(t, len(data), actualLen)

	lastByte := b[10] | b[9]&nibble
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

func BenchmarkEncode(b *testing.B) {
	inputs := [...]string{"https://github.com/psxzz/go-qr",
		"❤️ 💔 💌 💕 💞 💓 💗 💖 💘 💝 💟 💜 💛 💚 💙",
		"ثم نفس سقطت وبالتحديد،, جزيرتي باستخدام أن دنو. إذ هنا؟ الستار وتنصيب كان. أهّل ايطاليا، بريطانيا-فرنسا قد أخذ. سليمان، إتفاقية بين ما, يذكر الحدود أي بعد, معاملة بولندا، الإطلاق عل إيو.",
		"The only unicode alphabet to use a space which isn't empty but should still act like a space.",
		"Ｔｈｅ ｑｕｉｃｋ ｂｒｏｗｎ ｆｏｘ ｊｕｍｐｓ ｏｖｅｒ ｔｈｅ ｌａｚｙ ｄｏｇ",
		"ヽ༼ຈل͜ຈ༽ﾉ ヽ༼ຈل͜ຈ༽ﾉ",
		"찦차를 타고 온 펲시맨과 쑛다리 똠방각하"}

	b.ReportAllocs()
	encoder := NewEncoder()
	for i := 0; i < b.N; i++ {
		for _, inp := range inputs {
			_, err := encoder.Encode(inp)
			if err != nil {
				b.Fatal(err)
			}
		}
	}
}
