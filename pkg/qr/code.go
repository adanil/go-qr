package qr

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/psxzz/go-qr/pkg/algorithms"
)

const (
	versionCodeNotRequired = 5
	syncLinePosition       = 6
)

type Code struct {
	version      int
	correction   Correction
	mask         int
	maskF        func(int, int) int
	penaltyScore int

	alignments []int

	canvas [][]Module
	size   int
}

func newCode(data []byte, correction Correction, version int, mask int) *Code {
	canvasSize := 4*(version+1) + 17 // nolint:gomnd

	var canvas [][]Module = make([][]Module, canvasSize)
	for i := range canvas {
		canvas[i] = make([]Module, canvasSize)
	}

	code := &Code{
		version:      version,
		correction:   correction,
		mask:         mask,
		maskF:        maskFunctions[mask],
		penaltyScore: 0,
		canvas:       canvas,
		size:         canvasSize,
		alignments:   alignmentPatterns[version],
	}

	code.encode(data)

	return code
}

func (c *Code) String() string {
	var buf bytes.Buffer

	buf.WriteByte('{')
	fmt.Fprintf(&buf, "\nsize: %v", c.size)
	fmt.Fprintf(&buf, "\nversion: %v", c.version)
	fmt.Fprintf(&buf, "\nerror correction: %v", c.correction)
	fmt.Fprintf(&buf, "\nmask pattern: %v", c.mask)
	fmt.Fprintf(&buf, "\nalignments: %v", c.alignments)
	fmt.Fprintf(&buf, "\ndata: ")

	for _, row := range c.canvas {
		buf.WriteString("\n\t\t")
		for _, v := range row {
			fmt.Fprintf(&buf, "%v", v.String())
		}
	}

	buf.WriteString("\n}")

	return buf.String()
}

// GetImageWithColors generates an image representation of the QR code with specified colors
func (c *Code) GetImageWithColors(pixelSize int, colorOne, colorTwo color.RGBA) image.Image {
	canvasHeight, canvasWidth := len(c.canvas), len(c.canvas[0])
	imageHeight, imageWidth := canvasHeight*pixelSize, canvasWidth*pixelSize

	upLeft, lowRight := image.Point{X: 0, Y: 0}, image.Point{X: imageWidth, Y: imageHeight}

	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	r := image.Rect(0, 0, imageWidth, imageHeight)
	draw.Draw(img, r, image.NewUniform(colorOne), image.Point{}, draw.Src)

	for y := 0; y < canvasHeight; y++ {
		for x := 0; x < canvasWidth; x++ {
			if !c.canvas[y][x].value {
				continue
			}
			rect := image.Rect(x*pixelSize, y*pixelSize, (x+1)*pixelSize, (y+1)*pixelSize)
			draw.Draw(img, rect, image.NewUniform(colorTwo), image.Point{}, draw.Src)
		}
	}

	return img
}

// GetImage generates an image representation of the QR code using default colors (black and white)
func (c *Code) GetImage(pixelSize int) image.Image {
	return c.GetImageWithColors(pixelSize, color.RGBA{R: 255, G: 255, B: 255, A: 0xff}, color.RGBA{A: 0xff}) //nolint:gomnd
}

func (c *Code) encode(data []byte) {
	c.placeSearchPatterns()
	c.placeAlignments()
	c.placeSync()

	if c.version > versionCodeNotRequired {
		c.placeVersion()
	}

	c.placeMask()
	c.writeData(data)
	c.countPenalty()
}

