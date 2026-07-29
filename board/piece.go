package board

import "strings"

// PieceType is the type of chess piece, i.e. one of pawn, knight, bishop, rook, queen, king.
// It does NOT encode the colour of the piece, which is stored separately in the Piece type.
type PieceType uint8

const (
	NoPieceType PieceType = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

var (
	WhiteKing   = NewPiece(White, King)
	BlackKing   = NewPiece(Black, King)
	WhiteQueen  = NewPiece(White, Queen)
	BlackQueen  = NewPiece(Black, Queen)
	WhiteRook   = NewPiece(White, Rook)
	BlackRook   = NewPiece(Black, Rook)
	WhiteBishop = NewPiece(White, Bishop)
	BlackBishop = NewPiece(Black, Bishop)
	WhiteKnight = NewPiece(White, Knight)
	BlackKnight = NewPiece(Black, Knight)
	WhitePawn   = NewPiece(White, Pawn)
	BlackPawn   = NewPiece(Black, Pawn)
	NoPiece     = NewPiece(White, NoPieceType)
)

func (p PieceType) String() string {
	switch p {
	case Pawn:
		return "P"
	case Knight:
		return "N"
	case Bishop:
		return "B"
	case Rook:
		return "R"
	case Queen:
		return "Q"
	case King:
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
	if pieceType == NoPieceType { // avoid encoding a black non-piece as >0
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
