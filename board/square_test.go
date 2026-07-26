package board

import "testing"

func TestRankString(t *testing.T) {
	tests := []struct {
		rank     Rank
		expected string
	}{
		{Rank1, "1"},
		{Rank2, "2"},
		{Rank3, "3"},
		{Rank4, "4"},
		{Rank5, "5"},
		{Rank6, "6"},
		{Rank7, "7"},
		{Rank8, "8"},
		{Rank(8), "?"},
	}

	for _, test := range tests {
		if test.rank.String() != test.expected {
			t.Errorf("Expected rank %v to be %s, got %s", test.rank, test.expected, test.rank.String())
		}
	}
}

func TestFileString(t *testing.T) {
	tests := []struct {
		file     File
		expected string
	}{
		{FileA, "a"},
		{FileB, "b"},
		{FileC, "c"},
		{FileD, "d"},
		{FileE, "e"},
		{FileF, "f"},
		{FileG, "g"},
		{FileH, "h"},
		{File(8), "?"},
	}

	for _, test := range tests {
		if test.file.String() != test.expected {
			t.Errorf("Expected file %v to be %s, got %s", test.file, test.expected, test.file.String())
		}
	}
}

func TestSquareString(t *testing.T) {
	tests := []struct {
		square   Square
		expected string
	}{
		{A1, "a1"},
		{A2, "a2"},
		{A3, "a3"},
		{A4, "a4"},
		{A5, "a5"},
		{A6, "a6"},
		{A7, "a7"},
		{A8, "a8"},
		{B1, "b1"},
		{B2, "b2"},
		{B3, "b3"},
		{B4, "b4"},
		{B5, "b5"},
		{B6, "b6"},
		{B7, "b7"},
		{B8, "b8"},
		{C1, "c1"},
		{C2, "c2"},
		{C3, "c3"},
		{C4, "c4"},
		{C5, "c5"},
		{C6, "c6"},
		{C7, "c7"},
		{C8, "c8"},
		{D1, "d1"},
		{D2, "d2"},
		{D3, "d3"},
		{D4, "d4"},
		{D5, "d5"},
		{D6, "d6"},
		{D7, "d7"},
		{D8, "d8"},
		{E1, "e1"},
		{E2, "e2"},
		{E3, "e3"},
		{E4, "e4"},
		{E5, "e5"},
		{E6, "e6"},
		{E7, "e7"},
		{E8, "e8"},
		{F1, "f1"},
		{F2, "f2"},
		{F3, "f3"},
		{F4, "f4"},
		{F5, "f5"},
		{F6, "f6"},
		{F7, "f7"},
		{F8, "f8"},
		{G1, "g1"},
		{G2, "g2"},
		{G3, "g3"},
		{G4, "g4"},
		{G5, "g5"},
		{G6, "g6"},
		{G7, "g7"},
		{G8, "g8"},
		{H1, "h1"},
		{H2, "h2"},
		{H3, "h3"},
		{H4, "h4"},
		{H5, "h5"},
		{H6, "h6"},
		{H7, "h7"},
		{H8, "h8"},
		{NoSquare, "-"},
		{Square(-1), "?"},
		{Square(65), "?"},
	}

	for _, test := range tests {
		if test.square.String() != test.expected {
			t.Errorf("Expected square %v to be %s, got %s", test.square, test.expected, test.square.String())
		}
	}
}

func TestSquareFileAndRank(t *testing.T) {
	tests := []struct {
		square       Square
		expectedFile File
		expectedRank Rank
	}{
		{A1, FileA, Rank1},
		{B1, FileB, Rank1},
		{C1, FileC, Rank1},
		{D1, FileD, Rank1},
		{E1, FileE, Rank1},
		{F1, FileF, Rank1},
		{G1, FileG, Rank1},
		{H1, FileH, Rank1},
		{A2, FileA, Rank2},
		{B2, FileB, Rank2},
		{C2, FileC, Rank2},
		{D2, FileD, Rank2},
		{E2, FileE, Rank2},
		{F2, FileF, Rank2},
		{G2, FileG, Rank2},
		{H2, FileH, Rank2},
		{A3, FileA, Rank3},
		{B3, FileB, Rank3},
		{C3, FileC, Rank3},
		{D3, FileD, Rank3},
		{E3, FileE, Rank3},
		{F3, FileF, Rank3},
		{G3, FileG, Rank3},
		{H3, FileH, Rank3},
		{A4, FileA, Rank4},
		{B4, FileB, Rank4},
		{C4, FileC, Rank4},
		{D4, FileD, Rank4},
		{E4, FileE, Rank4},
		{F4, FileF, Rank4},
		{G4, FileG, Rank4},
		{H4, FileH, Rank4},
		{A5, FileA, Rank5},
		{B5, FileB, Rank5},
		{C5, FileC, Rank5},
		{D5, FileD, Rank5},
		{E5, FileE, Rank5},
		{F5, FileF, Rank5},
		{G5, FileG, Rank5},
		{H5, FileH, Rank5},
		{A6, FileA, Rank6},
		{B6, FileB, Rank6},
		{C6, FileC, Rank6},
		{D6, FileD, Rank6},
		{E6, FileE, Rank6},
		{F6, FileF, Rank6},
		{G6, FileG, Rank6},
		{H6, FileH, Rank6},
		{A7, FileA, Rank7},
		{B7, FileB, Rank7},
		{C7, FileC, Rank7},
		{D7, FileD, Rank7},
		{E7, FileE, Rank7},
		{F7, FileF, Rank7},
		{G7, FileG, Rank7},
		{H7, FileH, Rank7},
		{A8, FileA, Rank8},
		{B8, FileB, Rank8},
		{C8, FileC, Rank8},
		{D8, FileD, Rank8},
		{E8, FileE, Rank8},
		{F8, FileF, Rank8},
		{G8, FileG, Rank8},
		{H8, FileH, Rank8},
	}

	for _, test := range tests {
		if test.square.File() != test.expectedFile {
			t.Errorf("Expected square %v to have file %v, got %v", test.square, test.expectedFile, test.square.File())
		}
		if test.square.Rank() != test.expectedRank {
			t.Errorf("Expected square %v to have rank %v, got %v", test.square, test.expectedRank, test.square.Rank())
		}
	}
}

