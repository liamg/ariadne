package board

type directionFunc func(Bitboard) Bitboard

var (
	bishopDirections = []directionFunc{Bitboard.NorthEast, Bitboard.NorthWest, Bitboard.SouthEast, Bitboard.SouthWest}
	rookDirections   = []directionFunc{Bitboard.North, Bitboard.South, Bitboard.East, Bitboard.West}
	queenDirections  = []directionFunc{Bitboard.North, Bitboard.South, Bitboard.East, Bitboard.West, Bitboard.NorthEast, Bitboard.NorthWest, Bitboard.SouthEast, Bitboard.SouthWest}
)

// slowlyFinRdRayAttacks returns a bitboard of all squares attacked by a sliding piece (bishop, rook, or queen) from the given square, considering the occupancy of the board. It uses a simple loop to find all squares in each direction until it hits an occupied square or the edge of the board.
// it is named so because it is a slow implementation of the ray attacks, and should not be used in performance critical code.
// it will mainly be used to generate magics and for testing "fast" variants of the code
func slowlyFindRayAttacks(sq Square, occupancy Bitboard, directions []directionFunc) Bitboard {
	output := EmptyBitboard
	for _, dir := range directions {
		current := sq.Bitboard()
		for {
			current = dir(current)
			if current == EmptyBitboard {
				break
			}
			output |= current
			if current&occupancy != EmptyBitboard {
				break
			}
		}
	}
	return output
}
