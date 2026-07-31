package search

import (
	"testing"

	"github.com/liamg/chess/board"
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
			// depth 2 only - at depth 3 white can promote next move instead,
			// so the king moves tie with a7a8q on score and the winner is
			// decided by move ordering
			name:         "promotes to queen",
			fen:          "8/P6k/8/8/8/8/8/7K w - - 0 1",
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
				t.Errorf("expected searcher to explore fewer nodes than minimax searcher, got %d >= %d", searcher.state.NodeCount, minimaxSearcher.state.NodeCount)
			}

			if result.BestMove.String() != tt.expectedMove {
				t.Errorf("expected best move %s, got %s", tt.expectedMove, result.BestMove.String())
			}
		})
	}
}

func BenchmarkSearch(b *testing.B) {
	pos, err := board.ParseFEN(board.FenKiwiPete)
	if err != nil {
		b.Fatalf("failed to parse FEN: %v", err)
	}

	var nc uint64
	for b.Loop() {
		b.StopTimer()
		searcher := New()
		b.StartTimer()
		result := searcher.Search(b.Context(), pos, Limits{Depth: 6})
		nc += result.NodeCount
	}
	b.ReportMetric(float64(nc)/b.Elapsed().Seconds(), "nodes/s")
	b.ReportMetric(float64(nc)/float64(b.N), "nodes/op")
}
