package board

import (
	"math/rand/v2"
	"testing"
)

func TestLookupInvariants(t *testing.T) {
	tests := []struct {
		name   string
		slider slider
	}{
		{
			name:   "rook",
			slider: rookSlider,
		},
		{
			name:   "bishop",
			slider: bishopSlider,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for i, entry := range test.slider.magicEntries {
				sq := Square(i)
				if entry.magic == 0 {
					t.Errorf("magic entry for square %v has magic 0", sq)
				}
				if entry.shift != 64-entry.mask.Count() {
					t.Errorf("magic entry for square %v has shift %d; expected %d", sq, entry.shift, 64-entry.mask.Count())
				}
				if len(entry.attacks) != 1<<entry.mask.Count() {
					t.Errorf("magic entry for square %v has %d attacks; expected %d", sq, len(entry.attacks), 1<<entry.mask.Count())
				}
			}
		})
	}
}

func TestLookupStartingPosition(t *testing.T) {
	occ := StartingPosition().Occupancy()

	tests := []struct {
		name     string
		sq       Square
		expected Bitboard
		lookup   func(sq Square, occ Bitboard) Bitboard
	}{
		{
			"bishop on C1",
			C1,
			BitboardFromSquares(B2, D2),
			bishopLookup,
		},
		{
			"bishop on D4",
			D4,
			BitboardFromSquares(B2, F2, C3, E3, C5, E5, B6, F6, A7, G7),
			bishopLookup,
		},
		{
			"rook on A1",
			A1,
			BitboardFromSquares(B1, A2),
			rookLookup,
		},
		{
			"rook on H8",
			H8,
			BitboardFromSquares(G8, H7),
			rookLookup,
		},
		{
			"rook on D4",
			D4,
			BitboardFromSquares(D2, D3, D5, D6, D7, A4, B4, C4, E4, F4, G4, H4),
			rookLookup,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.lookup(test.sq, occ)
			if actual != test.expected {
				t.Errorf("lookup(%v, occ) = %v; want %v", test.sq, actual, test.expected)
			}
		})
	}
}

func TestLookupWithEnumeratedOccupancies(t *testing.T) {
	tests := []struct {
		name          string
		directions    []directionFunc
		lookup        func(sq Square, occ Bitboard) Bitboard
		expectedCount int
	}{
		{
			"rook",
			rookDirections,
			rookLookup,
			102400,
		},
		{
			"bishop",
			bishopDirections,
			bishopLookup,
			5248,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			total := 0
			for sq := A1; sq <= H8; sq++ {
				for _, pair := range computeOccupancyAttacks(sq, test.directions) {
					actual := test.lookup(sq, pair.occupancy)
					expected := pair.attacks
					if actual != expected {
						t.Fatalf("lookup(%v, %v) = %v; want %v", sq, pair.occupancy, actual, expected)
					}
					total++
				}
			}
			if total != test.expectedCount {
				t.Fatalf("Total lookups for %v = %d; want %d", test.name, total, test.expectedCount)
			}
		})
	}
}

func TestLookupFuzz(t *testing.T) {
	tests := []struct {
		name       string
		directions []directionFunc
		lookup     func(sq Square, occ Bitboard) Bitboard
	}{
		{
			"rook",
			rookDirections,
			rookLookup,
		},
		{
			"bishop",
			bishopDirections,
			bishopLookup,
		},
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
						expected := slowlyFindRayAttacks(sq, bb, test.directions)
						actual := test.lookup(sq, bb)
						if actual != expected {
							t.Fatalf("lookup(%v, %v) = %v; want %v", sq, bb, actual, expected)
						}
					}
				}
			}
		})
	}
}
