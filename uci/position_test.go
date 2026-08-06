package uci

import (
	"testing"
	"time"

	"github.com/liamg/ariadne/board"
)

func TestPosition(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedFen string
	}{
		{
			name:        "starting position",
			input:       "position startpos",
			expectedFen: board.FenStarting,
		},
		{
			name:        "starting position with moves",
			input:       "position startpos moves e2e4 e7e5 g1f3",
			expectedFen: "rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
		},
		{
			name:        "FEN position",
			input:       "position fen rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
			expectedFen: "rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2",
		},
		{
			name:        "FEN position with moves",
			input:       "position fen rnbqkbnr/pppp1ppp/8/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R b KQkq - 1 2 moves g8f6",
			expectedFen: "rnbqkb1r/pppp1ppp/5n2/4p3/4P3/5N2/PPPP1PPP/RNBQKB1R w KQkq - 2 3",
		},
		{
			name:        "white kingside castle",
			input:       "position fen r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1 moves e1g1",
			expectedFen: "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R4RK1 b kq - 1 1",
		},
		{
			name:        "white queenside castle",
			input:       "position fen r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQkq - 0 1 moves e1c1",
			expectedFen: "r3k2r/pppppppp/8/8/8/8/PPPPPPPP/2KR3R b kq - 1 1",
		},
		{
			name:        "black kingside castle",
			input:       "position fen r3k2r/pppppppp/8/8/8/8/PPPPPPPP/R3K2R b KQkq - 0 1 moves e8g8",
			expectedFen: "r4rk1/pppppppp/8/8/8/8/PPPPPPPP/R3K2R w KQ - 1 2",
		},
		{
			// rank 5 must be empty: the captured pawn on d5 is gone
			name:        "en passant capture",
			input:       "position startpos moves e2e4 a7a6 e4e5 d7d5 e5d6",
			expectedFen: "rnbqkbnr/1pp1pppp/p2P4/8/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 3",
		},
		{
			name:        "promotion to queen",
			input:       "position fen 4k3/P7/8/8/8/8/8/4K3 w - - 0 1 moves a7a8q",
			expectedFen: "Q3k3/8/8/8/8/8/8/4K3 b - - 0 1",
		},
		{
			name:        "promotion to rook",
			input:       "position fen 4k3/P7/8/8/8/8/8/4K3 w - - 0 1 moves a7a8r",
			expectedFen: "R3k3/8/8/8/8/8/8/4K3 b - - 0 1",
		},
		{
			name:        "promotion to bishop",
			input:       "position fen 4k3/P7/8/8/8/8/8/4K3 w - - 0 1 moves a7a8b",
			expectedFen: "B3k3/8/8/8/8/8/8/4K3 b - - 0 1",
		},
		{
			name:        "promotion to knight",
			input:       "position fen 4k3/P7/8/8/8/8/8/4K3 w - - 0 1 moves a7a8n",
			expectedFen: "N3k3/8/8/8/8/8/8/4K3 b - - 0 1",
		},
		{
			// the whole line is rejected, leaving the position untouched
			name:        "invalid move token",
			input:       "position startpos moves e2e4 zzzz",
			expectedFen: board.FenStarting,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHarness(t)
			defer h.Close()
			u := New(h, h)
			go func() { _ = u.Run(t.Context()) }()

			h.Send(tt.input)
			h.WaitUntilReady(time.Second)

			if actual := board.GenerateFEN(u.eng.position); actual != tt.expectedFen {
				t.Errorf("expected FEN %s, got %s", tt.expectedFen, actual)
			}
		})
	}
}

// TestPositionRepetition guards the "replay every move" property: repetition is
// detected from the zobrist history built up by MakeMove, so a position command
// that shortcut to a final board state would leave this undetectable.
func TestPositionRepetition(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedDraw bool
	}{
		{
			name:         "no moves",
			input:        "position startpos",
			expectedDraw: false,
		},
		{
			name:         "no repetition",
			input:        "position startpos moves e2e4 e7e5",
			expectedDraw: false,
		},
		{
			// one ply short of returning to the starting position
			name:         "partial repetition",
			input:        "position startpos moves g1f3 g8f6 f3g1",
			expectedDraw: false,
		},
		{
			// both knights out and back: the starting position occurs twice
			name:         "repetition",
			input:        "position startpos moves g1f3 g8f6 f3g1 f6g8",
			expectedDraw: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newTestHarness(t)
			defer h.Close()
			u := New(h, h)
			go func() { _ = u.Run(t.Context()) }()

			h.Send(tt.input)
			h.WaitUntilReady(time.Second)

			if actual := u.eng.position.IsDrawByRepetition(); actual != tt.expectedDraw {
				t.Errorf("expected IsDrawByRepetition %v, got %v", tt.expectedDraw, actual)
			}
		})
	}
}
