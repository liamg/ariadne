package eval

import (
	"strings"
	"testing"

	"github.com/liamg/ariadne/board"
)

// swap colours, flip the board vertically
func mirrorFEN(t *testing.T, fen string) string {
	t.Helper()

	fields := strings.Fields(fen)
	if len(fields) < 4 {
		t.Fatalf("malformed FEN %q", fen)
	}

	ranks := strings.Split(fields[0], "/")
	if len(ranks) != 8 {
		t.Fatalf("malformed FEN %q", fen)
	}
	for i, j := 0, len(ranks)-1; i < j; i, j = i+1, j-1 {
		ranks[i], ranks[j] = ranks[j], ranks[i]
	}
	pieces := strings.Map(swapCase, strings.Join(ranks, "/"))

	side := "w"
	if fields[1] == "w" {
		side = "b"
	}

	castling := "-"
	if fields[2] != "-" {
		var sb strings.Builder
		for _, r := range "KQkq" {
			if strings.ContainsRune(fields[2], swapCase(r)) {
				sb.WriteRune(r)
			}
		}
		castling = sb.String()
	}

	ep := "-"
	if fields[3] != "-" {
		// same file, rank 3 <-> rank 6
		ep = string(fields[3][0]) + string('1'+'8'-fields[3][1])
	}

	mirrored := []string{pieces, side, castling, ep}
	mirrored = append(mirrored, fields[4:]...)
	return strings.Join(mirrored, " ")
}

func swapCase(r rune) rune {
	switch {
	case r >= 'a' && r <= 'z':
		return r - 32
	case r >= 'A' && r <= 'Z':
		return r + 32
	}
	return r
}

// eval is side-to-move relative and mirroring swaps the mover, so the score
// comes back identical, not negated. catches pst orientation and sign bugs
func TestEvaluateIsColourSymmetric(t *testing.T) {
	fens := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
		"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
		"rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
		"2r3k1/p2q1pp1/1p2p2p/8/2PQ4/1P4P1/P4PKP/3R4 b - - 0 1",
		"4k3/8/8/3pP3/8/8/8/4K3 w - d6 0 1",
		"8/P1k5/8/8/8/8/8/7K w - - 0 1",
		"6k1/5ppp/8/8/8/8/8/R5K1 w - - 0 1",
	}

	for _, fen := range fens {
		pos, err := board.ParseFEN(fen)
		if err != nil {
			t.Fatalf("failed to parse FEN %q: %v", fen, err)
		}

		mirroredFEN := mirrorFEN(t, fen)
		mirrored, err := board.ParseFEN(mirroredFEN)
		if err != nil {
			t.Fatalf("failed to parse mirrored FEN %q (from %q): %v", mirroredFEN, fen, err)
		}

		if got, want := Evaluate(mirrored), Evaluate(pos); got != want {
			t.Errorf("%s scores %v, mirror %s scores %v", fen, want, mirroredFEN, got)
		}
	}
}

func allTables() map[string]*[2][7][64]int16 {
	return map[string]*[2][7][64]int16{
		"midgame": &pieceSquareTablesMidGame,
		"endgame": &pieceSquareTablesEndGame,
	}
}

// a short row shifts everything after it, which breaks file symmetry.
// queen is exempt - the published table is genuinely asymmetric
func TestPieceSquareTablesAreHorizontallySymmetric(t *testing.T) {
	for name, table := range allTables() {
		for pt := board.Pawn; pt <= board.King; pt++ {
			if pt == board.Queen {
				continue
			}
			for _, colour := range []board.Colour{board.White, board.Black} {
				for sq := board.A1; sq <= board.H8; sq++ {
					got := table[colour][pt][sq]
					want := table[colour][pt][sq^7]
					if got != want {
						t.Fatalf("%s piece %d colour %d: %v holds %d, file mirror holds %d", name, pt, colour, sq, got, want)
					}
				}
			}
		}
	}
}

// white's table must be black's vertical mirror - what init() builds
func TestPieceSquareTablesMirrorByColour(t *testing.T) {
	for name, table := range allTables() {
		for pt := board.Pawn; pt <= board.King; pt++ {
			for sq := board.A1; sq <= board.H8; sq++ {
				white := table[board.White][pt][sq]
				black := table[board.Black][pt][sq^56]
				if white != black {
					t.Fatalf("%s piece %d: white %v holds %d, black mirror holds %d", name, pt, sq, white, black)
				}
			}
		}
	}
}

// startpos is symmetric, so it must score 0
func TestStartingPositionIsBalanced(t *testing.T) {
	pos := board.StartingPosition()
	if score := Evaluate(pos); score != 0 {
		t.Errorf("expected 0, got %v", score)
	}
}

