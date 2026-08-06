package search

import (
	"testing"

	"github.com/liamg/ariadne/board"
	"github.com/liamg/ariadne/eval"
)

// bands, not exact values, so eval retuning doesn't churn these.
// mate and draw scores are exact - min == max.
// the minimax comparison is the real check, and is eval independent
func TestNegamax(t *testing.T) {
	tests := []struct {
		name     string
		fen      string
		depth    int
		minScore eval.Score
		maxScore eval.Score
	}{
		{
			name:     "initial position depth 1",
			fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			depth:    1,
			minScore: -100,
			maxScore: 100,
		},
		{
			name:     "initial position depth 2",
			fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			depth:    2,
			minScore: -100,
			maxScore: 100,
		},
		{
			name:     "initial position depth 3",
			fen:      "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			depth:    3,
			minScore: -100,
			maxScore: 100,
		},
		{
			name:     "bare kings",
			fen:      "4k3/8/8/8/8/8/8/4K3 w - - 0 1",
			depth:    3,
			minScore: -100,
			maxScore: 100,
		},
		{
			// rook takes a hanging queen: roughly +500 either side of tuning
			name:     "white wins queen depth 1",
			fen:      "4k3/8/8/8/3q4/8/8/3RK3 w - - 0 1",
			depth:    1,
			minScore: 300,
			maxScore: 700,
		},
		{
			name:     "white wins queen depth 3",
			fen:      "4k3/8/8/8/3q4/8/8/3RK3 w - - 0 1",
			depth:    3,
			minScore: 300,
			maxScore: 700,
		},
		{
			name:     "black wins queen depth 1",
			fen:      "3rk3/8/8/3Q4/8/8/8/4K3 b - - 0 1",
			depth:    1,
			minScore: 300,
			maxScore: 700,
		},
		{
			name:     "black wins queen depth 3",
			fen:      "3rk3/8/8/3Q4/8/8/8/4K3 b - - 0 1",
			depth:    3,
			minScore: 300,
			maxScore: 700,
		},
		{
			name:     "mate in 1 visible at depth 1 with quiescence",
			fen:      "6k1/5ppp/8/8/8/8/8/R5K1 w - - 0 1",
			depth:    1,
			minScore: eval.Mate - 1,
			maxScore: eval.Mate - 1,
		},
		{
			name:     "white mates in 1",
			fen:      "6k1/5ppp/8/8/8/8/8/R5K1 w - - 0 1",
			depth:    2,
			minScore: eval.Mate - 1,
			maxScore: eval.Mate - 1,
		},
		{
			name:     "black mates in 1",
			fen:      "r5k1/8/8/8/8/8/5PPP/6K1 b - - 0 1",
			depth:    2,
			minScore: eval.Mate - 1,
			maxScore: eval.Mate - 1,
		},
		{
			name:     "white already mated",
			fen:      "6k1/8/8/8/8/8/5PPP/r5K1 w - - 0 1",
			depth:    1,
			minScore: -eval.Mate,
			maxScore: -eval.Mate,
		},
		{
			name:     "black already mated",
			fen:      "R5k1/5ppp/8/8/8/8/8/6K1 b - - 0 1",
			depth:    1,
			minScore: -eval.Mate,
			maxScore: -eval.Mate,
		},
		{
			name:     "stalemate",
			fen:      "7k/5Q2/6K1/8/8/8/8/8 b - - 0 1",
			depth:    1,
			minScore: eval.Draw,
			maxScore: eval.Draw,
		},
		{
			name:     "mate in 2",
			fen:      "7k/p7/5K2/8/8/8/8/R7 w - - 0 1",
			depth:    4,
			minScore: eval.Mate - 3,
			maxScore: eval.Mate - 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := board.ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("failed to parse FEN: %v", err)
			}

			searcher := New()
			minimaxSearcher := New()
			minimaxScore := minimaxSearcher.minimax(pos, tt.depth, 0)
			if pos.SideToMove() == board.Black {
				minimaxScore = -minimaxScore
			}
			score := searcher.negamax(pos, tt.depth, 0, -eval.Infinity, eval.Infinity)
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("expected score in [%v, %v], got %v", tt.minScore, tt.maxScore, score)
			}
			if score != minimaxScore {
				t.Errorf("expected score to match minimax at %v, got %v", minimaxScore, score)
			}
		})
	}
}
