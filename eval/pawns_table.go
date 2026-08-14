package eval

import "github.com/liamg/ariadne/board"

const (
	pawnsTableLogSize uint64 = 12 // 2^12 = 4096
	pawnsTableSize    uint64 = 1 << pawnsTableLogSize
)

type pawnsTable struct {
	entries []pawnsTableEntry
}

func newPawnsTable() pawnsTable {
	e := make([]pawnsTableEntry, pawnsTableSize)
	for i := range e {
		e[i].kingSquare[0] = 0xFF
		e[i].kingSquare[1] = 0xFF
	}
	return pawnsTable{
		entries: e,
	}
}

type pawnsTableEntry struct {
	whitePawns   board.Bitboard
	blackPawns   board.Bitboard
	kingSquare   [2]byte
	shelter      [2]int16
	midGameScore Score
	endGameScore Score
}

func (p *pawnsTable) index(whitePawns, blackPawns board.Bitboard) uint64 {
	h := whitePawns*0x9E3779B97F4A7C15 ^ blackPawns*0xBF58476D1CE4E5B9
	return uint64(h) >> (64 - pawnsTableLogSize)
}

func (p *pawnsTable) probe(whitePawns, blackPawns board.Bitboard) (entry *pawnsTableEntry, ok bool) {
	index := p.index(whitePawns, blackPawns)
	e := &p.entries[index]
	return e, e.whitePawns == whitePawns && e.blackPawns == blackPawns
}
