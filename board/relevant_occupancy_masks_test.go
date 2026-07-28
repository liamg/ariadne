package board

import (
	"math/rand/v2"
	"testing"
)

func TestRelevantOccupancyMasksHandPicked(t *testing.T) {
	tests := []struct {
		name       string
		sq         Square
		directions []directionFunc
		expected   Bitboard
	}{
		{
			"rook interior",
			D4,
			rookDirections,
			BitboardFromSquares(B4, C4, D2, D3, D5, D6, D7, E4, F4, G4),
		},
		{
			"rook corner",
			A1,
			rookDirections,
			BitboardFromSquares(B1, C1, D1, E1, F1, G1, A2, A3, A4, A5, A6, A7),
		},
		{
			"rook rim",
			A4,
			rookDirections,
			BitboardFromSquares(A2, A3, A5, A6, A7, B4, C4, D4, E4, F4, G4),
		},
		{
			"rook first rank",
			E1,
			rookDirections,
			BitboardFromSquares(B1, C1, D1, F1, G1, E2, E3, E4, E5, E6, E7),
		},
		{
			"bishop interior",
			D4,
			bishopDirections,
			BitboardFromSquares(B2, C3, E3, F2, C5, E5, B6, F6, G7),
		},
		{
			"bishop corner",
			A1,
			bishopDirections,
			BitboardFromSquares(B2, C3, D4, E5, F6, G7),
		},
		{
			"bishop rim",
			A2,
			bishopDirections,
			BitboardFromSquares(B3, C4, D5, E6, F7),
		},
		{
			"bishop two blocked rays",
			C7,
			bishopDirections,
			BitboardFromSquares(B6, D6, E5, F4, G3),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mask := computeRelevantOccupancyMask(test.sq, test.directions)
			if mask != test.expected {
				t.Errorf("For square %v, expected mask %v, got %v", test.sq, test.expected, mask)
			}
		})
	}
}

func TestRelevantOccupancyMasksAgainstKnownBitCounts(t *testing.T) {
	for sq := A1; sq <= H8; sq++ {
		file := sq.File()
		rank := sq.Rank()

		isCorner := (file == FileA || file == FileH) && (rank == Rank1 || rank == Rank8)
		isEdge := (file == FileA || file == FileH || rank == Rank1 || rank == Rank8) && !isCorner

		rookMask := rookRelevantOccupancy[sq]
		rookAttacks := slowlyFindRayAttacks(sq, EmptyBitboard, rookDirections)

		expectedRookCount := 0

		switch {
		case isCorner:
			expectedRookCount = 12
		case isEdge:
			expectedRookCount = 11
		default:
			expectedRookCount = 10
		}

		if rookMask.Count() != expectedRookCount {
			t.Errorf("For square %v, expected rook mask count %d, got %d", sq, expectedRookCount, rookMask.Count())
		}

		if rookMask.Has(sq) {
			t.Errorf("For square %v, rook mask should not include the square itself", sq)
		}

		if rookMask&^rookAttacks != EmptyBitboard {
			t.Errorf("For square %v, rook mask has squares not in attacks: %v", sq, rookMask & ^rookAttacks)
		}

		bishopMask := bishopRelevantOccupancy[sq]
		bishopAttacks := slowlyFindRayAttacks(sq, EmptyBitboard, bishopDirections)

		expectedBishopCount := 0

		switch {
		case isCorner:
			expectedBishopCount = 6
		case (file == FileD || file == FileE) && (rank == Rank4 || rank == Rank5): // central 4
			expectedBishopCount = 9
		case file >= FileC && file <= FileF && rank >= Rank3 && rank <= Rank6: // inner 16
			expectedBishopCount = 7
		default:
			expectedBishopCount = 5
		}

		if bishopMask.Count() != expectedBishopCount {
			t.Errorf("For square %v, expected bishop mask count %d, got %d", sq, expectedBishopCount, bishopMask.Count())
		}

		if bishopMask.Has(sq) {
			t.Errorf("For square %v, bishop mask should not include the square itself", sq)
		}

		if bishopMask&^bishopAttacks != EmptyBitboard {
			t.Errorf("For square %v, bishop mask has squares not in attacks: %v", sq, bishopMask & ^bishopAttacks)
		}
	}
}

func TestRelevantOccupancyMasksFuzz(t *testing.T) {
	directionSets := []struct {
		directions []directionFunc
		table      [64]Bitboard
	}{
		{rookDirections, rookRelevantOccupancy},
		{bishopDirections, bishopRelevantOccupancy},
	}

	rnd := rand.New(rand.NewPCG(1, 2))
	for _, dirSet := range directionSets {
		for sparsity := range 5 {
			for range 1000 {
				bb := Bitboard(rnd.Uint64())
				for range sparsity {
					bb &= Bitboard(rnd.Uint64())
				}
				for sq := A1; sq <= H8; sq++ {
					attacks := slowlyFindRayAttacks(sq, bb, dirSet.directions)
					masked := slowlyFindRayAttacks(sq, bb&dirSet.table[sq], dirSet.directions)
					if attacks != masked {
						t.Fatalf("Fuzz test failed for square %v, occupancy %v: got %v, expected %v", sq, bb, masked, attacks)
					}
				}
			}
		}
	}
}
