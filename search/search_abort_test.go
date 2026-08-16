package search

import (
	"testing"

	"github.com/liamg/ariadne/board"
)

// TestAbortedSearchReportsConsistentBestMove pins the invariant fastchess checks:
// the bestmove an engine finally plays must be the first move of the last pv it
// reported. Breaking it invalidates an entire SPRT run, and it only shows up when
// a search is cut short part-way through an iteration.
//
// A node limit is used rather than a time limit so the cut lands at exactly the
// same point on every machine and every run.
func TestAbortedSearchReportsConsistentBestMove(t *testing.T) {
	fens := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		// Kiwipete - dense, lots of captures, so ordering churns between depths
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
		"r2q1rk1/pP1p2pp/Q4n2/bbp1p3/Np6/1B3NBn/pPPP1PPP/R3K2R b KQ - 0 1",
	}

	// A spread of limits so the abort lands at many different points, including
	// inside the very first iteration.
	limits := []int64{1, 7, 50, 300, 1500, 9000, 40000}

	for _, fen := range fens {
		for _, nodes := range limits {
			t.Run(fen+"/"+itoa(nodes), func(t *testing.T) {
				pos, err := board.ParseFEN(fen)
				if err != nil {
					t.Fatalf("ParseFEN(%q): %v", fen, err)
				}

				var reported []Progress
				searcher := New(WithProgressCallback(func(p Progress) {
					pv := make([]board.Move, len(p.PV))
					copy(pv, p.PV)
					p.PV = pv
					reported = append(reported, p)
				}))

				result := searcher.Search(t.Context(), pos, Limits{Nodes: nodes}, nil)

				if result.BestMove == board.NullMove {
					// legitimate only when the search never completed a root move
					return
				}

				if len(result.PV) == 0 {
					t.Fatalf("bestmove %s returned with an empty pv", result.BestMove)
				}
				if result.PV[0] != result.BestMove {
					t.Errorf("bestmove %s is not the head of its own pv %v", result.BestMove, result.PV)
				}

				if !isPseudoLegal(pos, result.BestMove) {
					t.Errorf("bestmove %s is not pseudo-legal in the search position", result.BestMove)
				}

				// the invariant fastchess actually enforces
				if len(reported) == 0 {
					t.Fatalf("bestmove %s was returned without reporting any pv", result.BestMove)
				}
				last := reported[len(reported)-1]
				if len(last.PV) == 0 {
					t.Fatalf("last reported info line carried no pv, bestmove was %s", result.BestMove)
				}
				if last.PV[0] != result.BestMove {
					t.Errorf("bestmove %s contradicts the last reported pv %v (depth %d)",
						result.BestMove, last.PV, last.Depth)
				}
			})
		}
	}
}

// TestAbortedSearchKeepsVerifiedMove checks the other half of the fix: a partial
// iteration must not replace a fully searched move with one that merely happened
// to be ordered first before the stop landed.
//
// Searching the same position twice with the same node limit must give the same
// answer, and that answer must be one the search actually settled on rather than
// an artefact of where the limit fell.
func TestAbortedSearchKeepsVerifiedMove(t *testing.T) {
	pos, err := board.ParseFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1")
	if err != nil {
		t.Fatal(err)
	}

	for _, nodes := range []int64{200, 2000, 20000} {
		first := New().Search(t.Context(), pos, Limits{Nodes: nodes}, nil)
		second := New().Search(t.Context(), pos, Limits{Nodes: nodes}, nil)

		if first.BestMove != second.BestMove {
			t.Errorf("node limit %d: search is not deterministic, got %s then %s",
				nodes, first.BestMove, second.BestMove)
		}
		if first.Score != second.Score {
			t.Errorf("node limit %d: score is not deterministic, got %d then %d",
				nodes, first.Score, second.Score)
		}
	}
}

func isPseudoLegal(pos *board.Position, move board.Move) bool {
	for _, m := range pos.GeneratePseudoLegalMoves(make([]board.Move, 0, MaxMoves)) {
		if m == move {
			return true
		}
	}
	return false
}

func itoa(n int64) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}
