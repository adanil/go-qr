package qr

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/psxzz/go-qr/pkg/algorithms"
	"go.uber.org/multierr"
	"golang.org/x/exp/slices"
)

const (
	versionCodeNotRequired = 5
	timingPosition         = 6
	versionCodeOffset      = 11
	baseScorePenalty1      = 3
	baseScorePenalty2      = 3
	baseScorePenalty3      = 40
)

// Correction is the level of QR code correction: L, M, Q, H
type Correction int

// Encoder encodes input data into a QR code
type Encoder struct {
	level                  Correction
	minVersion, maxVersion int
	minMask, maxMask       int
	version                int
}

// Encode encodes the given text into a QR code
func (e *Encoder) Encode(text string) (*Code, error) {
	data, err := e.dataEncode(text)

	if err != nil {
		return nil, fmt.Errorf("runtime error in data_encoder: %w", err)
	}

	return e.generateCode(data), nil
}

func (e *Encoder) dataEncode(text string) ([]byte, error) {
	byteLen := len(text)
	codeVersion, err := e.getVersion(byteLen)
	if err != nil {
		return nil, multierr.Combine(ErrVersionNotFound, err)
	}
	e.version = codeVersion

	currBuff := bytes.NewBuffer(make([]byte, 0, len(text)+10)) // nolint:gomnd
	e.fillBuffer(currBuff, []byte(text))

	blocks := e.divideIntoBlocks(currBuff)
	correctionBlocks := e.generateCorrectionBlocks(blocks)
	result := e.mergeBlocks(blocks, correctionBlocks)

	return result, nil
}

func (e *Encoder) generateCode(data []byte) *Code {
	var (
		currentCode *Code
		wg          sync.WaitGroup
	)
	codes := make(chan *Code, e.maxMask-e.minMask)

	wg.Add(e.maxMask - e.minMask)
	for mask := e.minMask; mask < e.maxMask; mask++ {
		go func(mask int) {
			defer wg.Done()
			code := newCode(data, e.level, e.version, mask)

			e.placeFinderPatterns(code)
			e.placeAlignments(code)
			e.placeTimings(code)

			if e.version > versionCodeNotRequired {
				e.placeVersion(code)
			}

			e.placeMask(code)
			e.placeData(code, data)
			e.countPenalty(code)

			codes <- code
		}(mask)
	}

	wg.Wait()
	close(codes)

	for code := range codes {
		if currentCode == nil || code.penaltyScore < currentCode.penaltyScore {
			currentCode = code
		}
	}

	return currentCode
}

func (e *Encoder) getVersion(byteLen int) (int, error) {
	bitLen := byteLen*8 + 4 // nolint:gomnd
	versionsArray := versionSize[e.level]
	version, err := algorithms.LowerBound(versionsArray[e.minVersion:e.maxVersion], bitLen)
	if err != nil {
		return -1, err
	}
	version += e.minVersion
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
		buff.WriteByte((headerNibble << 4) | ((dataLen >> 4) & nibble))
		currByte = dataLen & nibble
	} else {
		dataLen := uint16(len(data))
		buff.WriteByte((headerNibble << 4) | (byte(dataLen>>12) & nibble))
		buff.WriteByte(byte(dataLen >> 4))
		currByte = byte(dataLen) & nibble
	}

	for _, b := range data {
		currByte = (currByte << 4) | ((b >> 4) & nibble)
		buff.WriteByte(currByte)
		currByte = b & nibble
	}
	currByte <<= 4
	buff.WriteByte(currByte)

	idx := 0
	currByte = fillerBytes[idx]
	for buff.Len()*8 < versionSize[e.level][e.version] {
		buff.WriteByte(currByte)
		idx = (idx + 1) % 2
		currByte = fillerBytes[idx]
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
		correctionBytesNum := algorithms.Max(len(block), coefficientsNum)
		correctionBytes := make([]byte, correctionBytesNum+2*len(block))
		copy(correctionBytes, block)

		for i := 0; i < len(block); i++ {
			a := correctionBytes[0]
			correctionBytes = correctionBytes[1:]

			if a == 0 {
				continue
			}

			for j, c := range coefficients {
				correctionBytes[j] ^= GF[c+invGF[a]]
			}
		}

		result = append(result, correctionBytes[:coefficientsNum])
	}

	return result
}

