package qr

import (
	"github.com/psxzz/go-qr/pkg/algorithms"
)

const (
	baseScorePenalty1 = 3
	baseScorePenalty2 = 3
	baseScorePenalty3 = 40
)

func penalty1(canvas [][]Module) int {
	var score int

	// rows
	for _, row := range canvas {
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
	for i := 0; i < len(canvas); i++ {
		prev := canvas[0][i].value
		count := 1

		for j := 1; j < len(canvas); j++ {
			if canvas[j][i].value != prev {
				prev, count = canvas[j][i].value, 1
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

func penalty2(canvas [][]Module) int {
	var blocks int

	for i := 0; i < len(canvas)-1; i++ {
		for j := 0; j < len(canvas)-1; j++ {
			curr := canvas[i][j].value

			adjacent := canvas[i][j+1].value &&
				canvas[i+1][j].value &&
				canvas[i+1][j+1].value

			if curr == (curr && adjacent) {
				blocks++
			}
		}
	}

	return baseScorePenalty2 * blocks
}

// func penalty3(canvas [][]Module) int {
// 	var patternsCount int

// 	for _, row := range canvas {

// 	}

// 	return patternsCount
// }

func penalty4(canvas [][]Module) int {
	var score int
	total := len(canvas) * len(canvas)
	black := 0

	for _, row := range canvas {
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
