package board

// isSquareAttacked checks if a square is attacked by any piece of the given colour.
func (p *Position) isSquareAttacked(sq Square, by Colour) bool {
	if knightAttackers := knightAttacks[sq] & p.byType[Knight] & p.byColour[by]; knightAttackers != 0 {
		return true
	}

	if pawnAttackers := pawnAttacks[by.Opposite()][sq] & p.byType[Pawn] & p.byColour[by]; pawnAttackers != 0 {
		return true
	}

	if kingAttackers := kingAttacks[sq] & p.byType[King] & p.byColour[by]; kingAttackers != 0 {
		return true
	}

	occ := p.Occupancy()

	bishopRays := bishopLookup(sq, occ)
	if bishopAttackers := bishopRays & p.byType[Bishop] & p.byColour[by]; bishopAttackers != 0 {
		return true
	}

	rookRays := rookLookup(sq, occ)
	if rookAttackers := rookRays & p.byType[Rook] & p.byColour[by]; rookAttackers != 0 {
		return true
	}

	if queenAttackers := (bishopRays | rookRays) & p.byType[Queen] & p.byColour[by]; queenAttackers != 0 {
		return true
	}

	return false
}
