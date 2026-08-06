package eval

import (
	"testing"

	"github.com/liamg/ariadne/board"
)

// obvious per-file implementations, kept as the reference the bitboard versions
// are checked against. slow and boring on purpose

func naiveDoubledPawns(pawns board.Bitboard) int {
	var n int
	for f := board.FileA; f <= board.FileH; f++ {
		if c := (pawns & f.Mask()).Count(); c > 1 {
			n += c - 1
		}
	}
	return n
}

func naiveIsolatedPawns(pawns board.Bitboard) int {
	var n int
	for f := board.FileA; f <= board.FileH; f++ {
		c := (pawns & f.Mask()).Count()
		if c == 0 {
			continue
		}
		var adjacent board.Bitboard
		if f > board.FileA {
			adjacent |= (f - 1).Mask()
		}
		if f < board.FileH {
			adjacent |= (f + 1).Mask()
		}
		if pawns&adjacent == 0 {
			n += c
		}
	}
	return n
}

// a pawn is passed if no enemy pawn sits on its own or an adjacent file, on any
// rank strictly ahead of it. an enemy pawn level with ours cannot stop it, as
// pawns capture away from the direction we are heading
func naivePassedPawns(colour board.Colour, pawns, enemyPawns board.Bitboard) int {
	var n int
	for bb := pawns; bb != 0; {
		var sq board.Square
		sq, bb = bb.PopSquare()

		passed := true
		for file := board.FileA; file <= board.FileH && passed; file++ {
			if file != sq.File() && file != sq.File()-1 && file != sq.File()+1 {
				continue
			}
			for rank := board.Rank1; rank <= board.Rank8; rank++ {
				ahead := rank > sq.Rank()
				if colour == board.Black {
					ahead = rank < sq.Rank()
				}
				if !ahead {
					continue
				}
				blocker, err := board.SquareFromFileAndRank(file, rank)
				if err != nil {
					continue
				}
				if enemyPawns&blocker.Bitboard() != 0 {
					passed = false
					break
				}
			}
		}
		if passed {
			n++
		}
	}
	return n
}

// a pawn is backward if an enemy pawn attacks the square in front of it, and no
// friendly pawn on an adjacent file is at or behind its rank - so it can never
// advance safely and can never be defended by a pawn
func naiveBackwardPawns(colour board.Colour, pawns, enemyPawns board.Bitboard) int {
	forward := 1
	if colour == board.Black {
		forward = -1
	}

	pawnAt := func(file, rank int, bb board.Bitboard) bool {
		if file < 0 || file > 7 || rank < 0 || rank > 7 {
			return false
		}
		sq, err := board.SquareFromFileAndRank(board.File(file), board.Rank(rank))
		if err != nil {
			return false
		}
		return bb&sq.Bitboard() != 0
	}

	var n int
	for bb := pawns; bb != 0; {
		var sq board.Square
		sq, bb = bb.PopSquare()
		file, rank := int(sq.File()), int(sq.Rank())

		stopRank := rank + forward
		if stopRank < 0 || stopRank > 7 {
			continue
		}

		// an enemy pawn attacks our stop square from one rank beyond it
		attacked := pawnAt(file-1, stopRank+forward, enemyPawns) ||
			pawnAt(file+1, stopRank+forward, enemyPawns)
		if !attacked {
			continue
		}

		// a friendly pawn at or behind us on an adjacent file could come up
		supported := false
		for _, adjacent := range []int{file - 1, file + 1} {
			for r := 0; r <= 7; r++ {
				atOrBehind := r <= rank
				if colour == board.Black {
					atOrBehind = r >= rank
				}
				if atOrBehind && pawnAt(adjacent, r, pawns) {
					supported = true
				}
			}
		}
		if !supported {
			n++
		}
	}
	return n
}

// a pawn is unopposed if no enemy pawn stands anywhere ahead of it on its own file
func naiveUnopposedPawns(colour board.Colour, pawns, enemyPawns board.Bitboard) int {
	var n int
	for bb := pawns; bb != 0; {
		var sq board.Square
		sq, bb = bb.PopSquare()

		opposed := false
		for rank := board.Rank1; rank <= board.Rank8; rank++ {
			ahead := rank > sq.Rank()
			if colour == board.Black {
				ahead = rank < sq.Rank()
			}
			if !ahead {
				continue
			}
			blocker, err := board.SquareFromFileAndRank(sq.File(), rank)
			if err != nil {
				continue
			}
			if enemyPawns&blocker.Bitboard() != 0 {
				opposed = true
				break
			}
		}
		if !opposed {
			n++
		}
	}
	return n
}

// the fill and shift tricks are easy to get subtly wrong, and wrong in the same
// way for both colours - which the symmetry test cannot see. walking a real
// tree and comparing against the naive version is what catches it
func TestPawnStructureMatchesNaive(t *testing.T) {
	fens := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",
		"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
		"rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
		"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
		"2r3k1/p2q1pp1/1p2p2p/8/2PQ4/1P4P1/P4PKP/3R4 b - - 0 1",
	}

	for _, fen := range fens {
		pos, err := board.ParseFEN(fen)
		if err != nil {
			t.Fatalf("failed to parse FEN %q: %v", fen, err)
		}
		walkPawnStructure(t, pos, 3)
	}
}

func walkPawnStructure(t *testing.T, p *board.Position, depth int) {
	t.Helper()

	for _, colour := range []board.Colour{board.White, board.Black} {
		pawns := p.Pieces(colour, board.Pawn)
		enemyPawns := p.Pieces(colour.Opposite(), board.Pawn)

		if got, want := doubledPawns(colour, pawns).Count(), naiveDoubledPawns(pawns); got != want {
			t.Fatalf("%s\n  %v doubled: bitboard=%d naive=%d", board.GenerateFEN(p), colour, got, want)
		}
		if got, want := isolatedPawns(pawns).Count(), naiveIsolatedPawns(pawns); got != want {
			t.Fatalf("%s\n  %v isolated: bitboard=%d naive=%d", board.GenerateFEN(p), colour, got, want)
		}
		if got, want := passedPawns(colour, pawns, enemyPawns).Count(), naivePassedPawns(colour, pawns, enemyPawns); got != want {
			t.Fatalf("%s\n  %v passed: bitboard=%d naive=%d", board.GenerateFEN(p), colour, got, want)
		}
		if got, want := backwardPawns(colour, pawns, enemyPawns).Count(), naiveBackwardPawns(colour, pawns, enemyPawns); got != want {
			t.Fatalf("%s\n  %v backward: bitboard=%d naive=%d", board.GenerateFEN(p), colour, got, want)
		}
		if got, want := unopposedPawns(colour, pawns, enemyPawns).Count(), naiveUnopposedPawns(colour, pawns, enemyPawns); got != want {
			t.Fatalf("%s\n  %v unopposed: bitboard=%d naive=%d", board.GenerateFEN(p), colour, got, want)
		}
	}

	if depth == 0 {
		return
	}

	moves := p.GeneratePseudoLegalMoves(make([]board.Move, 0, 256))
	for _, move := range moves {
		undo := p.MakeMove(move)
		if !p.IsLastMoveIllegal() {
			walkPawnStructure(t, p, depth-1)
		}
		p.UnmakeMove(undo)
	}
}
