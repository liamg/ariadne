package board

import "fmt"

// Square represents a square on the chess board.
// Six bits used. Lower 3 bits are file, next 3 bits are rank.
// NOTE: no benefit to making it smaller like a uint8, creates footguns like wrapping,
// and occasionally costs extra truncation instructions.
// NOTE: when indexing using a Square, it's a good idea to do array[square&63] so the
// compiler can optimize the bounds check away.
type Square int

const (
	A1 Square = iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1
	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2
	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3
	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4
	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5
	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6
	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7
	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
	NoSquare
)

// File returns the file of the square (A-H)
func (s Square) File() File {
	// mask off the top 3 bits is effectively mod 8
	return File(s & 7)
}

// Rank returns the rank of the square (1-8)
// NOTE: the actual backing value of a Rank is 0-7, not 1-8, so be careful when using it as an index.
func (s Square) Rank() Rank {
	// shift 3 bits down is effectively divide by 8
	return Rank(s >> 3)
}

func (s Square) String() string {
	if s < A1 || s > H8 {
		if s == NoSquare {
			return "-"
		}
		return "?"
	}
	return s.File().String() + s.Rank().String()
}

func (s Square) Bitboard() Bitboard {
	return 1 << s
}

var ErrInvalidSquare = fmt.Errorf("invalid square")

func ParseSquare(s string) (Square, error) {
	if len(s) != 2 {
		return NoSquare, ErrInvalidSquare
	}
	file := s[0]
	rank := s[1]

	if file < 'a' || file > 'h' || rank < '1' || rank > '8' {
		return NoSquare, ErrInvalidSquare
	}

	return Square((file - 'a') + (rank-'1')*8), nil
}

func SquareFromFileAndRank(file File, rank Rank) (Square, error) {
	if file < FileA || file > FileH || rank < Rank1 || rank > Rank8 {
		return NoSquare, ErrInvalidSquare
	}

	return Square(file) + Square(rank*8), nil
}

// File is a vertical column of the board, lettered A-H left to right from White's perspective.
type File int

const (
	FileA File = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH
)

const (
	FileAMask Bitboard = 0b0000000100000001000000010000000100000001000000010000000100000001
	FileBMask Bitboard = 0b0000001000000010000000100000001000000010000000100000001000000010
	FileCMask Bitboard = 0b0000010000000100000001000000010000000100000001000000010000000100
	FileDMask Bitboard = 0b0000100000001000000010000000100000001000000010000000100000001000
	FileEMask Bitboard = 0b0001000000010000000100000001000000010000000100000001000000010000
	FileFMask Bitboard = 0b0010000000100000001000000010000000100000001000000010000000100000
	FileGMask Bitboard = 0b0100000001000000010000000100000001000000010000000100000001000000
	FileHMask Bitboard = 0b1000000010000000100000001000000010000000100000001000000010000000
)

func (f File) String() string {
	if f < FileA || f > FileH {
		return "?"
	}
	return string('a' + rune(f))
}

// Rank is a horizontal row of the board, numbered 1-8 from White's perspective.
// Rank1 is the first rank (White's home rank), and Rank8 is the eighth rank (Black's home rank).
// NOTE: the actual backing value of a Rank is 0-7, not 1-8, so be careful when using it as an index.
type Rank int

const (
	Rank1 Rank = iota
	Rank2
	Rank3
	Rank4
	Rank5
	Rank6
	Rank7
	Rank8
)

const (
	Rank1Mask Bitboard = 0b0000000000000000000000000000000000000000000000000000000011111111
	Rank2Mask Bitboard = 0b0000000000000000000000000000000000000000000000001111111100000000
	Rank3Mask Bitboard = 0b0000000000000000000000000000000000000000111111110000000000000000
	Rank4Mask Bitboard = 0b0000000000000000000000000000000011111111000000000000000000000000
	Rank5Mask Bitboard = 0b0000000000000000000000001111111100000000000000000000000000000000
	Rank6Mask Bitboard = 0b0000000000000000111111110000000000000000000000000000000000000000
	Rank7Mask Bitboard = 0b0000000011111111000000000000000000000000000000000000000000000000
	Rank8Mask Bitboard = 0b1111111100000000000000000000000000000000000000000000000000000000
)

func (r Rank) String() string {
	if r < Rank1 || r > Rank8 {
		return "?"
	}
	return string('1' + rune(r))
}
