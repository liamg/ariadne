package board

// kingAttacks is a precomputed table of king attacks for each square on the board.
var kingAttacks = computeKingAttacks()

func computeKingAttacks() [64]Bitboard {
	var attacks [64]Bitboard
	for square := A1; square <= H8; square++ {
		bb := square.Bitboard()
		attacks[square] = bb.North() | bb.NorthEast() | bb.East() | bb.SouthEast() | bb.South() | bb.SouthWest() | bb.West() | bb.NorthWest()
	}
	return attacks
}

// knightAttacks is a precomputed table of knight attacks for each square on the board.
var knightAttacks = computeKnightAttacks()

func computeKnightAttacks() [64]Bitboard {
	var attacks [64]Bitboard
	for square := A1; square <= H8; square++ {
		bb := square.Bitboard()
		attacks[square] = bb.North().North().East() |
			bb.North().North().West() |
			bb.North().East().East() |
			bb.North().West().West() |
			bb.South().South().East() |
			bb.South().South().West() |
			bb.South().East().East() |
			bb.South().West().West()
	}
	return attacks
}
