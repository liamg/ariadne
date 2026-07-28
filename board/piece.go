package board

import "strings"

// PieceType is the type of chess piece, i.e. one of pawn, knight, bishop, rook, queen, king.
// It does NOT encode the colour of the piece, which is stored separately in the Piece type.
type PieceType uint8

const (
	PieceTypeNone PieceType = iota
	PieceTypePawn
	PieceTypeKnight
	PieceTypeBishop
	PieceTypeRook
	PieceTypeQueen
	PieceTypeKing
)

func (p PieceType) String() string {
	switch p {
	case PieceTypePawn:
		return "P"
	case PieceTypeKnight:
		return "N"
	case PieceTypeBishop:
		return "B"
	case PieceTypeRook:
		return "R"
	case PieceTypeQueen:
		return "Q"
	case PieceTypeKing:
		return "K"
	default:
		return "?"
	}
}

// Colour encodes the colour of a piece, either white or black.
type Colour uint8

const (
	White Colour = iota
	Black
)

func (c Colour) String() string {
	if c == White {
		return "w"
	}
	return "b"
}

func (c Colour) Opposite() Colour {
	if c == White {
		return Black
	}
	return White
}

// Piece encodes a piece type and a colour e.g. a black knight
// The least significant 3 bits = piece type, the 4th bit = colour
type Piece uint8

// NewPiece creates a new Piece from a colour and piece type.
func NewPiece(colour Colour, pieceType PieceType) Piece {
	if pieceType == PieceTypeNone { // avoid encoding a black non-piece as >0
		return Piece(pieceType)
	}
	return Piece(colour<<3) | Piece(pieceType)
}

func (p Piece) Colour() Colour {
	return Colour((p & 8) >> 3)
}

func (p Piece) Type() PieceType {
	return PieceType(p & 0b111)
}

func (p Piece) String() string {
	pt := p.Type().String()
	if p.Colour() == Black {
		return strings.ToLower(pt)
	}
	return pt
}
