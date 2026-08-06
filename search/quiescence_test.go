package search

import (
	"testing"

	"github.com/liamg/ariadne/board"
	"github.com/liamg/ariadne/eval"
)

// relations against the static eval, so these survive any eval change
type want int

const (
	wantEqualsStatic want = iota // stood pat, or declined every capture
	wantBeatsStatic              // found something better than standing pat
	wantNotMate                  // ordinary score, not a mate score
	wantExact                    // mate scores only - eval independent
)

func TestQuiescence(t *testing.T) {
	tests := []struct {
		name       string
		fen        string
		ply        int
		want       want
		exactScore eval.Score
	}{
		{
			name: "simple position, no captures",
			fen:  "8/8/8/8/8/8/8/K6k w - - 0 1",
			want: wantEqualsStatic,
		},
		{
			// qsearch must find Qxd2 and resolve the capture rather than
			// evaluating a position with a hanging rook on it
			name: "resolves a free capture",
			fen:  "4k3/8/8/8/8/8/3r4/3QK3 w - - 0 1",
			want: wantBeatsStatic,
		},
		{
			// Qxd7 is met by Kxd7, so qsearch must decline it and return the
			// stand-pat score - fails if bestScore starts at -Infinity
			name: "stand-pat floor declines a losing capture",
			fen:  "4k3/3p4/8/8/8/8/8/3QK3 w - - 0 1",
			want: wantEqualsStatic,
		},
		{
			// white is in check with no captures available, only quiet king moves.
			// if the in-check branch used the capture generator there would be zero
			// legal moves and this would wrongly return a mate score
			name: "in check with only a quiet escape",
			fen:  "4k3/8/8/8/8/8/8/r3K3 w - - 0 1",
			want: wantNotMate,
		},
		{
			name:       "checkmate at ply 0",
			fen:        "R5k1/5ppp/8/8/8/8/8/6K1 b - - 0 1",
			want:       wantExact,
			exactScore: -eval.Mate,
		},
		{
			// same position, ply-relative mate score shifts by exactly the ply
			name:       "checkmate at ply 5",
			fen:        "R5k1/5ppp/8/8/8/8/8/6K1 b - - 0 1",
			ply:        5,
			want:       wantExact,
			exactScore: -eval.Mate + 5,
		},
		{
			// stands pat without quiet promotions in the generator
			name: "quiet promotion",
			fen:  "4k3/P7/8/8/8/8/8/4K3 w - - 0 1",
			want: wantBeatsStatic,
		},
		{
			// stands pat without en passant in the capture generator
			name: "en passant capture",
			fen:  "4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 1",
			want: wantBeatsStatic,
		},
		{
			// captures are available but the ply guard must return the static
			// eval untouched
			name: "ply limit returns static eval",
			fen:  "4k3/8/8/8/8/8/3r4/3QK3 w - - 0 1",
			ply:  MaxPly,
			want: wantEqualsStatic,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pos, err := board.ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("failed to parse FEN: %v", err)
			}

			static := eval.Evaluate(pos)

			searcher := New()
			score := searcher.quiescence(pos, tt.ply, -eval.Infinity, eval.Infinity)

			switch tt.want {
			case wantEqualsStatic:
				if score != static {
					t.Errorf("expected quiescence to return the static eval %v, got %v", static, score)
				}
			case wantBeatsStatic:
				if score <= static {
					t.Errorf("expected quiescence to beat the static eval %v, got %v", static, score)
				}
			case wantNotMate:
				if score > eval.Mate-MaxPly || score < -eval.Mate+MaxPly {
					t.Errorf("expected an ordinary score, got mate score %v", score)
				}
			case wantExact:
				if score != tt.exactScore {
					t.Errorf("expected score %v, got %v", tt.exactScore, score)
				}
			}
		})
	}
}

// not in check means you can always decline every capture, so the result can
// never be worse than the static eval
func TestQuiescenceStandPatFloor(t *testing.T) {
	fens := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
		"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
		"rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
		"2r3k1/p2q1pp1/1p2p2p/8/2PQ4/1P4P1/P4PKP/3R4 b - - 0 1",
	}

	for _, fen := range fens {
		pos, err := board.ParseFEN(fen)
		if err != nil {
			t.Fatalf("failed to parse FEN: %v", err)
		}
		if pos.InCheck(pos.SideToMove()) {
			continue // stand-pat does not apply when in check
		}

		static := eval.Evaluate(pos)
		searcher := New()
		score := searcher.quiescence(pos, 0, -eval.Infinity, eval.Infinity)

		if score < static {
			t.Errorf("%s: quiescence returned %v, below the stand-pat floor of %v", fen, score, static)
		}
	}
}
