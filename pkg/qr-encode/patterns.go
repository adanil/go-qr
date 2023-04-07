package qr_encode

const (
	bl = false
	wh = true
)

type Pattern struct {
	data  [][]Pixel
	xSize int
	ySize int
}

var (
	searchPatternTL = Pattern{
		data: [][]Pixel{
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
		data: [][]Pixel{
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
		data: [][]Pixel{
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
		data: [][]Pixel{
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
