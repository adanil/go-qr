package qr

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
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
