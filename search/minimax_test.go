package search

import (
	"testing"

	"github.com/liamg/chess/board"
	"github.com/liamg/chess/eval"
)

func TestMinimax(t *testing.T) {
	tests := []struct {
		name     string
		fen      string
		depth    int
		expected eval.Score
	}{
		{
			name:     "initial position depth 1",
			fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			depth:    1,
			expected: eval.Score(0),
		},
		{
			name:     "initial position depth 2",
			fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			depth:    2,
			expected: eval.Score(0),
		},
		{
			name:     "initial position depth 3",
			fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			depth:    3,
			expected: eval.Score(0),
		},
		{
			name:     "bare kings",
			fen:      "4k3/8/8/8/8/8/8/4K3 w - - 0 1",
			depth:    3,
			expected: eval.Score(0),
		},
		{
			name:     "white wins queen depth 1",
			fen:      "4k3/8/8/8/3q4/8/8/3RK3 w - - 0 1",
			depth:    1,
			expected: eval.Score(500),
		},
		{
			name:     "white wins queen depth 3",
			fen:      "4k3/8/8/8/3q4/8/8/3RK3 w - - 0 1",
			depth:    3,
			expected: eval.Score(500),
		},
		{
			name:     "black wins queen depth 1",
			fen:      "3rk3/8/8/3Q4/8/8/8/4K3 b - - 0 1",
			depth:    1,
			expected: eval.Score(-500),
		},
		{
			name:     "black wins queen depth 3",
			fen:      "3rk3/8/8/3Q4/8/8/8/4K3 b - - 0 1",
			depth:    3,
			expected: eval.Score(-500),
		},
		{
			// depth 0 returns the static eval without a terminal check,
			// so a mate one ply away is invisible - this is just material
			name:     "mate in 1 invisible at depth 1",
			fen:      "6k1/5ppp/8/8/8/8/8/R5K1 w - - 0 1",
			depth:    1,
			expected: eval.Score(200),
		},
		{
			name:     "white mates in 1",
			fen:      "6k1/5ppp/8/8/8/8/8/R5K1 w - - 0 1",
			depth:    2,
			expected: eval.Mate - 1,
		},
		{
			name:     "black mates in 1",
			fen:      "r5k1/8/8/8/8/8/5PPP/6K1 b - - 0 1",
			depth:    2,
			expected: -(eval.Mate - 1),
		},
		{
			name:     "white already mated",
			fen:      "6k1/8/8/8/8/8/5PPP/r5K1 w - - 0 1",
			depth:    1,
			expected: -eval.Mate,
		},
		{
			name:     "black already mated",
			fen:      "R5k1/5ppp/8/8/8/8/8/6K1 b - - 0 1",
			depth:    1,
			expected: eval.Mate,
		},
		{
			name:     "stalemate",
			fen:      "7k/5Q2/6K1/8/8/8/8/8 b - - 0 1",
			depth:    1,
			expected: eval.Draw,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := board.ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("failed to parse FEN: %v", err)
			}

			searcher := New()
			result := searcher.minimax(pos, tt.depth, 0)
			if result != tt.expected {
				t.Errorf("expected score %v, got %v", tt.expected, result)
			}
		})
	}
}
