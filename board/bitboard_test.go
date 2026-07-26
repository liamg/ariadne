package board

import (
	"slices"
	"testing"
)

func TestBitboardHasSetClear(t *testing.T) {
	tests := []struct {
		bitboard Bitboard
		square   Square
		expected bool
	}{
		{EmptyBitboard, A1, false},
		{EmptyBitboard, H8, false},
		{FullBitboard, A1, true},
		{FullBitboard, H8, true},
		{EmptyBitboard.Set(A1), A1, true},
		{EmptyBitboard.Set(A3), A3, true},
		{EmptyBitboard.Set(H8), H8, true},
		{FullBitboard.Clear(A1), A1, false},
		{FullBitboard.Clear(A3), A3, false},
		{FullBitboard.Clear(H8), H8, false},
	}

	for _, test := range tests {
		if test.bitboard.Has(test.square) != test.expected {
			t.Errorf("Expected bitboard %v to have square %v: %v, got %v", test.bitboard, test.square, test.expected, test.bitboard.Has(test.square))
		}

		bitboardWithSquare := test.bitboard.Set(test.square)
		if !bitboardWithSquare.Has(test.square) {
			t.Errorf("Expected bitboard %v to have square %v after setting it, but it does not", bitboardWithSquare, test.square)
		}

		bitboardCleared := bitboardWithSquare.Clear(test.square)
		if bitboardCleared.Has(test.square) {
			t.Errorf("Expected bitboard %v to not have square %v after clearing it, but it does", bitboardCleared, test.square)
		}
	}
}

func TestBitboardCount(t *testing.T) {
	tests := []struct {
		bitboard Bitboard
		expected int
	}{
		{EmptyBitboard, 0},
		{FullBitboard, 64},
		{EmptyBitboard.Set(A1), 1},
		{EmptyBitboard.Set(A1).Set(B2), 2},
		{EmptyBitboard.Set(A1).Set(B2).Set(C3), 3},
		{FullBitboard.Clear(A1), 63},
		{FullBitboard.Clear(A1).Clear(B2), 62},
		{FullBitboard.Clear(A1).Clear(B2).Clear(C3), 61},
	}

	for _, test := range tests {
		if count := test.bitboard.Count(); count != test.expected {
			t.Errorf("Expected bitboard %v to have %d bits set, got %d", test.bitboard, test.expected, count)
		}
	}
}

func TestBitboardString(t *testing.T) {
	tests := []struct {
		bitboard Bitboard
		expected string
	}{
		{
			EmptyBitboard,
			`
8 . . . . . . . . 
7 . . . . . . . . 
6 . . . . . . . . 
5 . . . . . . . . 
4 . . . . . . . . 
3 . . . . . . . . 
2 . . . . . . . . 
1 . . . . . . . . 
  a b c d e f g h 
`,
		},
		{
			FullBitboard,
			`
8 x x x x x x x x 
7 x x x x x x x x 
6 x x x x x x x x 
5 x x x x x x x x 
4 x x x x x x x x 
3 x x x x x x x x 
2 x x x x x x x x 
1 x x x x x x x x 
  a b c d e f g h 
`,
		},
		{
			FullBitboard.Clear(A5).Clear(B2),
			`
8 x x x x x x x x 
7 x x x x x x x x 
6 x x x x x x x x 
5 . x x x x x x x 
4 x x x x x x x x 
3 x x x x x x x x 
2 x . x x x x x x 
1 x x x x x x x x 
  a b c d e f g h 
`,
		},
	}

	for _, test := range tests {
		if str := test.bitboard.String(); str != test.expected {
			t.Errorf("Expected bitboard %v to have string representation:\n%s\nGot:\n%s", test.bitboard, test.expected, str)
		}
	}
}