func TestSquareParse(t *testing.T) {
	tests := []struct {
		input    string
		expected Square
		wantErr  bool
	}{
		{"a1", A1, false},
		{"b1", B1, false},
		{"c1", C1, false},
		{"d1", D1, false},
		{"e1", E1, false},
		{"f1", F1, false},
		{"g1", G1, false},
		{"h1", H1, false},
		{"a2", A2, false},
		{"b2", B2, false},
		{"c2", C2, false},
		{"d2", D2, false},
		{"e2", E2, false},
		{"f2", F2, false},
		{"g2", G2, false},
		{"h2", H2, false},
		{"a3", A3, false},
		{"b3", B3, false},
		{"c3", C3, false},
		{"d3", D3, false},
		{"e3", E3, false},
		{"f3", F3, false},
		{"g3", G3, false},
		{"h3", H3, false},
		{"a4", A4, false},
		{"b4", B4, false},
		{"c4", C4, false},
		{"d4", D4, false},
		{"e4", E4, false},
		{"f4", F4, false},
		{"g4", G4, false},
		{"h4", H4, false},
		{"a5", A5, false},
		{"b5", B5, false},
		{"c5", C5, false},
		{"d5", D5, false},
		{"e5", E5, false},
		{"f5", F5, false},
		{"g5", G5, false},
		{"h5", H5, false},
		{"a6", A6, false},
		{"b6", B6, false},
		{"c6", C6, false},
		{"d6", D6, false},
		{"e6", E6, false},
		{"f6", F6, false},
		{"g6", G6, false},
		{"h6", H6, false},
		{"a7", A7, false},
		{"b7", B7, false},
		{"c7", C7, false},
		{"d7", D7, false},
		{"e7", E7, false},
		{"f7", F7, false},
		{"g7", G7, false},
		{"h7", H7, false},
		{"a8", A8, false},
		{"b8", B8, false},
		{"c8", C8, false},
		{"d8", D8, false},
		{"e8", E8, false},
		{"f8", F8, false},
		{"g8", G8, false},
		{"h8", H8, false},
		{"x2", NoSquare, true},
		{"a9", NoSquare, true},
		{"", NoSquare, true},
		{"long", NoSquare, true},
		{"a", NoSquare, true},
		{"a12", NoSquare, true},
	}

	for _, test := range tests {
		square, err := ParseSquare(test.input)
		if err != nil {
			if !test.wantErr {
				t.Errorf("Unexpected error parsing square %s: %v", test.input, err)
			}
		} else if test.wantErr {
			t.Errorf("Expected error parsing square %s, but got none", test.input)
		}
		if square != test.expected {
			t.Errorf("Expected square %s to parse to %v, got %v", test.input, test.expected, square)
		}
	}
}

func TestSquareToBitboard(t *testing.T) {
	tests := []struct {
		square   Square
		expected Bitboard
	}{
		{A1, 1 << A1},
		{B1, 1 << B1},
		{C1, 1 << C1},
		{D1, 1 << D1},
		{E1, 1 << E1},
		{F1, 1 << F1},
		{G1, 1 << G1},
		{H1, 1 << H1},
		{A2, 1 << A2},
		{B2, 1 << B2},
		{C2, 1 << C2},
		{D2, 1 << D2},
		{E2, 1 << E2},
		{F2, 1 << F2},
		{G2, 1 << G2},
		{H2, 1 << H2},
		{A3, 1 << A3},
		{B3, 1 << B3},
		{C3, 1 << C3},
		{D3, 1 << D3},
		{E3, 1 << E3},
		{F3, 1 << F3},
		{G3, 1 << G3},
		{H3, 1 << H3},
		{A4, 1 << A4},
		{B4, 1 << B4},
		{C4, 1 << C4},
		{D4, 1 << D4},
		{E4, 1 << E4},
		{F4, 1 << F4},
		{G4, 1 << G4},
		{H4, 1 << H4},
		{A5, 1 << A5},
		{B5, 1 << B5},
		{C5, 1 << C5},
		{D5, 1 << D5},
		{E5, 1 << E5},
		{F5, 1 << F5},
		{G5, 1 << G5},
		{H5, 1 << H5},
		{A6, 1 << A6},
		{B6, 1 << B6},
		{C6, 1 << C6},
		{D6, 1 << D6},
		{E6, 1 << E6},
		{F6, 1 << F6},
		{G6, 1 << G6},
		{H6, 1 << H6},
		{A7, 1 << A7},
		{B7, 1 << B7},
		{C7, 1 << C7},
		{D7, 1 << D7},
		{E7, 1 << E7},
		{F7, 1 << F7},
		{G7, 1 << G7},
		{H7, 1 << H7},
		{A8, 1 << A8},
		{B8, 1 << B8},
		{C8, 1 << C8},
		{D8, 1 << D8},
		{E8, 1 << E8},
		{F8, 1 << F8},
		{G8, 1 << G8},
		{H8, 1 << H8},
		{NoSquare, 0},
	}

	for _, test := range tests {
		if test.square.Bitboard() != test.expected {
			t.Errorf("Expected square %v to convert to bitboard %v, got %v", test.square, test.expected, test.square.Bitboard())
		}
	}
}
