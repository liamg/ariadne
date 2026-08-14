package board

import (
	"math/bits"
	"strings"
)

// Bitboard is a 64-bit integer where each bit represents a square on the chess board.
// The least significant bit (LSB) represents square A1, and the most significant bit (MSB) represents square H8.
type Bitboard uint64

const (
	// EmptyBitboard represents a bitboard with no pieces on it.
	EmptyBitboard Bitboard = 0
	// FullBitboard represents a bitboard with all squares occupied.
	FullBitboard Bitboard = ^Bitboard(0)
)

// Has returns true if the bitboard has a piece on the given square.
func (b Bitboard) Has(sq Square) bool {
	return b&sq.Bitboard() != 0
}

// Set sets the bit for the given square in the bitboard.
func (b Bitboard) Set(sq Square) Bitboard {
	return b | sq.Bitboard()
}

// Clear clears the bit for the given square in the bitboard.
func (b Bitboard) Clear(sq Square) Bitboard {
	return b &^ sq.Bitboard()
}

// Count returns the number of bits (squares) set in the bitboard.
func (b Bitboard) Count() int {
	return bits.OnesCount64(uint64(b))
}

func BitboardFromSquares(squares ...Square) Bitboard {
	var b Bitboard
	for _, sq := range squares {
		b = b.Set(sq)
	}
	return b
}

func (b Bitboard) EachSquareSlow(f func(sq Square)) {
	e := b
	var sq Square
	for {
		sq, e = e.PopSquare()
		if sq == NoSquare {
			break
		}
		f(sq)
	}
}

// String returns a string representation of the bitboard, with 'x' for occupied squares and '.' for empty squares.
func (b Bitboard) String() string {
	var sb strings.Builder
	sb.WriteString("\n")
	for rank := Rank8; rank >= Rank1; rank-- {
		sb.WriteString(rank.String())
		sb.WriteString(" ")
		for file := FileA; file <= FileH; file++ {
			sq, _ := SquareFromFileAndRank(file, rank)
			if b.Has(sq) {
				sb.WriteString("x ")
			} else {
				sb.WriteString(". ")
			}
		}
		sb.WriteString("\n")
	}
	sb.WriteString("  a b c d e f g h \n")
	return sb.String()
}

// PopSquare returns the least significant square in the bitboard and a new bitboard with that square cleared.
func (b Bitboard) PopSquare() (Square, Bitboard) {
	sq := Square(bits.TrailingZeros64(uint64(b)))
	return sq, b & (b - 1)
}

// North shifts all set squares north one square
func (b Bitboard) North() Bitboard {
	return b << 8
}

// South shifts all set squares south one square
func (b Bitboard) South() Bitboard {
	return b >> 8
}

// East shifts all set squares east one square, without wrapping
func (b Bitboard) East() Bitboard {
	return (b &^ FileHMask) << 1
}

// West shifts all set squares west one square, without wrapping
func (b Bitboard) West() Bitboard {
	return (b &^ FileAMask) >> 1
}

// NorthEast shifts all set squares northeast one square, without wrapping
func (b Bitboard) NorthEast() Bitboard {
	return (b &^ FileHMask) << 9
}

// NorthWest shifts all set squares northwest one square, without wrapping
func (b Bitboard) NorthWest() Bitboard {
	return (b &^ FileAMask) << 7
}

// SouthEast shifts all set squares southeast one square, without wrapping
func (b Bitboard) SouthEast() Bitboard {
	return (b &^ FileHMask) >> 7
}

// SouthWest shifts all set squares southwest one square, without wrapping
func (b Bitboard) SouthWest() Bitboard {
	return (b &^ FileAMask) >> 9
}

func (b Bitboard) FlipVertical() Bitboard {
	return Bitboard(bits.ReverseBytes64(uint64(b)))
}
