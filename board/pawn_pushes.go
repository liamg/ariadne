package board

// findPawnPushes returns the bitboards of single and double pawn pushes for the given colour, pawns and occupancy.
// pawns is a bitboard of the pawns of the given colour, and bothOccupancy is a bitboard of all pieces on the board.
func findPawnPushes(colour Colour, pawns Bitboard, bothOccupancy Bitboard) (singles, doubles Bitboard) {
	switch colour {
	case White:
		singles = pawns.North() &^ bothOccupancy
		doubles = (singles.North() &^ bothOccupancy) & Rank4Mask
	case Black:
		singles = pawns.South() &^ bothOccupancy
		doubles = (singles.South() &^ bothOccupancy) & Rank5Mask
	}
	return singles, doubles
}
