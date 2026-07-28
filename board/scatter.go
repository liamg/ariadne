package board

// scatter takes the bits set in index, and disperses them amongst a new bitboard at the positions marked by the mask
// e.g. index = 0b101, mask = 0b101010000 would return 0b100010000
func scatter(index int, mask Bitboard) Bitboard {
	var result Bitboard
	var sq Square
	for i := 0; mask != EmptyBitboard; i++ {
		sq, mask = mask.PopSquare()
		if index&(1<<i) != 0 {
			result = result.Set(sq)
		}
	}
	return result
}
