package search

import (
	"context"
	"testing"
	"time"

	"github.com/liamg/ariadne/board"
	"github.com/liamg/ariadne/eval"
)

func TestSearch(t *testing.T) {
	tests := []struct {
		name         string
		fen          string
		expectedMove string
		depth        int
	}{
		{
			name:         "wins hanging queen",
			fen:          "4k3/8/8/8/3q4/8/8/3RK3 w - - 0 1",
			depth:        1,
			expectedMove: "d1d4",
		},
		{
			name:         "wins hanging queen deeper",
			fen:          "4k3/8/8/8/3q4/8/8/3RK3 w - - 0 1",
			depth:        3,
			expectedMove: "d1d4",
		},
		{
			name:         "white mates in 1",
			fen:          "6k1/5ppp/8/8/8/8/8/R5K1 w - - 0 1",
			depth:        2,
			expectedMove: "a1a8",
		},
		{
			// still found when there is spare depth to waste
			name:         "white mates in 1 deeper",
			fen:          "6k1/5ppp/8/8/8/8/8/R5K1 w - - 0 1",
			depth:        3,
			expectedMove: "a1a8",
		},
		{
			name:         "black wins hanging queen",
			fen:          "3rk3/8/8/3Q4/8/8/8/4K3 b - - 0 1",
			depth:        1,
			expectedMove: "d8d5",
		},
		{
			name:         "black mates in 1",
			fen:          "r5k1/8/8/8/8/8/5PPP/6K1 b - - 0 1",
			depth:        2,
			expectedMove: "a8a1",
		},
		{
			// black king can reach b7 and cover a8, so delaying costs the pawn.
			// with the king far away every first move ties - qsearch promotes
			// at the leaf anyway - and the king move wins on position
			name:         "promotes to queen",
			fen:          "8/P1k5/8/8/8/8/8/7K w - - 0 1",
			depth:        2,
			expectedMove: "a7a8q",
		},
		{
			name:         "captures en passant",
			fen:          "8/8/8/3pP3/8/8/8/K6k w - d6 0 1",
			depth:        1,
			expectedMove: "e5d6",
		},
		{
			name:         "king captures checking rook",
			fen:          "4k3/8/8/8/8/8/4r3/4K2R w K - 0 1",
			depth:        1,
			expectedMove: "e1e2",
		},
		{
			// no legal moves, so no move is ever recorded
			name:         "stalemate has no best move",
			fen:          "7k/5Q2/6K1/8/8/8/8/8 b - - 0 1",
			depth:        1,
			expectedMove: "0000",
		},
		{
			name:         "mate in 2",
			fen:          "7k/p7/5K2/8/8/8/8/R7 w - - 0 1",
			depth:        4,
			expectedMove: "f6f7",
		},
		{
			name:         "mate in 3",
			fen:          "6k1/1rr2ppp/8/8/8/8/8/R5K1 w - - 0 1",
			depth:        3,
			expectedMove: "a1a8",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := board.ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("failed to parse FEN: %v", err)
			}

			searcher := New()
			result := searcher.Search(t.Context(), pos, Limits{Depth: tt.depth})

			minimaxSearcher := New()
			minimaxResult := minimaxSearcher.minimax(pos, tt.depth, 0)
			if pos.SideToMove() == board.Black {
				minimaxResult = -minimaxResult
			}

			if result.Score != minimaxResult {
				t.Errorf("expected score %d, got %d", minimaxResult, result.Score)
			}

			if searcher.state.NodeCount > minimaxSearcher.state.NodeCount {
				t.Errorf("expected searcher to explore fewer nodes than minimax searcher, got %d > %d", searcher.state.NodeCount, minimaxSearcher.state.NodeCount)
			}

			if result.BestMove.String() != tt.expectedMove {
				t.Errorf("expected best move %s, got %s", tt.expectedMove, result.BestMove.String())
			}

			// reuse the same searcher and ensure result is the same
			result2 := searcher.Search(t.Context(), pos, Limits{Depth: tt.depth})
			if result.BestMove != result2.BestMove {
				t.Errorf("expected best move to be the same on reuse, got %s != %s", result.BestMove.String(), result2.BestMove.String())
			}
			if result.Score != result2.Score {
				t.Errorf("expected score to be the same on reuse, got %d != %d", result.Score, result2.Score)
			}
			if result2.NodeCount > result.NodeCount {
				t.Errorf("expected result to be the same or less on reuse, got %v != %v", result, result2)
			}
		})
	}
}