func (e *Encoder) mergeBlocks(blocks [][]byte, correctionBlocks [][]byte) []byte {
	maxBlockSize := 0
	for _, block := range blocks {
		maxBlockSize = algorithms.Max(maxBlockSize, len(block))
	}

	result := bytes.NewBuffer(make([]byte, 0, 2*maxBlockSize*len(blocks)))

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
		maxBlockSize = algorithms.Max(maxBlockSize, len(corrBlock))
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

func (e *Encoder) placeFinderPatterns(code *Code) {
	e.placePattern(code, 0, 0, &finderPatternTL)                               // Top left corner
	e.placePattern(code, 0, code.size-finderPatternTR.xSize, &finderPatternTR) // Top right
	e.placePattern(code, code.size-finderPatternBL.ySize, 0, &finderPatternBL) // Bottom left
}

func (e *Encoder) placeAlignments(code *Code) {
	perms := algorithms.GeneratePermutations(code.alignments)
	offset := alignmentPatternSize / 2 // nolint:gomnd

	for _, loc := range perms {
		x, y := loc[0]-offset, loc[1]-offset
		e.placePattern(code, x, y, &alignmentPattern)
	}
}

func (e *Encoder) placeTimings(code *Code) {
	lenTimingPixels := len(timingPixels)
	timingEnd := code.size - 7 // nolint:gomnd

	var i, locX, locY int
	// Vertical sync border
	for i, locX, locY = 0, 6, 8; locY < timingEnd; locY++ {
		if !code.canvas[locY][locX].isSet {
			code.canvas[locY][locX].Set(timingPixels[i])
		}
		i = (i + 1) % lenTimingPixels
	}

	// Horizontal sync border
	for i, locX, locY = 0, 8, 6; locX < timingEnd; locX++ {
		if !code.canvas[locY][locX].isSet {
			code.canvas[locY][locX].Set(timingPixels[i])
		}
		i = (i + 1) % lenTimingPixels
	}
}

func (e *Encoder) placeVersion(code *Code) {
	startX, startY := 0, code.size-versionCodeOffset

	versionBits := versionCodes[code.version]
	for y_offset, b := range versionBits {
		bits := algorithms.ToBoolArray(b)

		for x_offset, bit := range bits[2:] {
			x, y := startX+x_offset, startY+y_offset

			code.canvas[y][x].Set(bit) // Bottom left code
			code.canvas[x][y].Set(bit) // Top right code
		}
	}

}

// nolint:gomnd
func (e *Encoder) placeMask(code *Code) {
	maskCode := maskCodes[code.correction][code.mask]

	codeBits := make([]bool, 0, 15)
	msb := algorithms.ToBoolArray(byte(maskCode >> 8))
	lsb := algorithms.ToBoolArray(byte(maskCode))

	codeBits = append(codeBits, msb[1:]...)
	codeBits = append(codeBits, lsb[:]...)

	// Bottom left + Top right
	i := 0
	for x, y := 8, code.size-1; y > code.size-8; y-- {
		code.canvas[y][x].Set(codeBits[i])
		i++
	}

	code.canvas[code.size-8][8].Set(true) // This module is always black

	for x, y := code.size-8, 8; x < code.size; x++ {
		code.canvas[y][x].Set(codeBits[i])
		i++
	}

	// Top left
	i = 0
	for x, y := 0, 8; x < 9; x++ {
		if !code.canvas[y][x].isSet {
			code.canvas[y][x].Set(codeBits[i])
			i++
		}
	}

	for x, y := 8, 7; y > -1; y-- {
		if !code.canvas[y][x].isSet {
			code.canvas[y][x].Set(codeBits[i])
			i++
		}
	}

}

func (e *Encoder) placeData(code *Code, bytes []byte) {
	mask := code.maskF
	nextBit := e.bitFlow(bytes) // convert encoded data to bit flow

	xl, xr := code.size-2, code.size-1 // nolint:gomnd
	upwards := true
	for xl >= 0 {
		if xr == timingPosition { // skip vertical timing
			xl, xr = xl-1, xr-1
		}

		y, border := code.size-1, -1
		if !upwards {
			y, border = 0, code.size
		}

		for y != border {
			if !code.canvas[y][xr].isSet {
				bit := nextBit()

				if mask(xr, y) == 0 {
					bit = !bit
				}

				code.canvas[y][xr].Set(bit)
			}

			if !code.canvas[y][xl].isSet {
				bit := nextBit()

				if mask(xl, y) == 0 {
					bit = !bit
				}

				code.canvas[y][xl].Set(bit)
			}

			if upwards {
				y--
			} else {
				y++
			}
		}

		xl, xr = xl-2, xr-2 // nolint:gomnd
		upwards = !upwards
	}
}

func (e *Encoder) countPenalty(code *Code) {
	code.penaltyScore = e.penalty1(code) + e.penalty2(code) +
		e.penalty3(code) + e.penalty4(code)
}

func (e *Encoder) placePattern(c *Code, startX, startY int, p *Pattern) {
	pxLen, pyLen := p.xSize, p.ySize

	if !e.isUnused(c, startX, startY, startX+pxLen, startY+pyLen) {
		return
	}

	for i, pi := startX, 0; i < c.size && pi < pxLen; i, pi = i+1, pi+1 {
		for j, pj := startY, 0; j < c.size && pj < pyLen; j, pj = j+1, pj+1 {
			c.canvas[i][j].Set(p.data[pi][pj])
		}
	}
}

func (e *Encoder) isUnused(c *Code, startX, startY, endX, endY int) bool {

	// false, if arguments are out of canvas bounds
	if startX < 0 || startX >= c.size ||
		startY < 0 || startY >= c.size ||
		endX > c.size || endY > c.size {
		return false
	}

	for i := startX; i < endX; i++ {
		for j := startY; j < endY; j++ {
			if c.canvas[i][j].isSet {
				return false
			}
		}
	}

	return true
}

func (e *Encoder) bitFlow(data []byte) func() bool {
	dataBits := make([]bool, 0, len(data)*8) // nolint:gomnd
	for _, b := range data {
		bits := algorithms.ToBoolArray(b)
		dataBits = append(dataBits, bits[:]...)
	}

	i := 0
	return func() bool {
		if i >= len(dataBits) {
			return false
		}

		bit := dataBits[i]
		i++
		return bit
	}
}

// nolint:gomnd
func (e *Encoder) penalty1(c *Code) int {
	var score int

	// rows
	for _, row := range c.canvas {
		prev := row[0].value
		count := 1

		for _, m := range row[1:] {
			if m.value != prev {
				prev, count = m.value, 1
				continue
			}

			count++
			if count == 5 {
				score += baseScorePenalty2
			} else if count > 5 {
				score++
			}
		}
	}

	// columns
	for i := 0; i < len(c.canvas); i++ {
		prev := c.canvas[0][i].value
		count := 1

		for j := 1; j < len(c.canvas); j++ {
			if c.canvas[j][i].value != prev {
				prev, count = c.canvas[j][i].value, 1
				continue
			}

			count++
			if count == 5 {
				score += baseScorePenalty1
			} else if count > 5 {
				score++
			}
		}
	}

	return score
}

func (e *Encoder) penalty2(c *Code) int {
	var blocks int

	for i := 0; i < len(c.canvas)-1; i++ {
		for j := 0; j < len(c.canvas)-1; j++ {
			curr := c.canvas[i][j].value

			adjacent := c.canvas[i][j+1].value &&
				c.canvas[i+1][j].value &&
				c.canvas[i+1][j+1].value

			if curr == (curr && adjacent) {
				blocks++
			}
		}
	}

	return baseScorePenalty2 * blocks
}

func (e *Encoder) penalty3(c *Code) int {
	var patternsCount int
	steps := len(c.canvas) - len(penalty3Pattern)

	// rows
	for _, row := range c.canvas {
		for x := 0; x < steps; x++ {
			curr := make([]bool, 0, len(penalty3Pattern))
			for _, m := range row[x : x+len(penalty3Pattern)] {
				curr = append(curr, m.value)
			}

			if slices.Equal(penalty3Pattern[:], curr) ||
				slices.Equal(penalty3PatternReversed[:], curr) {
				patternsCount++
			}
		}
	}

	// columns
	for x := 0; x < len(c.canvas); x++ {
		for y := 0; y < steps; y++ {
			curr := make([]bool, 0, len(penalty3Pattern))
			for i := y; i < len(penalty3Pattern); i++ {
				curr = append(curr, c.canvas[i][x].value)
			}

			if slices.Equal(penalty3Pattern[:], curr) ||
				slices.Equal(penalty3PatternReversed[:], curr) {
				patternsCount++
			}
		}
	}

	return patternsCount * baseScorePenalty3
}

// nolint:gomnd
func (e *Encoder) penalty4(c *Code) int {
	var score int
	total := len(c.canvas) * len(c.canvas)
	black := 0

	for _, row := range c.canvas {
		for _, m := range row {
			if !m.value {
				black++
			}
		}
	}

	ratio := float64(black) / float64(total)
	ratio = algorithms.Floor(ratio*100 - 50)
	score = algorithms.Abs(int(ratio)) * 2

	return score
}
