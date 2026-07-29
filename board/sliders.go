package board

// slider holds everything needed to answer attack queries for one kind of sliding piece.
type slider struct {
	magicEntries [64]magicEntry
}

// magicEntry holds the magic bitboard lookup data for a single square: the relevant occupancy
// mask, the magic number that hashes an occupancy to a table index, the shift that reduces that
// hash to the low bits, and the table of precomputed attacks indexed by it.
type magicEntry struct {
	magic   uint64
	mask    Bitboard
	shift   int
	attacks []Bitboard
}

var (
	// rookSlider and bishopSlider are the magic lookup tables, built once at startup from the
	// generated magic numbers in magics_gen.go.
	rookSlider   = buildSlider(rookDirections, rookMagicNumbers)
	bishopSlider = buildSlider(bishopDirections, bishopMagicNumbers)
)

// buildSlider builds the complete lookup table for one kind of sliding piece, using the
// pre-generated magic numbers for each square.
func buildSlider(directions []directionFunc, magicNumbers [64]uint64) slider {
	return slider{
		magicEntries: buildSliderMagicEntries(directions, magicNumbers),
	}
}

// buildSliderMagicEntries builds a magicEntry for every square, filling each attack table by
// running every relevant occupancy through the square's magic to find the slot it will be read
// from at lookup time.
func buildSliderMagicEntries(directions []directionFunc, magicNumbers [64]uint64) [64]magicEntry {
	var entries [64]magicEntry
	for sq := A1; sq <= H8; sq++ {
		mask := computeRelevantOccupancyMask(sq, directions)
		n := mask.Count()
		shift := 64 - n
		attacks := make([]Bitboard, 1<<n)
		pairs := computeOccupancyAttacks(sq, directions)
		for _, pair := range pairs {
			index := int((uint64(pair.occupancy) * magicNumbers[sq]) >> shift)
			if existing := attacks[index]; existing > 0 && existing != pair.attacks {
				panic("magic number collision detected")
			}
			attacks[index] = pair.attacks
		}
		entries[sq] = magicEntry{
			magic:   magicNumbers[sq],
			mask:    mask,
			shift:   shift,
			attacks: attacks,
		}
	}
	return entries
}

// rookLookup returns the squares available to move to by a rook on the given square,
// given the occupancy of the whole board. Blockers are included in the result regardless
// of colour - callers wanting legal destinations must remove their own pieces.
func rookLookup(sq Square, occupancy Bitboard) Bitboard {
	entry := &rookSlider.magicEntries[sq]
	index := int((uint64(occupancy&entry.mask) * entry.magic) >> entry.shift)
	return entry.attacks[index]
}

// bishopLookup returns the squares available to move to by a bishop on the given square,
// given the occupancy of the whole board. Blockers are included in the result regardless
// of colour - callers wanting legal destinations must remove their own pieces.
func bishopLookup(sq Square, occupancy Bitboard) Bitboard {
	entry := &bishopSlider.magicEntries[sq]
	index := int((uint64(occupancy&entry.mask) * entry.magic) >> entry.shift)
	return entry.attacks[index]
}