func TestSearchWithExpiredContext(t *testing.T) {
	pos, err := board.ParseFEN("4k3/8/8/8/3q4/8/8/3RK3 w - - 0 1") // hanging black queen
	if err != nil {
		t.Fatalf("failed to parse FEN: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	searcher := New()
	result := searcher.Search(ctx, pos, Limits{Depth: 6})

	if result.BestMove == board.NullMove {
		t.Errorf("expected best move to be set even with expired context, got %s", result.BestMove.String())
	}
	// loose floor, not an exact value - the point is the queen was won
	if result.Score < 400 {
		t.Errorf("expected score to be >= 400 for a hanging queen, got %d", result.Score)
	}
}

func TestCancelledSearchDoesNotCorruptPosition(t *testing.T) {
	inputFEN := "4k3/8/8/8/3q4/8/8/3RK3 w - - 0 1"
	pos, err := board.ParseFEN(inputFEN) // hanging black queen
	if err != nil {
		t.Fatalf("failed to parse FEN: %v", err)
	}
	before := board.GenerateZobristHash(pos)

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()

	searcher := New()
	result := searcher.Search(ctx, pos, Limits{Depth: 20})

	if result.BestMove == board.NullMove {
		t.Errorf("expected best move to be set even with expired context, got %s", result.BestMove.String())
	}
	// loose floor, not an exact value - the point is the queen was won
	if result.Score < 400 {
		t.Errorf("expected score to be >= 400 for a hanging queen, got %d", result.Score)
	}

	// ensure the position is still valid after a cancelled search
	after := board.GenerateZobristHash(pos)
	if before != after {
		t.Errorf("expected position to be unchanged after cancelled search, got before: %d, after: %d", before, after)
	}
	if board.GenerateFEN(pos) != inputFEN {
		t.Errorf("expected FEN to be unchanged after cancelled search, got before: %s, after: %s", board.FenKiwiPete, board.GenerateFEN(pos))
	}
}

func TestSearchWithPlyRelativeScore(t *testing.T) {
	tests := []struct {
		name          string
		fen           string
		expectedMove  string
		expectedScore eval.Score
		depth         int
	}{
		{
			name:          "white mates in 2",
			fen:           "6k1/1r3ppp/8/8/8/8/8/R5K1 w - - 0 1",
			depth:         4,
			expectedMove:  "a1a8",
			expectedScore: eval.Mate - 3,
		},
		{
			name:          "black mates in 2",
			fen:           "r5k1/8/8/8/8/8/1R3PPP/6K1 b - - 0 1",
			depth:         4,
			expectedMove:  "a8a1",
			expectedScore: eval.Mate - 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := board.ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("failed to parse FEN: %v", err)
			}

			searcher := New()
			result := searcher.Search(t.Context(), pos, Limits{Depth: tt.depth})

			if result.Score != tt.expectedScore {
				t.Errorf("expected score %d, got %d", tt.expectedScore, result.Score)
			}
			if result.BestMove.String() != tt.expectedMove {
				t.Errorf("expected best move %s, got %s", tt.expectedMove, result.BestMove.String())
			}

			pos.MakeMove(result.BestMove)
			entry, ok := searcher.tt.probe(pos.ZobristHash(), 1)
			if !ok {
				t.Fatalf("expected transposition table to contain the position after making the best move")
			}
			if eval.Score(entry.Score) != (-tt.expectedScore) {
				t.Errorf("expected score after making best move to be %d (taking into account ply=1), got %d", (-tt.expectedScore), entry.Score)
			}
			entryPly5, ok := searcher.tt.probe(pos.ZobristHash(), 5)
			if !ok {
				t.Fatalf("expected transposition table to contain the position after making the best move")
			}
			if eval.Score(entryPly5.Score) != (-tt.expectedScore + 4) {
				t.Errorf("expected score after making best move to be %d (taking into account ply=5), got %d", (-tt.expectedScore + 4), entryPly5.Score)
			}
		})
	}
}

func BenchmarkSearch(b *testing.B) {
	pos, err := board.ParseFEN(board.FenKiwiPete)
	if err != nil {
		b.Fatalf("failed to parse FEN: %v", err)
	}

	searcher := New()
	_ = searcher.Search(b.Context(), pos, Limits{Depth: 5})
	b.ResetTimer()
	var nc uint64
	for b.Loop() {
		result := searcher.Search(b.Context(), pos, Limits{Depth: 5})
		nc += result.NodeCount
	}
	b.ReportMetric(float64(nc)/b.Elapsed().Seconds(), "nodes/s")
	b.ReportMetric(float64(nc)/float64(b.N), "nodes/op")
}
