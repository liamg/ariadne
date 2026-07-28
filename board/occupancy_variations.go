package board

// computeOccupancyVariations returns every possible arrangement of occupied squares within the given mask.
// There are 2^n of these, where n is the number of bits set in the mask, and the variation at index i is the one produced by scatter(i, mask).
func computeOccupancyVariations(relevantOccupancy Bitboard) []Bitboard {
	numRelevantBits := relevantOccupancy.Count()
	output := make([]Bitboard, 1<<numRelevantBits)
	for i := range len(output) {
		output[i] = scatter(i, relevantOccupancy)
	}
	return output
}

// occupancyAttacks pairs an occupancy with the attacks a sliding piece generates when that occupancy is on the board.
type occupancyAttacks struct {
	occupancy Bitboard
	attacks   Bitboard
}

// computeOccupancyAttacks returns every relevant occupancy for the given square and directions, each paired with the true attacks found by slowlyFindRayAttacks.
// The results are index-aligned with computeOccupancyVariations. This is the ground truth used to generate and verify magic numbers.
func computeOccupancyAttacks(sq Square, dirs []directionFunc) []occupancyAttacks {
	mask := computeRelevantOccupancyMask(sq, dirs)
	occupancyVariations := computeOccupancyVariations(mask)
	attacks := make([]occupancyAttacks, len(occupancyVariations))
	for i, occupancy := range occupancyVariations {
		attacks[i] = occupancyAttacks{
			occupancy: occupancy,
			attacks:   slowlyFindRayAttacks(sq, occupancy, dirs),
		}
	}
	return attacks
}