func TestBitboardPopSquare(t *testing.T) {
	tests := []struct {
		bitboard Bitboard
		expected []Square
	}{
		{EmptyBitboard, []Square{}},
		{
			FullBitboard, []Square{
				A1, B1, C1, D1, E1, F1, G1, H1,
				A2, B2, C2, D2, E2, F2, G2, H2,
				A3, B3, C3, D3, E3, F3, G3, H3,
				A4, B4, C4, D4, E4, F4, G4, H4,
				A5, B5, C5, D5, E5, F5, G5, H5,
				A6, B6, C6, D6, E6, F6, G6, H6,
				A7, B7, C7, D7, E7, F7, G7, H7,
				A8, B8, C8, D8, E8, F8, G8, H8,
			},
		},
		{
			EmptyBitboard.Set(A1).Set(B2).Set(C3).Set(D4).Set(E5).Set(F6).Set(G7).Set(H8),
			[]Square{A1, B2, C3, D4, E5, F6, G7, H8},
		},
		{
			EmptyBitboard.Set(A2),
			[]Square{A2},
		},
	}

	for _, test := range tests {
		bitboard := test.bitboard
		var poppedSquares []Square
		for {
			sq, newBitboard := bitboard.PopSquare()
			if sq == NoSquare {
				break
			}
			poppedSquares = append(poppedSquares, sq)
			bitboard = newBitboard
		}

		if !slices.Equal(poppedSquares, test.expected) {
			t.Errorf("Expected bitboard %v to pop squares %v, got %v", test.bitboard, test.expected, poppedSquares)
		}

	}
}

func TestBitboardNorth(t *testing.T) {
	tests := []struct {
		bitboard Bitboard
		expected Bitboard
	}{
		{EmptyBitboard, EmptyBitboard},
		{FullBitboard, 0xFFFFFFFFFFFFFF00},
		{EmptyBitboard.Set(A1), EmptyBitboard.Set(A2)},
		{EmptyBitboard.Set(A1).Set(B4), EmptyBitboard.Set(A2).Set(B5)},
		{EmptyBitboard.Set(H8), EmptyBitboard},
	}

	for _, test := range tests {
		if north := test.bitboard.North(); north != test.expected {
			t.Errorf("Expected bitboard %v to have north bitboard %v, got %v", test.bitboard, test.expected, north)
		}
	}
}

func TestBitboardSouth(t *testing.T) {
	tests := []struct {
		bitboard Bitboard
		expected Bitboard
	}{
		{EmptyBitboard, EmptyBitboard},
		{FullBitboard, 0x00FFFFFFFFFFFFFF},
		{EmptyBitboard.Set(A2), EmptyBitboard.Set(A1)},
		{EmptyBitboard.Set(A2).Set(B5), EmptyBitboard.Set(A1).Set(B4)},
		{EmptyBitboard.Set(H1), EmptyBitboard},
	}

	for _, test := range tests {
		if south := test.bitboard.South(); south != test.expected {
			t.Errorf("Expected bitboard %v to have south bitboard %v, got %v", test.bitboard, test.expected, south)
		}
	}
}

func TestBitboardEast(t *testing.T) {
	tests := []struct {
		bitboard Bitboard
		expected Bitboard
	}{
		{EmptyBitboard, EmptyBitboard},
		{FullBitboard, 0xFEFEFEFEFEFEFEFE},
		{EmptyBitboard.Set(A1), EmptyBitboard.Set(B1)},
		{EmptyBitboard.Set(A1).Set(G4), EmptyBitboard.Set(B1).Set(H4)},
		{EmptyBitboard.Set(D8), EmptyBitboard.Set(E8)},
		{EmptyBitboard.Set(H2), EmptyBitboard},
	}

	for _, test := range tests {
		if east := test.bitboard.East(); east != test.expected {
			t.Errorf("Expected bitboard %v to have east bitboard %v, got %v", test.bitboard, test.expected, east)
		}
	}
}

