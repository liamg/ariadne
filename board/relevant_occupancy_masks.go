package board

var (
	// these masks list all squares that could possibly block a sliding piece's attack originating at a given square from proceeding elsewhere
	rookRelevantOccupancy   = computeRelevantOccupancyMasks(rookDirections)
	bishopRelevantOccupancy = computeRelevantOccupancyMasks(bishopDirections)
)

func computeRelevantOccupancyMasks(directions []directionFunc) [64]Bitboard {
	var masks [64]Bitboard
	for sq := A1; sq <= H8; sq++ {
		masks[sq] = computeRelevantOccupancyMask(sq, directions)
	}
	return masks
}

func computeRelevantOccupancyMask(sq Square, directions []directionFunc) Bitboard {
	output := EmptyBitboard
	for _, dir := range directions {
		current := sq.Bitboard()
		for {
			current = dir(current)
			if current == EmptyBitboard || dir(current) == EmptyBitboard {
				break
			}
			output |= current
		}
	}
	return output
}