// promotions can push phase past maxPhase. clamped, the endgame weight is zero
// and the endgame tables cannot affect the score at all. unclamped the weight
// goes negative and they start subtracting.
//
// perturbing a table and checking nothing moves avoids restating the eval here,
// so this does not churn when new terms are added
func TestPhaseIsClampedToMaxPhase(t *testing.T) {
	// seven queens - phase 28 - white king on e3
	const fen = "k7/8/8/8/8/1Q2K3/8/1QQQQQQ1 w - - 0 1"

	pos, err := board.ParseFEN(fen)
	if err != nil {
		t.Fatalf("failed to parse FEN: %v", err)
	}

	before := Evaluate(pos)

	saved := pieceSquareTablesEndGame
	defer func() { pieceSquareTablesEndGame = saved }()
	pieceSquareTablesEndGame[board.White][board.King][board.E3] += 1000

	if after := Evaluate(pos); after != before {
		t.Errorf("endgame table affected the score at phase above maxPhase: %v became %v", before, after)
	}
}

func TestKingZones(t *testing.T) {
	for sq := board.A1; sq <= board.H8; sq++ {
		kingZone := kingZones[sq]
		if kingZone.Count() != 9 {
			t.Errorf("king zone for %v has %d squares, want 9", sq, kingZone.Count())
		}
		if !kingZone.Has(sq) {
			t.Errorf("king zone for %v does not contain the square itself", sq)
		}

		flippedKingZone := kingZones[sq^56]
		var mirror board.Bitboard
		var csq board.Square
		for flippedKingZone != 0 {
			csq, flippedKingZone = flippedKingZone.PopSquare()
			mirror |= (csq ^ 56).Bitboard()
		}
		if mirror != kingZones[sq] {
			t.Errorf("king zone for %v does not mirror king zone for %v", sq, sq^56)
		}
		flippedKingZone = kingZones[sq^7]
		mirror = 0
		for flippedKingZone != 0 {
			csq, flippedKingZone = flippedKingZone.PopSquare()
			mirror |= (csq ^ 7).Bitboard()
		}
		if mirror != kingZones[sq] {
			t.Errorf("king zone for %v does not mirror king zone for %v", sq, sq^7)
		}

	}
	if kingZones[board.G1] != kingZones[board.H1] {
		t.Errorf("g1 and h1 should provide the same king zones")
	}
	if kingZones[board.A1] != kingZones[board.B1] {
		t.Errorf("a1 and b1 should provide the same king zones")
	}
}

func TestKingSafetyFavoursTheAttackingSide(t *testing.T) {
	tests := []struct {
		fen      string
		favoured board.Colour
	}{
		{"4nrk1/ppp3b1/6q1/6NQ/8/3B4/PPP5/1K3R2 w - - 0 1", board.White},
		{"1k3r2/ppp5/3b4/8/6nq/6Q1/PPP3B1/4NRK1 w - - 0 1", board.Black},
	}

	for _, test := range tests {
		pos, err := board.ParseFEN(test.fen)
		if err != nil {
			t.Fatalf("failed to parse FEN %q: %v", test.fen, err)
		}

		score := Evaluate(pos)
		if test.favoured == board.White && score <= 0 {
			t.Errorf("%s: white has the attack but scores %v", test.fen, score)
		}
		if test.favoured == board.Black && score >= 0 {
			t.Errorf("%s: black has the attack but scores %v", test.fen, score)
		}
	}
}

func TestKingShelterPenalty(t *testing.T) {
	tests := []struct {
		king    board.Square
		pawns   []board.Square
		penalty int16
	}{
		{
			king:    board.G1,
			pawns:   []board.Square{board.F2, board.G2, board.H2},
			penalty: 0,
		},
		{
			king:    board.G1,
			pawns:   []board.Square{board.F2, board.H2},
			penalty: 30,
		},
		{
			king:    board.G1,
			pawns:   []board.Square{board.F2, board.G3, board.H2},
			penalty: 10,
		},
		{
			king:    board.G1,
			pawns:   []board.Square{board.F2, board.G4, board.H2},
			penalty: 20,
		},
		{
			king:    board.G1,
			pawns:   nil,
			penalty: 90,
		},
		{
			king:    board.H1,
			pawns:   []board.Square{board.F2, board.G2, board.H2},
			penalty: 0,
		},
		{
			king:    board.G2,
			pawns:   []board.Square{board.F2, board.G3, board.H2},
			penalty: 10,
		},
		{
			king:    board.G4,
			pawns:   []board.Square{board.F5, board.G5, board.H5},
			penalty: 90,
		},
		{
			king:    board.G6,
			pawns:   []board.Square{board.F7, board.G7, board.H7},
			penalty: 90,
		},
	}

	for _, test := range tests {
		calculatedPenalty := calculateKingShelterPenalty(test.king, board.BitboardFromSquares(test.pawns...))
		if calculatedPenalty != test.penalty {
			t.Errorf("king %v with pawns %v: expected penalty %d, got %d", test.king, test.pawns, test.penalty, calculatedPenalty)
		}
	}
}
