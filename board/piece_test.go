package board

import "testing"

func TestPieceTypeString(t *testing.T) {
	tests := []struct {
		pieceType PieceType
		expected  string
	}{
		{PieceTypePawn, "P"},
		{PieceTypeKnight, "N"},
		{PieceTypeBishop, "B"},
		{PieceTypeRook, "R"},
		{PieceTypeQueen, "Q"},
		{PieceTypeKing, "K"},
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
		{White, PieceTypePawn},
		{Black, PieceTypeKnight},
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
		{NewPiece(White, PieceTypePawn), "P"},
		{NewPiece(White, PieceTypeKnight), "N"},
		{NewPiece(White, PieceTypeBishop), "B"},
		{NewPiece(White, PieceTypeRook), "R"},
		{NewPiece(White, PieceTypeQueen), "Q"},
		{NewPiece(White, PieceTypeKing), "K"},
		{NewPiece(Black, PieceTypePawn), "p"},
		{NewPiece(Black, PieceTypeKnight), "n"},
		{NewPiece(Black, PieceTypeBishop), "b"},
		{NewPiece(Black, PieceTypeRook), "r"},
		{NewPiece(Black, PieceTypeQueen), "q"},
		{NewPiece(Black, PieceTypeKing), "k"},
	}

	for _, test := range tests {
		if test.piece.String() != test.expected {
			t.Errorf("Expected piece %v to be %s, got %s", test.piece, test.expected, test.piece.String())
		}
	}
}
