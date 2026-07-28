package board

// findEnPassantAttackers returns a bitboard of pawns of the given colour than
// can legally capture en passant on the given en passant square.
func findEnPassantAttackers(colour Colour, pawns Bitboard, epSquare Square) Bitboard {
	if epSquare == NoSquare {
		return EmptyBitboard
	}
	return pawnAttacks[colour.Opposite()][epSquare] & pawns
}
