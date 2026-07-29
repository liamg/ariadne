package board

import "testing"

func TestPieceTypeString(t *testing.T) {
	tests := []struct {
		pieceType PieceType
		expected  string
	}{
		{Pawn, "P"},
		{Knight, "N"},
		{Bishop, "B"},
		{Rook, "R"},
		{Queen, "Q"},
		{King, "K"},
	}

	for _, test := range tests {
		if test.pieceType.String() != test.expected {
			t.Errorf("Expected piece type %v to be %s, got %s", test.pieceType, test.expected, test.pieceType.String())
		}
	}
}

func TestColourString(t *testing.T) {
	tests := []struct {
		colour   Colour
		expected string
	}{
		{White, "w"},
		{Black, "b"},
	}

	for _, test := range tests {
		if test.colour.String() != test.expected {
			t.Errorf("Expected colour %v to be %s, got %s", test.colour, test.expected, test.colour.String())
		}
	}
}

func TestPieceRoundTrip(t *testing.T) {
	tests := []struct {
		colour    Colour
		pieceType PieceType
	}{
		{White, Pawn},
		{Black, Knight},
	}

	for _, test := range tests {
		piece := NewPiece(test.colour, test.pieceType)
		if piece.Colour() != test.colour {
			t.Errorf("Expected piece colour to be %v, got %v", test.colour, piece.Colour())
		}
		if piece.Type() != test.pieceType {
			t.Errorf("Expected piece type to be %v, got %v", test.pieceType, piece.Type())
		}
	}
}

func TestPieceString(t *testing.T) {
	tests := []struct {
		piece    Piece
		expected string
	}{
		{NewPiece(White, Pawn), "P"},
		{NewPiece(White, Knight), "N"},
		{NewPiece(White, Bishop), "B"},
		{NewPiece(White, Rook), "R"},
		{NewPiece(White, Queen), "Q"},
		{NewPiece(White, King), "K"},
		{NewPiece(Black, Pawn), "p"},
		{NewPiece(Black, Knight), "n"},
		{NewPiece(Black, Bishop), "b"},
		{NewPiece(Black, Rook), "r"},
		{NewPiece(Black, Queen), "q"},
		{NewPiece(Black, King), "k"},
	}

	for _, test := range tests {
		if test.piece.String() != test.expected {
			t.Errorf("Expected piece %v to be %s, got %s", test.piece, test.expected, test.piece.String())
		}
	}
}
