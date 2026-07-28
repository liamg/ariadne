package board

import (
	"math/rand/v2"
	"testing"
)

func TestSlowRayAttacks(t *testing.T) {
	tests := []struct {
		name       string
		sq         Square
		occupancy  Bitboard
		directions []directionFunc
		expected   Bitboard
	}{
		{
			name:       "bishop on empty board",
			sq:         E4,
			occupancy:  EmptyBitboard.Set(E4),
			directions: bishopDirections,
			expected:   BitboardFromSquares(D5, C6, B7, A8, F5, G6, H7, D3, C2, B1, F3, G2, H1),
		},
		{
			name:       "rook on empty board",
			sq:         E4,
			occupancy:  EmptyBitboard.Set(E4),
			directions: rookDirections,
			expected:   (FileEMask | Rank4Mask) &^ E4.Bitboard(),
		},
		{
			name:       "queen on empty board",
			sq:         E4,
			occupancy:  EmptyBitboard,
			directions: queenDirections,
			expected:   BitboardFromSquares(D5, C6, B7, A8, F5, G6, H7, D3, C2, B1, F3, G2, H1) | (FileEMask|Rank4Mask)&^E4.Bitboard(),
		},
		{
			name:       "bishop blocked on all sides",
			sq:         E4,
			occupancy:  FullBitboard,
			directions: bishopDirections,
			expected:   BitboardFromSquares(D3, F5, F3, D5),
		},
		{
			name:       "rook blocked on all sides",
			sq:         E4,
			occupancy:  FullBitboard,
			directions: rookDirections,
			expected:   BitboardFromSquares(D4, F4, E3, E5),
		},
		{
			name:       "rook blocked last square",
			sq:         A1,
			occupancy:  EmptyBitboard.Set(A1).Set(A8),
			directions: rookDirections,
			expected:   (FileAMask | Rank1Mask) &^ A1.Bitboard(),
		},
		{
			name:       "bishop blocked on last square",
			sq:         E4,
			occupancy:  EmptyBitboard.Set(A8),
			directions: bishopDirections,
			expected:   BitboardFromSquares(D5, C6, B7, A8, F5, G6, H7, D3, C2, B1, F3, G2, H1),
		},
		{
			name:       "rook on west edge",
			sq:         A4,
			occupancy:  EmptyBitboard,
			directions: rookDirections,
			expected:   (FileAMask | Rank4Mask) &^ A4.Bitboard(),
		},
		{
			name:       "rook on east edge",
			sq:         H4,
			occupancy:  EmptyBitboard,
			directions: rookDirections,
			expected:   (FileHMask | Rank4Mask) &^ H4.Bitboard(),
		},
		{
			name:       "rook on north edge",
			sq:         D8,
			occupancy:  EmptyBitboard,
			directions: rookDirections,
			expected:   (FileDMask | Rank8Mask) &^ D8.Bitboard(),
		},
		{
			name:       "rook on south edge",
			sq:         D1,
			occupancy:  EmptyBitboard,
			directions: rookDirections,
			expected:   (FileDMask | Rank1Mask) &^ D1.Bitboard(),
		},
		{
			name:       "bishop in nw corner",
			sq:         A8,
			occupancy:  EmptyBitboard,
			directions: bishopDirections,
			expected:   BitboardFromSquares(B7, C6, D5, E4, F3, G2, H1),
		},
		{
			name:       "bishop in se corner",
			sq:         H1,
			occupancy:  EmptyBitboard,
			directions: bishopDirections,
			expected:   BitboardFromSquares(G2, F3, E4, D5, C6, B7, A8),
		},
		{
			name:       "bishop in sw corner",
			sq:         A1,
			occupancy:  EmptyBitboard,
			directions: bishopDirections,
			expected:   BitboardFromSquares(B2, C3, D4, E5, F6, G7, H8),
		},
		{
			name:       "bishop in ne corner",
			sq:         H8,
			occupancy:  EmptyBitboard,
			directions: bishopDirections,
			expected:   BitboardFromSquares(G7, F6, E5, D4, C3, B2, A1),
		},
		{
			name:       "bishop on a4",
			sq:         A4,
			occupancy:  EmptyBitboard,
			directions: bishopDirections,
			expected:   BitboardFromSquares(B5, C6, D7, E8, B3, C2, D1),
		},
		{
			name:       "bishop on h5",
			sq:         H5,
			occupancy:  EmptyBitboard,
			directions: bishopDirections,
			expected:   BitboardFromSquares(G6, F7, E8, G4, F3, E2, D1),
		},
		{
			name:       "rook attack blocked halfway",
			sq:         D4,
			occupancy:  EmptyBitboard.Set(D4).Set(F4),
			directions: rookDirections,
			expected:   BitboardFromSquares(D1, D2, D3, D5, D6, D7, D8, C4, B4, A4, E4, F4),
		},
		{
			name:       "bishop attack blocked halfway",
			sq:         D4,
			occupancy:  EmptyBitboard.Set(D4).Set(G7),
			directions: bishopDirections,
			expected:   BitboardFromSquares(C5, B6, A7, E5, F6, G7, C3, B2, A1, E3, F2, G1),
		},
		{
			name:       "rook attack blocked halfway - two blockers",
			sq:         D4,
			occupancy:  EmptyBitboard.Set(D4).Set(F4).Set(G4),
			directions: rookDirections,
			expected:   BitboardFromSquares(D1, D2, D3, D5, D6, D7, D8, C4, B4, A4, E4, F4),
		},
		{
			name:       "bishop attack blocked halfway - two blockers",
			sq:         D4,
			occupancy:  EmptyBitboard.Set(D4).Set(G7).Set(H8),
			directions: bishopDirections,
			expected:   BitboardFromSquares(C5, B6, A7, E5, F6, G7, C3, B2, A1, E3, F2, G1),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := slowlyFindRayAttacks(test.sq, test.occupancy, test.directions)
			if result != test.expected {
				t.Errorf("slowlyFindRayAttacks(%v, %v, [...]) = %v; expected %v", test.sq, test.occupancy, result, test.expected)
			}
		})
	}
}

