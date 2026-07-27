package board

import (
	"math/bits"
	"testing"
)

func TestFindEnPassantAttackers(t *testing.T) {
	tests := []struct {
		name     string
		colour   Colour
		pawns    Bitboard
		epSquare Square
		expected Bitboard
	}{
		{
			name:     "no en passant square (white)",
			colour:   White,
			pawns:    EmptyBitboard.Set(D5),
			epSquare: NoSquare,
			expected: EmptyBitboard,
		},
		{
			name:     "no en passant square (black)",
			colour:   Black,
			pawns:    EmptyBitboard.Set(D4),
			epSquare: NoSquare,
			expected: EmptyBitboard,
		},
		{
			name:     "white pawn can capture en passant",
			colour:   White,
			pawns:    EmptyBitboard.Set(D5),
			epSquare: E6,
			expected: EmptyBitboard.Set(D5),
		},
		{
			name:     "black pawn can capture en passant",
			colour:   Black,
			pawns:    EmptyBitboard.Set(D4),
			epSquare: E3,
			expected: EmptyBitboard.Set(D4),
		},
		{
			name:     "white pawns everywhere - only two in attacking position can capture",
			colour:   White,
			pawns:    FullBitboard.Clear(E6),
			epSquare: E6,
			expected: EmptyBitboard.Set(D5).Set(F5),
		},
		{
			name:     "black pawns everywhere - only two in attacking position can capture",
			colour:   Black,
			pawns:    FullBitboard.Clear(E3),
			epSquare: E3,
			expected: EmptyBitboard.Set(D4).Set(F4),
		},
		{
			name:     "white pawn on wrong side of en passant square",
			colour:   White,
			pawns:    EmptyBitboard.Set(D7),
			epSquare: E6,
			expected: EmptyBitboard,
		},
		{
			name:     "black pawn on wrong side of en passant square",
			colour:   Black,
			pawns:    EmptyBitboard.Set(D2),
			epSquare: E3,
			expected: EmptyBitboard,
		},
		{
			name:     "pawn on right file, wrong rank",
			colour:   White,
			pawns:    EmptyBitboard.Set(D4),
			epSquare: E6,
			expected: EmptyBitboard,
		},
		{
			name:     "edge file",
			colour:   White,
			pawns:    EmptyBitboard.Set(B5).Set(H5),
			epSquare: A6,
			expected: EmptyBitboard.Set(B5),
		},
		{
			name:     "no pawns",
			colour:   White,
			pawns:    EmptyBitboard,
			epSquare: E6,
			expected: EmptyBitboard,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := findEnPassantAttackers(test.colour, test.pawns, test.epSquare)
			if result != test.expected {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
			if result&^test.pawns != EmptyBitboard {
				t.Errorf("Result contains pawns that were not in the input pawns bitboard. Result: %v, Input pawns: %v", result, test.pawns)
			}
			if bits.OnesCount64(uint64(result)) > 2 {
				t.Errorf("Result contains more than 2 pawns, which is impossible. Result: %v", result)
			}
		})
	}
}

func TestFindEnPassantAttackersAllPossibilities(t *testing.T) {
	check := func(t *testing.T, colour Colour, pawns Bitboard, ep Square, expected Bitboard) {
		t.Helper()
		result := findEnPassantAttackers(colour, pawns, ep)
		if result != expected {
			t.Errorf("Expected %v, got %v when pawns are %v for square %s (%s)", expected, result, pawns, ep, colour)
		}
		if result&^pawns != EmptyBitboard {
			t.Errorf("Result contains pawns that were not in the input pawns bitboard. Result: %v, Input pawns: %v", result, pawns)
		}
		if bits.OnesCount64(uint64(result)) > 2 {
			t.Errorf("Result contains more than 2 pawns, which is impossible. Result: %v", result)
		}
	}

	for _, colour := range []Colour{White, Black} {
		for file := FileA; file <= FileH; file++ {

			epRank := Rank6
			pawnRank := Rank5
			if colour == Black {
				epRank = Rank3
				pawnRank = Rank4
			}
			epSq := mustSquareFromFileAndRank(file, epRank)

			if file > FileA {
				// pawns attack from east and west
				if file < FileH {
					check(
						t,
						colour,
						FullBitboard.Clear(epSq),
						epSq,
						EmptyBitboard.
							Set(mustSquareFromFileAndRank(file-1, pawnRank)).
							Set(mustSquareFromFileAndRank(file+1, pawnRank)),
					)
				} else {
					check(
						t,
						colour,
						FullBitboard.Clear(epSq),
						epSq,
						EmptyBitboard.Set(mustSquareFromFileAndRank(file-1, pawnRank)),
					)
				}
			} else if file < FileH {
				check(
					t,
					colour,
					FullBitboard.Clear(epSq),
					epSq,
					EmptyBitboard.Set(mustSquareFromFileAndRank(file+1, pawnRank)),
				)
			}

		}
	}
}