func TestBitboardWest(t *testing.T) {
	tests := []struct {
		bitboard Bitboard
		expected Bitboard
	}{
		{EmptyBitboard, EmptyBitboard},
		{FullBitboard, 0x7F7F7F7F7F7F7F7F},
		{EmptyBitboard.Set(B1), EmptyBitboard.Set(A1)},
		{EmptyBitboard.Set(B1).Set(H4), EmptyBitboard.Set(A1).Set(G4)},
		{EmptyBitboard.Set(A4), EmptyBitboard},
	}

	for _, test := range tests {
		if west := test.bitboard.West(); west != test.expected {
			t.Errorf("Expected bitboard %v to have west bitboard %v, got %v", test.bitboard, test.expected, west)
		}
	}
}

func TestBitboardNorthEast(t *testing.T) {
	tests := []struct {
		bitboard Bitboard
		expected Bitboard
	}{
		{EmptyBitboard, EmptyBitboard},
		{FullBitboard, FullBitboard &^ FileAMask &^ Rank1Mask},
		{EmptyBitboard.Set(A1), EmptyBitboard.Set(B2)},
		{EmptyBitboard.Set(A1).Set(G4), EmptyBitboard.Set(B2).Set(H5)},
		{EmptyBitboard.Set(D8).Set(H2), EmptyBitboard},
	}

	for _, test := range tests {
		if ne := test.bitboard.NorthEast(); ne != test.expected {
			t.Errorf("Expected bitboard %v to have north-east bitboard %v, got %v", test.bitboard, test.expected, ne)
		}
	}
}

func TestBitboardNorthWest(t *testing.T) {
	tests := []struct {
		bitboard Bitboard
		expected Bitboard
	}{
		{EmptyBitboard, EmptyBitboard},
		{FullBitboard, FullBitboard &^ FileHMask &^ Rank1Mask},
		{EmptyBitboard.Set(B1), EmptyBitboard.Set(A2)},
		{EmptyBitboard.Set(B1).Set(H4), EmptyBitboard.Set(A2).Set(G5)},
		{EmptyBitboard.Set(A1).Set(H8), EmptyBitboard},
	}

	for _, test := range tests {
		if nw := test.bitboard.NorthWest(); nw != test.expected {
			t.Errorf("Expected bitboard %v to have north-west bitboard %v, got %v", test.bitboard, test.expected, nw)
		}
	}
}

func TestBitboardSouthEast(t *testing.T) {
	tests := []struct {
		bitboard Bitboard
		expected Bitboard
	}{
		{EmptyBitboard, EmptyBitboard},
		{FullBitboard, FullBitboard &^ FileAMask &^ Rank8Mask},
		{EmptyBitboard.Set(A2), EmptyBitboard.Set(B1)},
		{EmptyBitboard.Set(A2).Set(G5), EmptyBitboard.Set(B1).Set(H4)},
		{EmptyBitboard.Set(D1).Set(H4), EmptyBitboard},
	}

	for _, test := range tests {
		if se := test.bitboard.SouthEast(); se != test.expected {
			t.Errorf("Expected bitboard %v to have south-east bitboard %v, got %v", test.bitboard, test.expected, se)
		}
	}
}

func TestBitboardSouthWest(t *testing.T) {
	tests := []struct {
		bitboard Bitboard
		expected Bitboard
	}{
		{EmptyBitboard, EmptyBitboard},
		{FullBitboard, FullBitboard &^ FileHMask &^ Rank8Mask},
		{EmptyBitboard.Set(B2), EmptyBitboard.Set(A1)},
		{EmptyBitboard.Set(B2).Set(H5), EmptyBitboard.Set(A1).Set(G4)},
		{EmptyBitboard.Set(A8).Set(H1), EmptyBitboard},
	}

	for _, test := range tests {
		if sw := test.bitboard.SouthWest(); sw != test.expected {
			t.Errorf("Expected bitboard %v to have south-west bitboard %v, got %v", test.bitboard, test.expected, sw)
		}
	}
}
