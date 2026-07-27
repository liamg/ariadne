package board

import "strings"

type Position struct {
	byType   [7]Bitboard // Bitboards for each piece type (Pawn, Knight, Bishop, Rook, Queen, King)
	byColour [2]Bitboard // Bitboards for each color (White, Black)
}

func StartingPosition() *Position {
	return &Position{
		byType: [7]Bitboard{
			0x0000000000000000,
			0x00FF00000000FF00, // Pawns
			0x4200000000000042, // Knights
			0x2400000000000024, // Bishops
			0x8100000000000081, // Rooks
			0x0800000000000008, // Queens
			0x1000000000000010, // Kings
		},
		byColour: [2]Bitboard{
			0x000000000000FFFF, // White pieces
			0xFFFF000000000000, // Black pieces
		},
	}
}

func (p *Position) String() string {
	var sb strings.Builder
	sb.WriteString("\n")
	for rank := Rank8; rank >= Rank1; rank-- {
		sb.WriteString(rank.String())
		sb.WriteString(" ")
		for file := FileA; file <= FileH; file++ {
			sq, _ := SquareFromFileAndRank(file, rank)
			var piece Piece
			for pieceType := PieceTypePawn; pieceType <= PieceTypeKing; pieceType++ {
				if p.byType[pieceType].Has(sq) {
					if p.byColour[Black].Has(sq) {
						piece = NewPiece(Black, pieceType)
					} else {
						piece = NewPiece(White, pieceType)
					}
					break
				}
			}
			if piece.Type() == PieceTypeNone {
				sb.WriteString(".")
			} else {
				sb.WriteString(piece.String())
			}
			sb.WriteString(" ")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("  a b c d e f g h \n")
	return sb.String()
}

func (p *Position) PiecesByColour(colour Colour) Bitboard {
	return p.byColour[colour]
}

func (p *Position) Pieces(colour Colour, pieceType PieceType) Bitboard {
	return p.byType[pieceType] & p.byColour[colour]
}

func (p *Position) Occupancy() Bitboard {
	return p.byColour[White] | p.byColour[Black]
}
