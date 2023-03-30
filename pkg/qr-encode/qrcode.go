package qr_encode

import (
	"bytes"
	"fmt"

	"github.com/psxzz/go-qr/pkg/algorithms"
)

// TODO: Get rid of magic numbers and make it constant
const (
	versionPadding = 11
	syncPadding    = 7
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

	canvas [][]Module
	size   uint

	alignments []uint
}

func NewQRCode(e *Encoder, data []byte) *QRCode {
	canvasSize := uint(21) // base canvas size for e.version == 1
	alignments := alignmentPatterns[e.version]

	if e.version != 1 {
		aLen := len(alignments)
		canvasSize = alignments[aLen-1] + 7
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
	fmt.Fprintf(&buf, "\nalignments: %v", g.alignments)
	fmt.Fprintf(&buf, "\ndata: ")

	for i, row := range g.canvas {
		fmt.Fprintf(&buf, "\n%d:\t", i)
		for _, v := range row {
			fmt.Fprintf(&buf, "%v", v.value)
		}
	}

	buf.WriteByte('}')

	return buf.String()
}

func (g *QRCode) MakeLayout() {
	g.placeSearchPatterns()
	g.placeAlignments()
	g.placeSync()

	if g.version > 6 {
		g.placeVersion()
	}

	// g.Write(g.data)

	g.placeMask()
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
	var locX, locY uint
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
	var locX, locY uint = 0, g.size - versionPadding
	var x, y uint

	versionBinary := versionCodes[g.version]

	for y_offset, b := range versionBinary {
		bits := algorithms.ToBoolArray(b)

		for x_offset, bit := range bits {
			x, y = locX+uint(x_offset), locY+uint(y_offset)

			g.canvas[y][x].Set(Pixel(!bit)) // Bottom left code
			g.canvas[x][y].Set(Pixel(!bit)) // Top right code
		}
	}

}

func (g *QRCode) placeMask() {
	// maskPattern := 0
	// code := maskCodes[g.correction][maskPattern]
}

func (g *QRCode) isUnused(startX, startY, endX, endY uint) bool {

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

func (g *QRCode) placePattern(locX, locY uint, p *Pattern) {
	pxLen, pyLen := p.xSize, p.ySize

	if !g.isUnused(locX, locY, locX+pxLen, locY+pyLen) {
		return
	}

	for i, pi := locX, 0; i < g.size && uint(pi) < pxLen; i, pi = i+1, pi+1 {
		for j, pj := locY, 0; j < g.size && uint(pj) < pyLen; j, pj = j+1, pj+1 {
			g.canvas[i][j].Set(p.data[pi][pj])
		}
	}
}
