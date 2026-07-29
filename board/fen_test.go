package board

import (
	"math/rand/v2"
	"testing"
)

func TestParseFEN(t *testing.T) {
	tests := []struct {
		name     string
		fen      string
		expected *Position
	}{
		{
			name:     "starting position",
			fen:      fenStarting,
			expected: StartingPosition(),
		},
		{
			name:     "empty board",
			fen:      fenEmpty,
			expected: EmptyPosition(),
		},
		{
			name:     "single rook",
			fen:      "8/8/8/8/8/8/8/R7 w - - 0 1",
			expected: EmptyPosition().addPiece(A1, WhiteRook),
		},
		{
			name:     "opposite rooks on same file",
			fen:      "r7/8/8/8/8/8/8/R7 w - - 0 1",
			expected: EmptyPosition().addPiece(A1, WhiteRook).addPiece(A8, BlackRook),
		},
		{
			name:     "opposite rooks on same rank",
			fen:      "8/8/8/8/8/8/8/R6r w - - 0 1",
			expected: EmptyPosition().addPiece(A1, WhiteRook).addPiece(H1, BlackRook),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pos, err := ParseFEN(test.fen)
			if err != nil {
				t.Fatalf("ParseFEN returned error: %v", err)
			}
			assertPositionsMatch(t, test.expected, pos)
			if err := pos.Validate(); err != nil {
				t.Errorf("Position validation failed: %v", err)
			}
		})
	}
}

func TestParseFENKiwiPete(t *testing.T) {
	pos, err := ParseFEN(fenKiwiPete)
	if err != nil {
		t.Fatalf("ParseFEN returned error: %v", err)
	}
	// placement here will be covered b ythe round trip, so we just need to assert the other stuff
	if pos.sideToMove != White {
		t.Errorf("ActiveColor = %v; want %v", pos.sideToMove, White)
	}

	if pos.state.castlingRights != CastlingAll {
		t.Errorf("CastlingRights = %v; want %v", pos.state.castlingRights, CastlingAll)
	}

	if pos.state.enPassantSquare != NoSquare {
		t.Errorf("EnPassantSquare = %v; want %v", pos.state.enPassantSquare, NoSquare)
	}

	if pos.state.halfMoveClock != 0 {
		t.Errorf("HalfMoveClock = %v; want %v", pos.state.halfMoveClock, 0)
	}

	if pos.fullMoveNumber != 1 {
		t.Errorf("FullMoveNumber = %v; want %v", pos.fullMoveNumber, 1)
	}
}

func TestFENRoundTripFuzz(t *testing.T) {
	rnd := rand.New(rand.NewPCG(1, 2))
	for range 1000 {
		pos := RandomPosition(rnd)

		if err := pos.validateStateSlow(); err != nil {
			t.Errorf("Position validation failed: %v", err)
		}

		fen := GenerateFEN(pos)
		parsedPos, err := ParseFEN(fen)
		if err != nil {
			t.Fatalf("ParseFEN returned error: %v", err)
		}

		if err := parsedPos.validateStateSlow(); err != nil {
			t.Errorf("Position validation failed: %v", err)
		}
		assertPositionsMatch(t, pos, parsedPos)
	}
}

func assertPositionsMatch(t *testing.T, expected, actual *Position) {
	if expected == nil && actual == nil {
		return
	}
	if expected == nil || actual == nil {
		t.Fatalf("One of the positions is nil: expected %v, actual %v", expected, actual)
	}
	if expected.sideToMove != actual.sideToMove {
		t.Fatalf("Side to move mismatch: expected %v, actual %v", expected.sideToMove, actual.sideToMove)
	}
	if expected.state.castlingRights != actual.state.castlingRights {
		t.Fatalf("Castling rights mismatch: expected %v, actual %v", expected.state.castlingRights, actual.state.castlingRights)
	}
	if expected.state.enPassantSquare != actual.state.enPassantSquare {
		t.Fatalf("En passant square mismatch: expected %v, actual %v", expected.state.enPassantSquare, actual.state.enPassantSquare)
	}
	if expected.state.halfMoveClock != actual.state.halfMoveClock {
		t.Fatalf("Half move clock mismatch: expected %v, actual %v", expected.state.halfMoveClock, actual.state.halfMoveClock)
	}
	if expected.fullMoveNumber != actual.fullMoveNumber {
		t.Fatalf("Full move number mismatch: expected %v, actual %v", expected.fullMoveNumber, actual.fullMoveNumber)
	}
	for i := range 2 {
		if expected.byColour[i] != actual.byColour[i] {
			t.Fatalf("Bitboard mismatch for colour %v: expected %v, actual %v", Colour(i), expected.byColour[i], actual.byColour[i])
		}
	}
	for i := range 6 {
		if expected.byType[i] != actual.byType[i] {
			t.Fatalf("Bitboard mismatch for piece type %v: expected %v, actual %v", PieceType(i), expected.byType[i], actual.byType[i])
		}
	}
}

func TestParseFENInvalid(t *testing.T) {
	tests := []struct {
		name string
		fen  string
	}{
		{
			name: "invalid piece character",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNX w KQkq - 0 1",
		},
		{
			name: "invalid rank count",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPP w KQkq - 0 1",
		},
		{
			name: "invalid side to move",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR x KQkq - 0 1",
		},
		{
			name: "invalid castling rights",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkx - 0 1",
		},
		{
			name: "invalid en passant square",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq z9 0 1",
		},
		{
			name: "invalid halfmove clock",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - -1 1",
		},
		{
			name: "invalid fullmove number",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 0",
		},
		{
			name: "incomplete rank",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPP/RNBQKBNR w KQkq - 0 1",
		},
		{
			name: "extra rank",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPPP/RNBQKBNRR w KQkq - 0 1",
		},
		{
			name: "extra file",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPPP/RNBQKBNR w KQkq - 0 1",
		},
		{
			name: "missing field(s)",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq",
		},
		{
			name: "empty",
			fen:  "",
		},
		{
			name: "double slash",
			fen:  "rnbqkbnr/pppppppp//8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseFEN(tt.fen)
			if err == nil {
				t.Errorf("ParseFEN(%q) did not return an error, but was expected to", tt.fen)
			}
		})
	}
}
