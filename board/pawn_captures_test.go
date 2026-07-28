package board

import (
	"math/rand/v2"
	"testing"
)

func TestFindPawnCaptures(t *testing.T) {
	tests := []struct {
		name         string
		colour       Colour
		pawns        Bitboard
		enemies      Bitboard
		expectedWest Bitboard
		expectedEast Bitboard
	}{
		{
			name:         "capture does not wrap west (white)",
			colour:       White,
			pawns:        EmptyBitboard.Set(A2),
			enemies:      FullBitboard.Clear(A2),
			expectedWest: EmptyBitboard,
			expectedEast: EmptyBitboard.Set(B3),
		},
		{
			name:         "capture does not wrap west (black)",
			colour:       Black,
			pawns:        EmptyBitboard.Set(A7),
			enemies:      FullBitboard.Clear(A7),
			expectedWest: EmptyBitboard,
			expectedEast: EmptyBitboard.Set(B6),
		},
		{
			name:         "capture does not wrap east (white)",
			colour:       White,
			pawns:        EmptyBitboard.Set(H2),
			enemies:      FullBitboard.Clear(H2),
			expectedWest: EmptyBitboard.Set(G3),
			expectedEast: EmptyBitboard,
		},
		{
			name:         "capture does not wrap east (black)",
			colour:       Black,
			pawns:        EmptyBitboard.Set(H7),
			enemies:      FullBitboard.Clear(H7),
			expectedWest: EmptyBitboard.Set(G6),
			expectedEast: EmptyBitboard,
		},
		{
			name:         "two pawns target same enemy (white)",
			colour:       White,
			pawns:        EmptyBitboard.Set(A2).Set(C2),
			enemies:      EmptyBitboard.Set(B3),
			expectedWest: EmptyBitboard.Set(B3),
			expectedEast: EmptyBitboard.Set(B3),
		},
		{
			name:         "two pawns target same enemy (black)",
			colour:       Black,
			pawns:        EmptyBitboard.Set(A7).Set(C7),
			enemies:      EmptyBitboard.Set(B6),
			expectedWest: EmptyBitboard.Set(B6),
			expectedEast: EmptyBitboard.Set(B6),
		},
		{
			name:         "enemy in front",
			colour:       White,
			pawns:        EmptyBitboard.Set(B2),
			enemies:      EmptyBitboard.Set(B3),
			expectedWest: EmptyBitboard,
			expectedEast: EmptyBitboard,
		},
		{
			name:         "capture onto promotion row",
			colour:       White,
			pawns:        EmptyBitboard.Set(G7),
			enemies:      EmptyBitboard.Set(H8),
			expectedWest: EmptyBitboard,
			expectedEast: EmptyBitboard.Set(H8),
		},
		{
			name:         "starting position (white)",
			colour:       White,
			pawns:        Rank2Mask,
			enemies:      Rank7Mask | Rank8Mask,
			expectedWest: EmptyBitboard,
			expectedEast: EmptyBitboard,
		},
		{
			name:         "starting position (black)",
			colour:       Black,
			pawns:        Rank7Mask,
			enemies:      Rank1Mask | Rank2Mask,
			expectedWest: EmptyBitboard,
			expectedEast: EmptyBitboard,
		},
		{
			name:         "no enemies",
			colour:       White,
			pawns:        FullBitboard,
			enemies:      EmptyBitboard,
			expectedWest: EmptyBitboard,
			expectedEast: EmptyBitboard,
		},
		{
			name:         "no pawns",
			colour:       White,
			pawns:        EmptyBitboard,
			enemies:      FullBitboard,
			expectedWest: EmptyBitboard,
			expectedEast: EmptyBitboard,
		},
		{
			name:         "one pawn, two targets",
			colour:       White,
			pawns:        EmptyBitboard.Set(E4),
			enemies:      EmptyBitboard.Set(D5).Set(F5),
			expectedWest: EmptyBitboard.Set(D5),
			expectedEast: EmptyBitboard.Set(F5),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			west, east := findPawnCaptures(test.colour, test.pawns, test.enemies)
			if west != test.expectedWest {
				t.Errorf("Expected west captures %v, got %v", test.expectedWest, west)
			}
			if east != test.expectedEast {
				t.Errorf("Expected east captures %v, got %v", test.expectedEast, east)
			}
		})
	}
}

func TestPawnCapturesFuzzAgainstSlowImpl(t *testing.T) {
	slowImpl := func(colour Colour, pawns Bitboard, enemyOccupancy Bitboard) (west, east Bitboard) {
		var sq Square
		for {
			sq, pawns = pawns.PopSquare()
			if sq == NoSquare {
				break
			}
			switch colour {
			case White:
				rank := sq.Rank()
				file := sq.File()
				target, err := SquareFromFileAndRank(file+1, rank+1)
				if err == nil {
					capture := target.Bitboard()
					if enemyOccupancy&capture != 0 {
						east |= capture
					}
				}
				target, err = SquareFromFileAndRank(file-1, rank+1)
				if err == nil {
					capture := target.Bitboard()
					if enemyOccupancy&capture != 0 {
						west |= capture
					}
				}
			case Black:
				rank := sq.Rank()
				file := sq.File()
				target, err := SquareFromFileAndRank(file+1, rank-1)
				if err == nil {
					capture := target.Bitboard()
					if enemyOccupancy&capture != 0 {
						east |= capture
					}
				}
				target, err = SquareFromFileAndRank(file-1, rank-1)
				if err == nil {
					capture := target.Bitboard()
					if enemyOccupancy&capture != 0 {
						west |= capture
					}
				}

			}
		}
		return west, east
	}

	rnd := rand.New(rand.NewPCG(1, 2))

	var colour Colour
	for sparsity := range 3 {
		for range 1000 {
			colour = colour.Opposite()
			pawns := Bitboard(rnd.Uint64() & rnd.Uint64() & rnd.Uint64()) // always realistic pawn density
			enemyOccupancy := Bitboard(rnd.Uint64())
			for range sparsity {
				enemyOccupancy &= Bitboard(rnd.Uint64()) // reduce density
			}
			enemyOccupancy &= ^pawns // pawns cannot be occupied by enemy pieces

			actualWest, actualEast := findPawnCaptures(colour, pawns, enemyOccupancy)
			expectedWest, expectedEast := slowImpl(colour, pawns, enemyOccupancy)

			if actualWest != expectedWest {
				t.Fatalf("Fuzz test failed for colour %v, pawns %s, enemy occupancy %s: expected west captures %s, got %s", colour, pawns, enemyOccupancy, expectedWest, actualWest)
			}
			if actualEast != expectedEast {
				t.Fatalf("Fuzz test failed for colour %v, pawns %s, enemy occupancy %s: expected east captures %s, got %s", colour, pawns, enemyOccupancy, expectedEast, actualEast)
			}
		}
	}
}
