package board

import (
	"math/rand/v2"
	"testing"
)

func TestPawnPushes(t *testing.T) {
	tests := []struct {
		name            string
		colour          Colour
		pawns           Bitboard
		occupancy       Bitboard
		expectedSingles Bitboard
		expectedDoubles Bitboard
	}{
		{
			name:            "starting position (white)",
			occupancy:       StartingPosition().Occupancy(),
			pawns:           StartingPosition().Pieces(White, PieceTypePawn),
			colour:          White,
			expectedSingles: Rank3Mask,
			expectedDoubles: Rank4Mask,
		},
		{
			name:            "starting position (black)",
			occupancy:       StartingPosition().Occupancy(),
			pawns:           StartingPosition().Pieces(Black, PieceTypePawn),
			colour:          Black,
			expectedSingles: Rank6Mask,
			expectedDoubles: Rank5Mask,
		},
		{
			name:            "blocked pawns (white)",
			occupancy:       EmptyBitboard.Set(E3).Set(D4),
			pawns:           EmptyBitboard.Set(E2).Set(D2),
			colour:          White,
			expectedSingles: EmptyBitboard.Set(D3),
			expectedDoubles: EmptyBitboard,
		},
		{
			name:            "blocked pawns (black)",
			occupancy:       EmptyBitboard.Set(E6).Set(D5),
			pawns:           EmptyBitboard.Set(E7).Set(D7),
			colour:          Black,
			expectedSingles: EmptyBitboard.Set(D6),
			expectedDoubles: EmptyBitboard,
		},
		{
			name:            "pawns without double available (white)",
			occupancy:       EmptyBitboard,
			pawns:           EmptyBitboard.Set(E3).Set(D3),
			colour:          White,
			expectedSingles: EmptyBitboard.Set(E4).Set(D4),
		},
		{
			name:            "pawns without double available (black)",
			occupancy:       EmptyBitboard,
			pawns:           EmptyBitboard.Set(E6).Set(D6),
			colour:          Black,
			expectedSingles: EmptyBitboard.Set(E5).Set(D5),
			expectedDoubles: EmptyBitboard,
		},
		{
			name:            "pawns on final rank (white)",
			occupancy:       EmptyBitboard,
			pawns:           EmptyBitboard.Set(E8).Set(D8),
			colour:          White,
			expectedSingles: EmptyBitboard,
			expectedDoubles: EmptyBitboard,
		},
		{
			name:            "pawns on final rank (black)",
			occupancy:       EmptyBitboard,
			pawns:           EmptyBitboard.Set(E1).Set(D1),
			colour:          Black,
			expectedSingles: EmptyBitboard,
			expectedDoubles: EmptyBitboard,
		},
		{
			name:            "no pawns (white)",
			occupancy:       EmptyBitboard,
			pawns:           EmptyBitboard,
			colour:          White,
			expectedSingles: EmptyBitboard,
			expectedDoubles: EmptyBitboard,
		},
		{
			name:            "no pawns (black)",
			occupancy:       EmptyBitboard,
			pawns:           EmptyBitboard,
			colour:          Black,
			expectedSingles: EmptyBitboard,
			expectedDoubles: EmptyBitboard,
		},
		{
			name:            "full occupancy (white)",
			occupancy:       FullBitboard,
			pawns:           FullBitboard,
			colour:          White,
			expectedSingles: EmptyBitboard,
			expectedDoubles: EmptyBitboard,
		},
		{
			name:            "full occupancy (black)",
			occupancy:       FullBitboard,
			pawns:           FullBitboard,
			colour:          Black,
			expectedSingles: EmptyBitboard,
			expectedDoubles: EmptyBitboard,
		},
		{
			name:            "pawns blocking on same file (white)",
			occupancy:       EmptyBitboard.Set(E2).Set(E3),
			pawns:           EmptyBitboard.Set(E2).Set(E3),
			colour:          White,
			expectedSingles: EmptyBitboard.Set(E4),
			expectedDoubles: EmptyBitboard,
		},
		{
			name:            "pawns blocking on same file (black)",
			occupancy:       EmptyBitboard.Set(E7).Set(E6),
			pawns:           EmptyBitboard.Set(E7).Set(E6),
			colour:          Black,
			expectedSingles: EmptyBitboard.Set(E5),
			expectedDoubles: EmptyBitboard,
		},
		{
			name:            "mixed pawns (white)",
			occupancy:       EmptyBitboard.Set(A2).Set(B3).Set(C2).Set(C3),
			pawns:           EmptyBitboard.Set(A2).Set(B3).Set(C2),
			colour:          White,
			expectedSingles: EmptyBitboard.Set(A3).Set(B4),
			expectedDoubles: EmptyBitboard.Set(A4),
		},
		{
			name:            "white pawn on rank 7",
			occupancy:       EmptyBitboard,
			pawns:           EmptyBitboard.Set(E7),
			colour:          White,
			expectedSingles: EmptyBitboard.Set(E8),
			expectedDoubles: EmptyBitboard,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			singles, doubles := findPawnPushes(test.colour, test.pawns, test.occupancy|test.pawns)
			if singles != test.expectedSingles {
				t.Errorf("Expected single pawn pushes to be %s, got %s", test.expectedSingles, singles)
			}
			if doubles != test.expectedDoubles {
				t.Errorf("Expected double pawn pushes to be %s, got %s", test.expectedDoubles, doubles)
			}
		})
	}
}