func TestSlowRayAttacks_RookAlwaysAttacks14SquaresOnEmptyBoard(t *testing.T) {
	for sq := A1; sq <= H8; sq++ {
		attacks := slowlyFindRayAttacks(sq, EmptyBitboard.Set(sq), rookDirections)
		if attacks.Count() != 14 {
			t.Errorf("Rook attacks from square %v on empty board: got %v squares, expected 14", sq, attacks.Count())
		}
	}
}

func TestSlowRayAttacks_RookAttackSymmetry(t *testing.T) {
	for sq := A1; sq <= H8; sq++ {
		attacks := slowlyFindRayAttacks(sq, EmptyBitboard.Set(sq), rookDirections)
		for {
			var asq Square
			asq, attacks = attacks.PopSquare()
			if asq == NoSquare {
				break
			}
			reverseAttacks := slowlyFindRayAttacks(asq, EmptyBitboard.Set(asq), rookDirections)
			if !reverseAttacks.Has(sq) {
				t.Errorf("Rook attack symmetry failed: square %v attacks %v, but reverse does not", sq, asq)
			}
		}
	}
}

func TestSlowRayAttacks_BishopAttackSymmetry(t *testing.T) {
	for sq := A1; sq <= H8; sq++ {
		attacks := slowlyFindRayAttacks(sq, EmptyBitboard.Set(sq), bishopDirections)
		for {
			var asq Square
			asq, attacks = attacks.PopSquare()
			if asq == NoSquare {
				break
			}
			reverseAttacks := slowlyFindRayAttacks(asq, EmptyBitboard.Set(asq), bishopDirections)
			if !reverseAttacks.Has(sq) {
				t.Errorf("Bishop attack symmetry failed: square %v attacks %v, but reverse does not", sq, asq)
			}
		}
	}
}

func TestSlowRayAttacks_Fuzz(t *testing.T) {
	bishopVectors := [][2]int{
		{1, 1},
		{1, -1},
		{-1, 1},
		{-1, -1},
	}

	rookVectors := [][2]int{
		{1, 0},
		{-1, 0},
		{0, 1},
		{0, -1},
	}

	queenVectors := make([][2]int, 0, len(bishopVectors)+len(rookVectors))
	queenVectors = append(queenVectors, bishopVectors...)
	queenVectors = append(queenVectors, rookVectors...)

	tests := []struct {
		name       string
		directions []directionFunc
		vectors    [][2]int
	}{
		{
			name:       "queen",
			directions: queenDirections,
			vectors:    queenVectors,
		},
		{
			name:       "rook",
			directions: rookDirections,
			vectors:    rookVectors,
		},
		{
			name:       "bishop",
			directions: bishopDirections,
			vectors:    bishopVectors,
		},
	}

	oracle := func(from Square, occupancy Bitboard, vecs [][2]int) Bitboard {
		output := EmptyBitboard

		for _, vec := range vecs {
			rank := from.Rank()
			file := from.File()
			for {
				file += File(vec[0])
				rank += Rank(vec[1])
				sq, err := SquareFromFileAndRank(file, rank)
				if err != nil {
					break
				}
				output |= sq.Bitboard()
				if occupancy.Has(sq) {
					break
				}
			}
		}

		return output
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rnd := rand.New(rand.NewPCG(1, 2))
			for sparsity := range 5 {
				for range 1000 {
					bb := Bitboard(rnd.Uint64())
					for range sparsity {
						bb &= Bitboard(rnd.Uint64())
					}
					for sq := A1; sq <= H8; sq++ {
						attacks := slowlyFindRayAttacks(sq, bb, test.directions)
						expected := oracle(sq, bb, test.vectors)
						if attacks != expected {
							t.Fatalf("Fuzz test for %q failed for square %v, occupancy %v: got %v, expected %v", test.name, sq, bb, attacks, expected)
						}
					}
				}
			}
		})
	}
}
