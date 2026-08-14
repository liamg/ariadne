package eval

import (
	"testing"

	"github.com/liamg/ariadne/board"
)

// obvious per-piece attack maps, kept as the reference the hoisted pawn and
// king shifts in Evaluate are checked against. slow and boring on purpose
func naiveAttackMaps(p *board.Position, c board.Colour) (reach [7]board.Bitboard, all, twice board.Bitboard) {
	occ := p.Occupancy()
	for pt := board.Pawn; pt <= board.King; pt++ {
		pieces := p.Pieces(c, pt)
		for pieces != 0 {
			var sq board.Square
			sq, pieces = pieces.PopSquare()
			att := p.AttacksWithCustomOccupancy(pt, sq, occ, c)
			twice |= all & att
			all |= att
			reach[pt] |= att
		}
	}
	return reach, all, twice
}

// Evaluate builds pawn attacks with two shifts instead of a lookup per pawn,
// and relies on a square taking at most one pawn from each diagonal
func TestPawnAttackShiftsMatchPerPawnLookups(t *testing.T) {
	fens := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
		"4k3/pp4pp/2p2p2/1P1PP1P1/8/8/8/4K3 w - - 0 1",
		"8/8/8/2P1P3/3P4/2P1P3/8/k6K w - - 0 1",
	}

	for _, fen := range fens {
		pos, err := board.ParseFEN(fen)
		if err != nil {
			t.Fatalf("failed to parse FEN %q: %v", fen, err)
		}

		for _, c := range []board.Colour{board.White, board.Black} {
			pawns := pos.Pieces(c, board.Pawn)
			gotAll, gotTwice := pawnAttackMaps(c, pawns)

			var wantAll, wantTwice board.Bitboard
			rest := pawns
			for rest != 0 {
				var sq board.Square
				sq, rest = rest.PopSquare()
				att := pos.AttacksWithCustomOccupancy(board.Pawn, sq, pos.Occupancy(), c)
				wantTwice |= wantAll & att
				wantAll |= att
			}

			if gotAll != wantAll {
				t.Errorf("%s colour %d: shifts give %v, per-pawn gives %v", fen, c, gotAll, wantAll)
			}
			if gotTwice != wantTwice {
				t.Errorf("%s colour %d: doubled shifts give %v, per-pawn gives %v", fen, c, gotTwice, wantTwice)
			}
		}
	}
}

// counts are hand-computed. the last pair is the double attack case - the king
// is the only defender of h2, so the check only counts when a second black
// piece hits it
func TestSafeCheckWeight(t *testing.T) {
	tests := []struct {
		name     string
		fen      string
		attacker board.Colour
		want     map[board.PieceType]int
	}{
		{
			name:     "knight check on an undefended square",
			fen:      "6k1/5ppp/8/3N4/8/8/8/6K1 w - - 0 1",
			attacker: board.White,
			want:     map[board.PieceType]int{board.Knight: 1},
		},
		{
			name:     "same knight check, now covered by a rook",
			fen:      "6k1/r4ppp/8/3N4/8/8/8/6K1 w - - 0 1",
			attacker: board.White,
			want:     map[board.PieceType]int{},
		},
		{
			name:     "rook reaches an undefended back rank check",
			fen:      "6k1/5ppp/8/8/8/8/8/R5K1 w - - 0 1",
			attacker: board.White,
			want:     map[board.PieceType]int{board.Rook: 1},
		},
		{
			name:     "same rook check, now defended",
			fen:      "1r4k1/5ppp/8/8/8/8/8/R5K1 w - - 0 1",
			attacker: board.White,
			want:     map[board.PieceType]int{},
		},
		{
			name:     "queen contact checks are king covered, distant ones are not",
			fen:      "6k1/5ppp/8/3Q4/8/8/8/6K1 w - - 0 1",
			attacker: board.White,
			want:     map[board.PieceType]int{board.Queen: 2},
		},
		{
			name:     "knight joins the queen on h2, so the contact check counts",
			fen:      "4k3/8/3q4/8/6n1/8/8/7K w - - 0 1",
			attacker: board.Black,
			want:     map[board.PieceType]int{board.Knight: 1, board.Queen: 5},
		},
		{
			name:     "queen alone on h2, so the contact check does not",
			fen:      "4k3/8/3q4/8/8/8/8/7K w - - 0 1",
			attacker: board.Black,
			want:     map[board.PieceType]int{board.Queen: 4},
		},
		{
			// b2 is hit twice by black but also covered twice by white, so
			// recapturing is not enough - a7, a8, e5 and h8 carry the count
			name:     "b2 is defended twice, so double attacking it is not enough",
			fen:      "1q6/8/6k1/8/8/3n4/8/K2N4 w - - 0 1",
			attacker: board.Black,
			want:     map[board.PieceType]int{board.Queen: 4},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pos, err := board.ParseFEN(test.fen)
			if err != nil {
				t.Fatalf("failed to parse FEN %q: %v", test.fen, err)
			}

			defender := board.White
			if test.attacker == board.White {
				defender = board.Black
			}

			attackerReach, _, attackerTwice := naiveAttackMaps(pos, test.attacker)
			_, defenderAll, defenderTwice := naiveAttackMaps(pos, defender)

			var want int16
			var wantAny bool
			for pt, n := range test.want {
				want += int16(n) * safeCheckWeights[pt]
				if n > 0 {
					wantAny = true
				}
			}

			got, gotAny := calculateSafeCheckWeight(pos, test.attacker, pos.KingSquare(defender),
				pos.Occupancy(), &attackerReach, defenderAll, defenderTwice, attackerTwice)

			if got != want {
				t.Errorf("weight is %d, want %d", got, want)
			}
			if gotAny != wantAny {
				t.Errorf("any safe check is %v, want %v", gotAny, wantAny)
			}
		})
	}
}
