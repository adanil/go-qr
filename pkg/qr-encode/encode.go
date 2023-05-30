package qr_encode

import (
	"bytes"
	"fmt"

	"github.com/psxzz/go-qr/pkg/algorithms"
	"go.uber.org/multierr"
)

type Encoder struct {
	level   CodeLevel
	version int
}

func NewEncoder(level CodeLevel) *Encoder {
	return &Encoder{level: level}
}

func (e *Encoder) Encode(text string) ([]byte, error) {
	byteLen := len(text)
	codeVersion, err := e.getVersion(byteLen)
	if err != nil {
		return nil, multierr.Combine(ErrVersionNotFound, err)
	}
	e.version = codeVersion

	currBuff := bytes.NewBuffer(make([]byte, 0))
	e.fillBuffer(currBuff, []byte(text))

	blocks := e.divideIntoBlocks(currBuff)
	correctionBlocks := e.generateCorrectionBlocks(blocks)
	result := e.mergeBlocks(blocks, correctionBlocks)

	return result, nil
}

func (e *Encoder) Encode2D(text string) (*QRCode, error) {
	data, err := e.Encode(text)

	if err != nil {
		return nil, fmt.Errorf("encode1d: %v", err)
	}

	grid := NewQRCode(e, data)
	grid.MakeLayout()

	_, err = grid.Write(data)

	if err != nil {
		return nil, fmt.Errorf("write: %v", err)
	}

	return grid, nil
}

func (e *Encoder) getVersion(byteLen int) (int, error) {
	bitLen := byteLen*8 + 4 // nolint:gomnd
	versionsArray := versionSize[e.level]
	version, err := algorithms.LowerBound(versionsArray[:], bitLen)
	if err != nil {
		return -1, err
	}
	maxVersionSize := versionSize[e.level][version]
	if version < 9 { // nolint:gomnd
		bitLen += 8
	} else {
		bitLen += 16
	}
	if bitLen > maxVersionSize {
		version++
		if version >= len(versionSize[e.level]) {
			return -1, ErrTooLargeSize
		}
	}
	return version, nil
}

// nolint:gomnd
func (e *Encoder) fillBuffer(buff *bytes.Buffer, data []byte) {
	var currByte byte

	if e.version < 9 {
		dataLen := uint8(len(data))
		buff.WriteByte((headerNibble << 4) | ((dataLen >> 4) & Nible))
		currByte = dataLen & Nible
	} else {
		dataLen := uint16(len(data))
		buff.WriteByte((headerNibble << 4) | (byte(dataLen>>12) & Nible))
		buff.WriteByte(byte(dataLen >> 4))
		currByte = byte(dataLen) & Nible
	}

	for _, b := range data {
		currByte = (currByte << 4) | ((b >> 4) & Nible)
		buff.WriteByte(currByte)
		currByte = b & Nible
	}
	currByte <<= 4
	buff.WriteByte(currByte)

	idx := 0
	currByte = FillerBytes[idx]
	for buff.Len()*8 < versionSize[e.level][e.version] {
		buff.WriteByte(currByte)
		idx = (idx + 1) % 2
		currByte = FillerBytes[idx]
	}
}

func (e *Encoder) divideIntoBlocks(buff *bytes.Buffer) [][]byte {
	blocksNum := numberOfBlocks[e.level][e.version]
	blockSize := buff.Len() / blocksNum
	rem := buff.Len() % blocksNum
	result := make([][]byte, blocksNum)

	data := buff.Bytes()
	currIdx := 0
	for i := 0; i < blocksNum-rem; i++ {
		result[i] = data[currIdx : currIdx+blockSize]
		currIdx += blockSize
	}
	for i := 0; i < rem; i++ {
		result[blocksNum-rem+i] = data[currIdx : currIdx+blockSize+1]
		currIdx += blockSize + 1
	}

	return result
}

// nolint:gomnd
func (e *Encoder) generateCorrectionBlocks(dataBlocks [][]byte) [][]byte {
	coefficientsNum := numberOfCorrectionBytes[e.level][e.version]
	coefficients := polynomialCoefficients[coefficientsNum]

	result := make([][]byte, 0, len(dataBlocks))
	for _, block := range dataBlocks {
		correctionBytesNum := algorithms.MaxInt(len(block), coefficientsNum)
		correctionBytes := make([]byte, 0, correctionBytesNum+len(block))
		correctionBytes = append(correctionBytes, block...)

		for i := len(correctionBytes); i < correctionBytesNum; i++ {
			correctionBytes = append(correctionBytes, 0)
		}

		for i := 0; i < len(block); i++ {
			a := correctionBytes[0]
			correctionBytes = append(correctionBytes[1:], 0)

			if a == 0 {
				continue
			}

			b := int(invGF[a])
			for j := 0; j < coefficientsNum; j++ {
				c := (coefficients[j] + b) % 255
				t := GF[c]
				correctionBytes[j] ^= t
			}
		}

		result = append(result, correctionBytes[:coefficientsNum])
	}

	return result
}

func (e *Encoder) mergeBlocks(blocks [][]byte, correctionBlocks [][]byte) []byte {
	result := bytes.NewBuffer(make([]byte, 0))

	maxBlockSize := 0
	for _, block := range blocks {
		maxBlockSize = algorithms.MaxInt(maxBlockSize, len(block))
	}

	currByteIdx := 0
	for currByteIdx < maxBlockSize {
		for i := 0; i < len(blocks); i++ {
			if currByteIdx >= len(blocks[i]) {
				continue
			}
			result.WriteByte(blocks[i][currByteIdx])
		}
		currByteIdx++
	}

	maxBlockSize = 0
	for _, corrBlock := range correctionBlocks {
		maxBlockSize = algorithms.MaxInt(maxBlockSize, len(corrBlock))
	}

	currByteIdx = 0
	for currByteIdx < maxBlockSize {
		for i := 0; i < len(correctionBlocks); i++ {
			if currByteIdx >= len(correctionBlocks[i]) {
				continue
			}
			result.WriteByte(correctionBlocks[i][currByteIdx])
		}
		currByteIdx++
	}

	return result.Bytes()
}
