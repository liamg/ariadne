package search

import (
	"math/rand/v2"
	"testing"

	"github.com/liamg/ariadne/board"
)

// seeCases are positions with a hand-resolved exchange on a single square.
//
// The 4a/4b pair differ only by the white rook on e1: without a working x-ray
// update the second rook never joins the exchange and both cases return -400,
// so 4b passing alone proves nothing.
var seeCases = []struct {
	name string
	fen  string
	move string
	want int32
}{
	{
		// nothing defends e5, so the exchange ends immediately
		name: "undefended pawn",
		fen:  "k7/8/8/4p3/8/8/8/K3R3 w - - 0 1",
		move: "e1e5",
		want: 100,
	},
	{
		name: "equal trade, rook for rook",
		fen:  "k3r3/8/8/4r3/8/8/8/K3R3 w - - 0 1",
		move: "e1e5",
		want: 0,
	},
	{
		name: "queen takes pawn defended by pawn",
		fen:  "k7/8/4p3/3p4/8/8/8/K2Q4 w - - 0 1",
		move: "d1d5",
		want: -800,
	},
	{
		// Rxe5 Nxe5 Rxe5 - the second rook only joins once the first vacates e2
		name: "x-ray, doubled rooks behind the capture",
		fen:  "k7/3n4/8/4p3/8/8/4R3/K3R3 w - - 0 1",
		move: "e2e5",
		want: -80,
	},
	{
		name: "x-ray, same position without the rook behind",
		fen:  "k7/3n4/8/4p3/8/8/4R3/K7 w - - 0 1",
		move: "e2e5",
		want: -400,
	},
	{
		// black must recapture with the e6 pawn, not the d8 queen;
		// recapturing queen-first returns +100 instead
		name: "least valuable attacker recaptures first",
		fen:  "k2q4/8/4p3/3p4/2P5/8/8/K2R4 w - - 0 1",
		move: "d1d5",
		want: -400,
	},
	{
		// the black king is the only defender, but e5 is still covered by the
		// rook on e1, so Kxe5 is not available and white simply wins the queen
		name: "king may not recapture into a defended square",
		fen:  "8/8/3k4/4q3/8/8/4R3/K3R3 w - - 0 1",
		move: "e2e5",
		want: 900,
	},
	{
		// the captured pawn stands on d5, not on d6; clearing it is what lets
		// the d4 rook see through to the landing square
		name: "en passant clears the captured pawn's square",
		fen:  "k7/8/8/3pP3/3r4/8/8/K7 w - d6 0 1",
		move: "e5d6",
		want: 0,
	},
	{
		// after fxe5 white can play Qxe5, but the e8 rook then wins the queen,
		// so white stops instead and the exchange resolves at -400
		name: "side to move declines the final recapture",
		fen:  "k3r3/8/5p2/4p3/8/8/4R3/K3Q3 w - - 0 1",
		move: "e2e5",
		want: -400,
	},
	{
		name: "quiet move onto a square attacked by a pawn",
		fen:  "7k/8/1p6/8/8/8/8/R6K w - - 0 1",
		move: "a1a5",
		want: -500,
	},
}

func TestSEE(t *testing.T) {
	for _, tc := range seeCases {
		t.Run(tc.name, func(t *testing.T) {
			pos := mustParseFEN(t, tc.fen)
			move := mustFindMove(t, pos, tc.move)

			if got := see(pos, move); got != tc.want {
				t.Errorf("see(%s) = %d; want %d\n%s", tc.move, got, tc.want, pos)
			}
		})
	}
}

// TestSEESlowReference runs the same table through the reference implementation.
// When TestSEE fails, this says which of the two is wrong.
func TestSEESlowReference(t *testing.T) {
	for _, tc := range seeCases {
		t.Run(tc.name, func(t *testing.T) {
			pos := mustParseFEN(t, tc.fen)
			move := mustFindMove(t, pos, tc.move)

			if got := seeSlow(pos, move); got != tc.want {
				t.Errorf("seeSlow(%s) = %d; want %d\n%s", tc.move, got, tc.want, pos)
			}
		})
	}
}

