package qr_encode

import (
	"bytes"
	"fmt"

	"github.com/psxzz/go-qr/pkg/algorithms"
)

// TODO: Get rid of magic numbers and make it constant
const (
	versionOneSize     = 21
	versionSizePadding = 7
	versionPadding     = 11
	syncPadding        = 7
)

type Pixel bool

func (q Pixel) String() string {
	if q {
		return "██"
	}
	return "  "
}

type Module struct {
	value  Pixel
	isUsed bool
}

// Sets the Module color
func (m *Module) Set(value Pixel) {
	m.value = value
	m.isUsed = true
}

type QRCode struct {
	data       []byte
	version    int
	correction CodeLevel
	pattern    int

	canvas [][]Module
	size   int

	alignments []int
}

func NewQRCode(e *Encoder, data []byte) *QRCode {
	e.version += 1               // FIXME: remove after versioning conflicts
	canvasSize := versionOneSize // base canvas size for e.version == 1
	alignments := alignmentPatterns[e.version]

	if e.version != 1 {
		aLen := len(alignments)
		canvasSize = alignments[aLen-1] + versionSizePadding
	}

	var canvas [][]Module = make([][]Module, canvasSize)
	for i := range canvas {
		canvas[i] = make([]Module, canvasSize)
	}

	for i, row := range canvas {
		for j := range row {
			canvas[i][j].value = true
			canvas[i][j].isUsed = false
		}
	}

	return &QRCode{
		data:       data,
		version:    e.version,
		pattern:    0,
		correction: e.level,
		canvas:     canvas,
		size:       canvasSize,
		alignments: alignments,
	}
}

func (g *QRCode) String() string {
	var buf bytes.Buffer

	buf.WriteByte('{')
	fmt.Fprintf(&buf, "\nversion: %v", g.version)
	fmt.Fprintf(&buf, "\nerror correction: %v", g.correction)
	fmt.Fprintf(&buf, "\nmask pattern: %v", g.pattern)
	fmt.Fprintf(&buf, "\nalignments: %v", g.alignments)
	fmt.Fprintf(&buf, "\ndata: ")

	for _, row := range g.canvas {
		buf.WriteString("\n\t\t")
		// fmt.Fprintf(&buf, "\n%d:\t", i)
		for _, v := range row {
			fmt.Fprintf(&buf, "%v", v.value)
		}
	}

	buf.WriteString("\n}")

	return buf.String()
}

func (g *QRCode) MakeLayout() {
	g.placeSearchPatterns()
	g.placeAlignments()
	g.placeSync()

	if g.version > 6 {
		g.placeVersion()
	}

	g.placeMask()

	// g.Write(g.data)

}

func (g *QRCode) Write(bytes []byte) (int, error) {
	var n int
	xl, xr := g.size-2, g.size-1
	upwards := true

	mask := maskFunctions[g.pattern]
	nextBit := g.bitsGenerator() // convert encoded data to bit flow

	for xl >= 0 {
		if xr == 6 { // skip vertical synchronization line
			xl, xr = xl-1, xr-1
		}

		y, border := g.size-1, -1
		if !upwards {
			y, border = 0, g.size
		}

		for y != border {
			if !g.canvas[y][xr].isUsed {
				bit := !nextBit()

				if mask(xr, y) == 0 {
					bit = !bit
				}

				g.canvas[y][xr].Set(Pixel(bit))

			}

			if !g.canvas[y][xl].isUsed {
				bit := !nextBit()

				if mask(xl, y) == 0 {
					bit = !bit
				}

				g.canvas[y][xl].Set(Pixel(bit))
			}

			if upwards {
				y--
			} else {
				y++
			}
		}

		xl, xr = xl-2, xr-2
		upwards = !upwards
	}

	return n, nil
}

