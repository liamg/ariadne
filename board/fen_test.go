package board

import (
	"math/rand/v2"
	"strings"
	"testing"
)

func TestParseFEN(t *testing.T) {
	tests := []struct {
		name        string
		fen         string
		expected    *Position
		expectedErr string
	}{
		{
			name:     "starting position",
			fen:      fenStarting,
			expected: StartingPosition(),
		},
		{
			name:     "only kings",
			fen:      "7k/8/8/8/8/8/8/K7 w - - 0 1",
			expected: EmptyPosition().addPiece(A1, WhiteKing).addPiece(H8, BlackKing),
		},
		{
			name:     "opposite kings on same file",
			fen:      "k7/8/8/8/8/8/8/K7 w - - 0 1",
			expected: EmptyPosition().addPiece(A1, WhiteKing).addPiece(A8, BlackKing),
		},
		{
			name:     "opposite kings on same rank",
			fen:      "8/8/8/8/8/8/8/K6k w - - 0 1",
			expected: EmptyPosition().addPiece(A1, WhiteKing).addPiece(H1, BlackKing),
		},

		// Classic FEN records the target square after every double push, whether or
		// not an enemy pawn is in position to use it. These must parse.
		{
			name: "en passant target with no capturer available",
			fen:  "7k/8/8/8/2P5/8/8/K7 b - c3 0 1",
			expected: func() *Position {
				p := EmptyPosition().addPiece(A1, WhiteKing).addPiece(H8, BlackKing).addPiece(C4, WhitePawn)
				p.sideToMove = Black
				return p
			}(),
		},
		{
			name: "black en passant target with no capturer available",
			fen:  "7k/8/8/2p5/8/8/8/K7 w - c6 0 1",
			expected: func() *Position {
				p := EmptyPosition().addPiece(A1, WhiteKing).addPiece(H8, BlackKing).addPiece(C5, BlackPawn)
				return p
			}(),
		},
		{
			name: "en passant target with a capturer available",
			fen:  "7k/8/8/8/2Pp4/8/8/K7 b - c3 0 1",
			expected: func() *Position {
				p := EmptyPosition().addPiece(A1, WhiteKing).addPiece(H8, BlackKing).
					addPiece(C4, WhitePawn).addPiece(D4, BlackPawn)
				p.sideToMove = Black
				p.state.enPassantSquare = C3
				return p
			}(),
		},

		// The target sits behind the pawn that pushed, so that pawn is always one
		// rank beyond it - never on the square it pushed from.
		{
			name:        "en passant target with no pawn behind it",
			fen:         "7k/8/8/8/8/8/8/K7 b - c3 0 1",
			expectedErr: "there is no white pawn on c4",
		},
		{
			name:        "en passant pawn still on its starting square",
			fen:         "7k/8/8/8/8/8/2P5/K7 b - c3 0 1",
			expectedErr: "there is no white pawn on c4",
		},
		{
			// occupancy is not enough - it has to be a pawn of the colour that
			// could have made the push
			name:        "en passant target with a wrongly coloured pawn behind it",
			fen:         "7k/8/8/8/2p5/8/8/K7 b - c3 0 1",
			expectedErr: "there is no white pawn on c4",
		},
		{
			name:        "black en passant target with no pawn behind it",
			fen:         "7k/2p5/8/8/8/8/8/K7 w - c6 0 1",
			expectedErr: "there is no black pawn on c5",
		},
		{
			// only a double push creates a target, so only ranks 3 and 6 are reachable
			name:        "en passant target on an impossible rank",
			fen:         "7k/8/8/8/8/8/8/K7 b - c4 0 1",
			expectedErr: "is not on rank 3 or 6",
		},
		{
			name:        "en passant target on the back rank",
			fen:         "7k/8/8/8/8/8/8/K7 b - c8 0 1",
			expectedErr: "is not on rank 3 or 6",
		},
		{
			// the pawn passed over the target, so nothing can be standing on it
			name:        "piece standing on the en passant target",
			fen:         "7k/8/8/8/2P5/2n5/8/K7 b - c3 0 1",
			expectedErr: "there is a piece on that square",
		},
		{
			// the pawn vacated this square, so it cannot still be occupied
			name:        "piece left on the square the pawn pushed from",
			fen:         "7k/8/8/8/2P5/8/2P5/K7 b - c3 0 1",
			expectedErr: "there is a piece on c2",
		},
		{
			// a rank 3 target implies a white double push, which means it is black's move
			name:        "en passant rank contradicts the side to move",
			fen:         "7k/8/8/8/2P5/8/8/K7 w - c3 0 1",
			expectedErr: "invalid en passant square: c3 is set, but it is not black's turn to move",
		},

		{
			name:        "no white king",
			fen:         "7k/8/8/8/8/8/8/8 w - - 0 1",
			expectedErr: "no white king on the board",
		},
		{
			name:        "castling right set without the rook",
			fen:         "4k3/8/8/8/8/8/8/4K3 w K - 0 1",
			expectedErr: "H1 does not have a white rook",
		},
		{
			name:        "pawns on invalid rank",
			fen:         "7k/8/8/8/8/8/8/K6P w - - 0 1",
			expectedErr: "pawns cannot be on the first or last rank",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pos, err := ParseFEN(test.fen)
			if err != nil {
				if test.expectedErr == "" {
					t.Fatalf("ParseFEN returned error: %v", err)
					return
				}
				if !strings.Contains(err.Error(), test.expectedErr) {
					t.Fatalf("ParseFEN returned error: %v, expected error containing: %v", err, test.expectedErr)
				}
				return
			} else if test.expectedErr != "" {
				t.Fatalf("ParseFEN did not return an error, but was expected to: %v", test.expectedErr)
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

		if pos.Validate() != nil {
			continue
		}

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
