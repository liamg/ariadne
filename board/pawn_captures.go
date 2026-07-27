package board

// findPawnCaptures gives all capture squares for a given colour's pawns, given the enemy occupancy.
// It returns two bitboards: one for captures to the west and one for captures to the east.
func findPawnCaptures(colour Colour, pawns Bitboard, enemyOccupancy Bitboard) (west, east Bitboard) {
	if colour == White {
		west = pawns.NorthWest() & enemyOccupancy
		east = pawns.NorthEast() & enemyOccupancy
	} else {
		west = pawns.SouthWest() & enemyOccupancy
		east = pawns.SouthEast() & enemyOccupancy
	}
	return west, east
}