func (g *QRCode) bitsGenerator() func() bool {
	dataBits := make([]bool, 0, len(g.data)*8)
	for _, b := range g.data {
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

func (g *QRCode) placeSearchPatterns() {
	g.placePattern(0, 0, &searchPatternTL)                            // Top left corner
	g.placePattern(0, g.size-searchPatternTR.xSize, &searchPatternTR) // Top right
	g.placePattern(g.size-searchPatternBL.ySize, 0, &searchPatternBL) // Bottom left
}

func (g *QRCode) placeAlignments() {
	perms := algorithms.GeneratePermutations(g.alignments)

	for _, loc := range perms {
		g.placePattern(loc[0]-(alignPattern.xSize/2), loc[1]-(alignPattern.ySize/2), &alignPattern)
	}
}

func (g *QRCode) placeSync() {
	var locX, locY int
	var i int
	syncPixels := [2]Pixel{bl, wh}

	// Vertical sync border
	for i, locX, locY = 0, 6, 8; locY < g.size-syncPadding; i, locY = (i+1)%2, locY+1 {
		if !g.canvas[locY][locX].isUsed {
			g.canvas[locY][locX].Set(syncPixels[i])
		}
	}

	// Horizontal sync border
	for i, locX, locY = 0, 8, 6; locX < g.size-syncPadding; i, locX = (i+1)%2, locX+1 {
		if !g.canvas[locY][locX].isUsed {
			g.canvas[locY][locX].Set(syncPixels[i])
		}
	}

}

func (g *QRCode) placeVersion() {
	var locX, locY int = 0, g.size - versionPadding
	var x, y int

	versionBits := versionCodes[g.version]

	for y_offset, b := range versionBits {
		bits := algorithms.ToBoolArray(b)

		for x_offset, bit := range bits[2:] {
			x, y = locX+x_offset, locY+y_offset

			g.canvas[y][x].Set(Pixel(!bit)) // Bottom left code
			g.canvas[x][y].Set(Pixel(!bit)) // Top right code
		}
	}

}

func (g *QRCode) placeMask() {
	code := maskCodes[g.correction][g.pattern]

	codeBits := make([]bool, 0, 15)
	msb := algorithms.ToBoolArray(byte(code >> 8))
	lsb := algorithms.ToBoolArray(byte(code))

	codeBits = append(codeBits, msb[1:]...)
	codeBits = append(codeBits, lsb[:]...)

	// Bottom left + Top right
	i := 0
	for x, y := 8, g.size-1; y > g.size-8; y = y - 1 {
		g.canvas[y][x].Set(Pixel(!codeBits[i]))
		i++
	}

	g.canvas[g.size-8][8].Set(false) // This module is always black

	for x, y := g.size-8, 8; x < g.size; x = x + 1 {
		g.canvas[y][x].Set(Pixel(!codeBits[i]))
		i++
	}

	// Top left
	i = 0
	for x, y := 0, 8; x < 9; x++ {
		if !g.canvas[y][x].isUsed {
			g.canvas[y][x].Set(Pixel(!codeBits[i]))
			i++
		}
	}

	for x, y := 8, 7; y > -1; y-- {
		if !g.canvas[y][x].isUsed {
			g.canvas[y][x].Set(Pixel(!codeBits[i]))
			i++
		}
	}

}

func (g *QRCode) isUnused(startX, startY, endX, endY int) bool {

	// false, if arguments are out of canvas bounds
	if startX >= g.size || startY >= g.size || endX > g.size || endY > g.size {
		return false
	}

	for i := startX; i < endX; i++ {
		for j := startY; j < endY; j++ {
			if g.canvas[i][j].isUsed {
				return false
			}
		}
	}

	return true
}

func (g *QRCode) placePattern(locX, locY int, p *Pattern) {
	pxLen, pyLen := p.xSize, p.ySize

	if !g.isUnused(locX, locY, locX+pxLen, locY+pyLen) {
		return
	}

	for i, pi := locX, 0; i < g.size && pi < pxLen; i, pi = i+1, pi+1 {
		for j, pj := locY, 0; j < g.size && pj < pyLen; j, pj = j+1, pj+1 {
			g.canvas[i][j].Set(p.data[pi][pj])
		}
	}
}
