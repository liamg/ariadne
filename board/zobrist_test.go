package board

import (
	"strings"
	"testing"
)

func TestZobristHash(t *testing.T) {
	var recurse func(pos *Position, depth int, hash uint64, preMoves []Move)
	recurse = func(pos *Position, depth int, hash uint64, preMoves []Move) {
		if depth == 0 {
			return
		}
		moves := pos.GenerateLegalMoves()
		for _, move := range moves {
			undo := pos.MakeMove(move)
			if GenerateZobristHash(pos) != pos.ZobristHash() {
				t.Errorf("Zobrist hash mismatch after making move %s: got %x, want %x", move.String(), pos.ZobristHash(), GenerateZobristHash(pos))
			}

			postHash := pos.ZobristHash()
			if postHash == hash {
				ms := make([]string, 0, len(preMoves)+1)
				for _, m := range preMoves {
					ms = append(ms, m.String())
				}
				ms = append(ms, move.String())
				t.Errorf("Zobrist hash did not change after making moves: %s", strings.Join(ms, " "))
			}
			m := make([]Move, len(preMoves)+1)
			copy(m, preMoves)
			m[len(preMoves)] = move
			recurse(pos, depth-1, postHash, m)
			pos.UnmakeMove(undo)
			if hash != pos.ZobristHash() {
				ms := make([]string, 0, len(preMoves)+1)
				for _, m := range preMoves {
					ms = append(ms, m.String())
				}
				ms = append(ms, move.String())
				t.Errorf("Zobrist hash mismatch after unmaking moves %s: got %x, want %x", strings.Join(ms, " "), pos.ZobristHash(), hash)
			}
		}
	}

	for _, test := range perftCases {
		t.Run(test.name, func(t *testing.T) {
			pos, err := ParseFEN(test.fen)
			if err != nil {
				t.Fatalf("failed to parse FEN: %v", err)
			}
			recurse(pos, 4, pos.ZobristHash(), nil)
		})
	}
}

func TestZobristTransposition(t *testing.T) {
	// 1.Nf3 d5 2.d4 versus 1.d4 d5 2.Nf3
	pos1, _ := ParseFEN(FenStarting)
	pos1.MakeMove(NewMove(G1, F3, QuietMove))
	pos1.MakeMove(NewMove(D7, D5, QuietMove))
	pos1.MakeMove(NewMove(D2, D4, QuietMove))

	pos2, _ := ParseFEN(FenStarting)
	pos2.MakeMove(NewMove(D2, D4, QuietMove))
	pos2.MakeMove(NewMove(D7, D5, QuietMove))
	pos2.MakeMove(NewMove(G1, F3, QuietMove))

	hash1 := pos1.ZobristHash()
	hash2 := pos2.ZobristHash()

	if hash1 != hash2 {
		t.Errorf("Zobrist hashes do not match: %x vs %x", hash1, hash2)
	}
}

func TestZobristSideToMove(t *testing.T) {
	pos, _ := ParseFEN(FenStarting)
	hash1 := pos.ZobristHash()
	pos.sideToMove = pos.sideToMove.Opposite()
	hash2 := GenerateZobristHash(pos)

	if hash1 == hash2 {
		t.Errorf("Zobrist hashes should not match after changing side to move: %x vs %x", hash1, hash2)
	}
}

func TestZobristOnRoundtrip(t *testing.T) {
	pos, _ := ParseFEN(FenKiwiPete)
	hash1 := pos.ZobristHash()
	fen := GenerateFEN(pos)
	pos2, _ := ParseFEN(fen)
	hash2 := pos2.ZobristHash()

	if hash1 != hash2 {
		t.Errorf("Zobrist hashes do not match after FEN roundtrip: %x vs %x", hash1, hash2)
	}
}

func BenchmarkZobrist(b *testing.B) {
	pos, _ := ParseFEN(FenKiwiPete)

	for b.Loop() {
		_ = pos.ZobristHash()
	}
}
