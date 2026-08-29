package search

import "math"

var lmrTable [MaxPly + 1][MaxMoves + 1]int

func init() {
	for depth := 0; depth <= MaxPly; depth++ {
		for move := 0; move <= MaxMoves; move++ {
			// NOTE: needs SPSA for proper values, but this is a good enough approximation for now
			lmrTable[depth][move] = int(0.75 + math.Log(float64(depth))*math.Log(float64(move))/2.25)
		}
	}
}
