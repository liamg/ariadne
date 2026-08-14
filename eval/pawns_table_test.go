package eval

import (
	"testing"

	"github.com/liamg/ariadne/board"
)

func TestPawnsTableScore(t *testing.T) {
	// test that the cache is actually used by manually populating a value
	// and checking it gets used
	pos := board.StartingPosition()
	pawns := pos.Pawns()
	whitePawns := pawns & pos.PiecesByColour(board.White)
	blackPawns := pawns & pos.PiecesByColour(board.Black)

	e := New()
	before := e.Evaluate(pos) // populate the cache

	entry, _ := e.pawnsTable.probe(whitePawns, blackPawns)
	entry.midGameScore = 30_000

	after := e.Evaluate(pos)

	if before == after {
		t.Errorf("Expected cached value to be used, but got the same value before and after: before=%d, after=%d", before, after)
	}

	// then check if the evaluator writes to the cache

	pos, err := board.ParseFEN("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1")
	if err != nil {
		t.Fatalf("Failed to parse FEN: %v", err)
	}

	whitePawns = pos.Pawns() & pos.PiecesByColour(board.White)
	blackPawns = pos.Pawns() & pos.PiecesByColour(board.Black)

	_ = e.Evaluate(pos) // populate the cache

	var ok bool
	entry, ok = e.pawnsTable.probe(whitePawns, blackPawns)
	if !ok {
		t.Errorf("Expected cache entry to be present, but it was not found")
	}

	if entry.midGameScore == 0 {
		t.Errorf("Expected mid game score value to be populated, but it is zero")
	}
}

func TestPawnsTableShelter(t *testing.T) {
	// test that the cache is actually used by manually populating a value
	// and checking it gets used
	pos := board.StartingPosition()
	pawns := pos.Pawns()
	whitePawns := pawns & pos.PiecesByColour(board.White)
	blackPawns := pawns & pos.PiecesByColour(board.Black)

	e := New()
	before := e.Evaluate(pos) // populate the cache

	entry, _ := e.pawnsTable.probe(whitePawns, blackPawns)
	entry.shelter[board.White] = 30_000

	after := e.Evaluate(pos)

	if before == after {
		t.Errorf("Expected cached value to be used, but got the same value before and after: before=%d, after=%d", before, after)
	}

	// then check if the evaluator writes to the cache

	pos, err := board.ParseFEN(board.FenKiwiPete)
	if err != nil {
		t.Fatalf("Failed to parse FEN: %v", err)
	}

	whitePawns = pos.Pawns() & pos.PiecesByColour(board.White)
	blackPawns = pos.Pawns() & pos.PiecesByColour(board.Black)

	_ = e.Evaluate(pos) // populate the cache

	var ok bool
	entry, ok = e.pawnsTable.probe(whitePawns, blackPawns)
	if !ok {
		t.Errorf("Expected cache entry to be present, but it was not found")
	}

	if entry.shelter[board.White] == 0 && entry.shelter[board.Black] == 0 {
		t.Errorf("Expected shelter values to be populated, but they are zero")
	}
}
