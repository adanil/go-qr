package qr

const (
	wh = false
	bl = true
)

type Pattern struct {
	data  [][]bool
	xSize int
	ySize int
}

const (
	finderPatternSize    = 8
	alignmentPatternSize = 5
)

var (
	finderPatternTL = Pattern{
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

	finderPatternTR = Pattern{
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

	finderPatternBL = Pattern{
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

	alignmentPattern = Pattern{
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
