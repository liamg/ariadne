package board

import "math/rand/v2"

var (
	zobristTable          [64][16]uint64
	zobristBlackToMove    uint64
	zobristCastlingRights [16]uint64
	zobristEnPassant      [65]uint64
)

func init() {
	rnd := rand.New(rand.NewPCG(1, 2))
	for square := range 64 {
		for piece := 1; piece < 16; piece++ { // NOTE: index 0 intentionally left zero
			zobristTable[square][piece] = rnd.Uint64()
		}
	}
	for ep := range 64 { // NOTE: index 64 intentionally left zero
		zobristEnPassant[ep] = rnd.Uint64()
	}
	for cr := 1; cr < 16; cr++ { // NOTE: index 0 intentionally left zero
		zobristCastlingRights[cr] = rnd.Uint64()
	}
	zobristBlackToMove = rnd.Uint64()
}

func GenerateZobristHash(p *Position) uint64 {
	hash := uint64(0)
	if p.sideToMove == Black {
		hash = zobristBlackToMove
	}

	hash ^= zobristCastlingRights[p.state.castlingRights]
	// NoSquare holds a zero key, so this is a no-op when there's no ep square
	hash ^= zobristEnPassant[p.state.enPassantSquare]

	occ := p.Occupancy()
	var sq Square

	for occ != 0 {
		sq, occ = occ.PopSquare()
		piece := p.PieceAt(sq)
		hash ^= zobristTable[sq][piece]
	}

	return hash
}