func TestPawnPushesFuzzAgainstSlowImpl(t *testing.T) {
	slowImpl := func(colour Colour, pawns Bitboard, bothOccupancy Bitboard) (singles, doubles Bitboard) {
		var sq Square
		for {
			sq, pawns = pawns.PopSquare()
			if sq == NoSquare {
				break
			}
			switch colour {
			case White:
				if rank := sq.Rank(); rank < Rank8 {
					file := sq.File()
					target, err := SquareFromFileAndRank(file, rank+1)
					if err == nil {
						singlePush := target.Bitboard()
						if bothOccupancy&singlePush == 0 {
							singles |= singlePush
							if rank == Rank2 {
								doublePush := singlePush << 8
								if bothOccupancy&doublePush == 0 {
									doubles |= doublePush
								}
							}
						}
					}
				}
			case Black:
				if rank := sq.Rank(); rank > Rank1 {
					target, err := SquareFromFileAndRank(sq.File(), rank-1)
					if err == nil {
						singlePush := target.Bitboard()
						if bothOccupancy&singlePush == 0 {
							singles |= singlePush
							if sq.Rank() == Rank7 {
								doublePush := singlePush >> 8
								if bothOccupancy&doublePush == 0 {
									doubles |= doublePush
								}
							}
						}
					}
				}
			}
		}
		return singles, doubles
	}

	rnd := rand.New(rand.NewPCG(1, 2))

	var colour Colour
	for sparsity := range 3 {
		for range 1000 {
			colour = colour.Opposite()
			pawns := Bitboard(rnd.Uint64() & rnd.Uint64() & rnd.Uint64()) // always realistic pawn density
			occupancy := Bitboard(rnd.Uint64())
			for range sparsity {
				occupancy &= Bitboard(rnd.Uint64()) // reduce density
			}
			occupancy |= pawns

			actualSingles, actualDoubles := findPawnPushes(colour, pawns, occupancy)
			expectedSingles, expectedDoubles := slowImpl(colour, pawns, occupancy)

			if actualSingles != expectedSingles {
				t.Fatalf("Fuzz test failed for colour %v, pawns %s, occupancy %s: expected singles %s, got %s", colour, pawns, occupancy, expectedSingles, actualSingles)
			}
			if actualDoubles != expectedDoubles {
				t.Fatalf("Fuzz test failed for colour %v, pawns %s, occupancy %s: expected doubles %s, got %s", colour, pawns, occupancy, expectedDoubles, actualDoubles)
			}
		}
	}
}
