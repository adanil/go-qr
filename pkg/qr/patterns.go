package qr

const (
	wh = false
	bl = true
)

type qrPattern struct {
	data  [][]bool
	xSize int
	ySize int
}

const (
	finderPatternSize    = 8
	alignmentPatternSize = 5
)

var (
	// Finder pattern with separator located in top left corner of QR code
	finderPatternTL = qrPattern{
		data: [][]bool{
			{bl, bl, bl, bl, bl, bl, bl, wh},
			{bl, wh, wh, wh, wh, wh, bl, wh},
			{bl, wh, bl, bl, bl, wh, bl, wh},
			{bl, wh, bl, bl, bl, wh, bl, wh},
			{bl, wh, bl, bl, bl, wh, bl, wh},
			{bl, wh, wh, wh, wh, wh, bl, wh},
			{bl, bl, bl, bl, bl, bl, bl, wh},
			{wh, wh, wh, wh, wh, wh, wh, wh},
		},
		xSize: finderPatternSize,
		ySize: finderPatternSize,
	}

	// Finder pattern with separator located in top right corner of QR code
	finderPatternTR = qrPattern{
		data: [][]bool{
			{wh, bl, bl, bl, bl, bl, bl, bl},
			{wh, bl, wh, wh, wh, wh, wh, bl},
			{wh, bl, wh, bl, bl, bl, wh, bl},
			{wh, bl, wh, bl, bl, bl, wh, bl},
			{wh, bl, wh, bl, bl, bl, wh, bl},
			{wh, bl, wh, wh, wh, wh, wh, bl},
			{wh, bl, bl, bl, bl, bl, bl, bl},
			{wh, wh, wh, wh, wh, wh, wh, wh},
		},
		xSize: finderPatternSize,
		ySize: finderPatternSize,
	}

	// Finder pattern with separator located in bottom left corner of QR code
	finderPatternBL = qrPattern{
		data: [][]bool{
			{wh, wh, wh, wh, wh, wh, wh, wh},
			{bl, bl, bl, bl, bl, bl, bl, wh},
			{bl, wh, wh, wh, wh, wh, bl, wh},
			{bl, wh, bl, bl, bl, wh, bl, wh},
			{bl, wh, bl, bl, bl, wh, bl, wh},
			{bl, wh, bl, bl, bl, wh, bl, wh},
			{bl, wh, wh, wh, wh, wh, bl, wh},
			{bl, bl, bl, bl, bl, bl, bl, wh},
		},
		xSize: finderPatternSize,
		ySize: finderPatternSize,
	}

	// General alignment pattern arranged according to the version
	alignmentPattern = qrPattern{
		data: [][]bool{
			{bl, bl, bl, bl, bl},
			{bl, wh, wh, wh, bl},
			{bl, wh, bl, wh, bl},
			{bl, wh, wh, wh, bl},
			{bl, bl, bl, bl, bl},
		},
		xSize: alignmentPatternSize,
		ySize: alignmentPatternSize,
	}
)