func (c *Code) writeData(bytes []byte) {
	mask := c.maskF
	nextBit := c.bitsGenerator(bytes) // convert encoded data to bit flow

	xl, xr := c.size-2, c.size-1 // nolint:gomnd
	upwards := true
	for xl >= 0 {
		if xr == syncLinePosition { // skip vertical synchronization line
			xl, xr = xl-1, xr-1
		}

		y, border := c.size-1, -1
		if !upwards {
			y, border = 0, c.size
		}

		for y != border {
			if !c.canvas[y][xr].isSet {
				bit := nextBit()

				if mask(xr, y) == 0 {
					bit = !bit
				}

				c.canvas[y][xr].Set(bit)
			}

			if !c.canvas[y][xl].isSet {
				bit := nextBit()

				if mask(xl, y) == 0 {
					bit = !bit
				}

				c.canvas[y][xl].Set(bit)
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

func (c *Code) bitsGenerator(data []byte) func() bool {
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

func (c *Code) countPenalty() int {
	// TODO implement
	return 0
}

func (c *Code) placeSearchPatterns() {
	c.placePattern(0, 0, &searchPatternTL)                            // Top left corner
	c.placePattern(0, c.size-searchPatternTR.xSize, &searchPatternTR) // Top right
	c.placePattern(c.size-searchPatternBL.ySize, 0, &searchPatternBL) // Bottom left
}

func (c *Code) placeAlignments() {
	perms := algorithms.GeneratePermutations(c.alignments)
	offset := alignmentPatternSize / 2 // nolint:gomnd

	for _, loc := range perms {
		x, y := loc[0]-offset, loc[1]-offset
		c.placePattern(x, y, &alignmentPattern)
	}
}

func (c *Code) placeSync() {
	syncPixels := [2]bool{bl, wh}
	lenSyncPixels := len(syncPixels)
	syncEnd := c.size - 7 // nolint:gomnd

	var i, locX, locY int
	// Vertical sync border
	for i, locX, locY = 0, 6, 8; locY < syncEnd; locY++ {
		if !c.canvas[locY][locX].isSet {
			c.canvas[locY][locX].Set(syncPixels[i])
		}
		i = (i + 1) % lenSyncPixels
	}

	// Horizontal sync border
	for i, locX, locY = 0, 8, 6; locX < syncEnd; locX++ {
		if !c.canvas[locY][locX].isSet {
			c.canvas[locY][locX].Set(syncPixels[i])
		}
		i = (i + 1) % lenSyncPixels
	}
}

func (c *Code) placeVersion() {
	versionPadding := 11
	locX, locY := 0, c.size-versionPadding

	versionBits := versionCodes[c.version]
	for y_offset, b := range versionBits {
		bits := algorithms.ToBoolArray(b)

		for x_offset, bit := range bits[2:] {
			x, y := locX+x_offset, locY+y_offset

			c.canvas[y][x].Set(bit) // Bottom left code
			c.canvas[x][y].Set(bit) // Top right code
		}
	}

}

// nolint:gomnd
func (c *Code) placeMask() {
	maskCode := maskCodes[c.correction][c.mask]

	codeBits := make([]bool, 0, 15)
	msb := algorithms.ToBoolArray(byte(maskCode >> 8))
	lsb := algorithms.ToBoolArray(byte(maskCode))

	codeBits = append(codeBits, msb[1:]...)
	codeBits = append(codeBits, lsb[:]...)

	// Bottom left + Top right
	i := 0
	for x, y := 8, c.size-1; y > c.size-8; y-- {
		c.canvas[y][x].Set(codeBits[i])
		i++
	}

	c.canvas[c.size-8][8].Set(false) // This module is always black

	for x, y := c.size-8, 8; x < c.size; x++ {
		c.canvas[y][x].Set(codeBits[i])
		i++
	}

	// Top left
	i = 0
	for x, y := 0, 8; x < 9; x++ {
		if !c.canvas[y][x].isSet {
			c.canvas[y][x].Set(codeBits[i])
			i++
		}
	}

	for x, y := 8, 7; y > -1; y-- {
		if !c.canvas[y][x].isSet {
			c.canvas[y][x].Set(codeBits[i])
			i++
		}
	}

}

func (c *Code) isUnused(startX, startY, endX, endY int) bool {

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

func (c *Code) placePattern(locX, locY int, p *Pattern) {
	pxLen, pyLen := p.xSize, p.ySize

	if !c.isUnused(locX, locY, locX+pxLen, locY+pyLen) {
		return
	}

	for i, pi := locX, 0; i < c.size && pi < pxLen; i, pi = i+1, pi+1 {
		for j, pj := locY, 0; j < c.size && pj < pyLen; j, pj = j+1, pj+1 {
			c.canvas[i][j].Set(p.data[pi][pj])
		}
	}
}