// TestSEEMatchesSlowReference is the test that finds what reading the code does
// not. Every pseudo-legal move of a large number of random positions goes
// through both implementations and the results must agree exactly.
//
// RandomPosition randomises the side to move and sets a real en passant square
// about a fifth of the time, so both colours and the en passant path are
// covered here as well as in the table above.
//
// Most generated positions are rejected by Validate, almost always for castling
// rights that do not match the piece placement, so the position count has to be
// an order of magnitude above the sample size wanted.
func TestSEEMatchesSlowReference(t *testing.T) {
	rnd := rand.New(rand.NewPCG(1, 2))

	positions, minChecked := 20000, 20000
	if testing.Short() {
		positions, minChecked = 2000, 2000
	}

	moves := make([]board.Move, 0, 256)
	var checked int

	for range positions {
		pos := board.RandomPosition(rnd)
		if pos.Validate() != nil {
			continue
		}

		moves = pos.GeneratePseudoLegalMoves(moves[:0])
		for _, move := range moves {
			want := seeSlow(pos, move)
			got := see(pos, move)
			if got != want {
				t.Fatalf("see disagrees with the reference implementation\nfen:  %s\nmove: %s\nsee:  %d\nslow: %d\n%s",
					board.GenerateFEN(pos), move, got, want, pos)
			}
			checked++
		}
	}

	if checked < minChecked {
		t.Fatalf("only %d moves checked - the corpus is too small to mean anything", checked)
	}
	t.Logf("checked %d moves across %d positions", checked, positions)
}

// seeSlow is the deliberately naive reference implementation of static exchange
// evaluation, kept so that see() has something to be verified against.
//
// It recurses instead of building a swap list, and rebuilds the whole attacker
// set from scratch at every step instead of maintaining it incrementally. That
// makes it far too slow for the search and short enough to check by reading.
func seeSlow(pos *board.Position, move board.Move) int32 {
	from := move.From()
	to := move.To()
	side := pos.SideToMove()

	occ := pos.Occupancy() &^ from.Bitboard()

	var captured int32
	if move.Kind() == board.EnPassantCapture {
		captured = seeValues[board.Pawn]
		if side == board.White {
			occ &^= (to - 8).Bitboard()
		} else {
			occ &^= (to + 8).Bitboard()
		}
	} else {
		captured = seeValues[pos.PieceAt(to).Type()]
	}

	// the initial capture is forced by the caller, so it is not subject to the
	// max(0, ...) that lets every later side walk away
	return captured - seeSlowExchange(pos, to, side.Opposite(), occ, pos.PieceAt(from).Type())
}

// seeSlowExchange returns the material the given side wins from the exchange on
// `to`, where a piece of type `victim` is currently standing there. Returning
// zero means the side declines to capture at all.
func seeSlowExchange(pos *board.Position, to board.Square, side board.Colour, occ board.Bitboard, victim board.PieceType) int32 {
	attackers := pos.AllAttackersForSquare(to, occ)

	sideAttackers := attackers & pos.PiecesByColour(side)
	if sideAttackers == 0 {
		return 0
	}

	var from board.Square
	var attacker board.PieceType
	for pt := board.Pawn; pt <= board.King; pt++ {
		if set := sideAttackers & pos.PiecesByType(pt); set != 0 {
			from, _ = set.PopSquare()
			attacker = pt
			break
		}
	}

	if attacker == board.King && attackers&pos.PiecesByColour(side.Opposite()) != 0 {
		return 0 // the king cannot capture into a square the enemy still covers
	}

	return max(0, seeValues[victim]-seeSlowExchange(pos, to, side.Opposite(), occ&^from.Bitboard(), attacker))
}

func mustParseFEN(t *testing.T, fen string) *board.Position {
	t.Helper()
	pos, err := board.ParseFEN(fen)
	if err != nil {
		t.Fatalf("ParseFEN(%q): %v", fen, err)
	}
	return pos
}

func mustFindMove(t *testing.T, pos *board.Position, uci string) board.Move {
	t.Helper()
	for _, move := range pos.GeneratePseudoLegalMoves(make([]board.Move, 0, 256)) {
		if move.String() == uci {
			return move
		}
	}
	t.Fatalf("%s is not pseudo-legal here\n%s", uci, pos)
	return board.NullMove
}

func BenchmarkSEE(b *testing.B) {
	pos, err := board.ParseFEN("k7/3n4/8/4p3/8/8/4R3/K3R3 w - - 0 1")
	if err != nil {
		b.Fatal(err)
	}

	var move board.Move
	for _, m := range pos.GeneratePseudoLegalMoves(make([]board.Move, 0, 256)) {
		if m.String() == "e2e5" {
			move = m
			break
		}
	}

	b.ResetTimer()
	for range b.N {
		if see(pos, move) != -80 {
			b.Fatal("unexpected result")
		}
	}
}
