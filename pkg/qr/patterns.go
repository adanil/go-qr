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

var (
	searchPatternTL = Pattern{
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
		xSize: 8,
		ySize: 8,
	}

	searchPatternTR = Pattern{
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
		xSize: 8,
		ySize: 8,
	}

	searchPatternBL = Pattern{
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
		xSize: 8,
		ySize: 8,
	}

	alignPattern = Pattern{
		data: [][]bool{
			{bl, bl, bl, bl, bl},
			{bl, wh, wh, wh, bl},
			{bl, wh, bl, wh, bl},
			{bl, wh, wh, wh, bl},
			{bl, bl, bl, bl, bl},
		},
		xSize: 5,
		ySize: 5,
	}
)
