package board

import (
	"strings"
	"testing"
)

func TestPositionString(t *testing.T) {
	tests := []struct {
		pos      *Position
		expected string
	}{
		{
			StartingPosition(),
			`
8 r n b q k b n r 
7 p p p p p p p p 
6 . . . . . . . . 
5 . . . . . . . . 
4 . . . . . . . . 
3 . . . . . . . . 
2 P P P P P P P P 
1 R N B Q K B N R 
  a b c d e f g h 
`,
		},
	}

	for _, test := range tests {
		if str := test.pos.String(); str != test.expected {
			t.Errorf("Expected position to have string representation:\n%s\nGot:\n%s", test.expected, str)
		}
		if err := test.pos.validateStateSlow(); err != nil {
			t.Errorf("Position validation failed: %v", err)
		}
	}
}

func TestStateCorruptionIsDetected(t *testing.T) {
	tests := []struct {
		name        string
		corruptFunc func() *Position
		expectedErr string
	}{
		{
			name: "piece missing from byType",
			corruptFunc: func() *Position {
				pos := StartingPosition()
				pos.byType[King] = 0 // Corrupt the state by removing the kings from the piece type bitboard
				return pos
			},
			expectedErr: "mailbox has piece not present in byType",
		},
		{
			name: "piece missing from byColour",
			corruptFunc: func() *Position {
				pos := StartingPosition()
				// remove piece from A1
				pos.byColour[White] &= ^A1.Bitboard()
				return pos
			},
			expectedErr: "mailbox has piece not present in byColour",
		},
		{
			name: "multiple pieces in byType for a square",
			corruptFunc: func() *Position {
				pos := StartingPosition()
				// Add a pawn to A1, which already has a rook
				pos.byType[Pawn] |= A1.Bitboard()
				return pos
			},
			expectedErr: "multiple pieces in byType for square",
		},
		{
			name: "byType has value not in mailbox",
			corruptFunc: func() *Position {
				pos := EmptyPosition()
				// Add a pawn to byType but not to mailbox
				pos.byType[Pawn] |= A1.Bitboard()
				return pos
			},
			expectedErr: "byType has piece not in mailbox, square: a1",
		},
		{
			name: "byColour has value not in mailbox",
			corruptFunc: func() *Position {
				pos := EmptyPosition()
				// Add a white piece to byColour but not to mailbox
				pos.byColour[White] |= A1.Bitboard()
				return pos
			},
			expectedErr: "byColour has piece not in mailbox, square: a1",
		},
		{
			name: "both colours on same square",
			corruptFunc: func() *Position {
				pos := EmptyPosition()
				pos.mailbox[A4] = WhitePawn
				pos.byType[Pawn] = A4.Bitboard()
				pos.byColour[White] |= A4.Bitboard()
				pos.byColour[Black] |= A4.Bitboard()
				return pos
			},
			expectedErr: "square has pieces of both colours, square: a4",
		},
		{
			name: "byType[0] non-empty",
			corruptFunc: func() *Position {
				pos := EmptyPosition()
				pos.byType[0] |= A4.Bitboard() // Corrupt the state by adding a square to the NoPieceType bitboard
				return pos
			},
			expectedErr: "byType has NoPieceType set",
		},
		{
			name: "invalid castling rights",
			corruptFunc: func() *Position {
				pos := EmptyPosition()
				pos.state.castlingRights = 0x10 // Corrupt the state by setting an invalid castling rights value
				return pos
			},
			expectedErr: "invalid castling rights",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pos := test.corruptFunc()
			if err := pos.validateStateSlow(); err == nil {
				t.Errorf("Expected position validation to fail due to corrupted state, but it passed")
			} else if !strings.Contains(err.Error(), test.expectedErr) && test.expectedErr != "" {
				t.Errorf("Expected error to contain '%s', but got: %v", test.expectedErr, err)
			}
		})
	}
}

func TestInCheck(t *testing.T) {
	// Ray generation is covered by TestIsSquareAttacked. These cases exist to pin
	// down what InCheck adds on top of it: that it looks up the king of the colour
	// it was asked about, and asks whether the *opposing* colour attacks it. Both
	// colours are asserted on every position, so a hardcoded colour cannot pass.
	tests := []struct {
		name  string
		fen   string
		white bool
		black bool
	}{
		{"neither king in check", "7k/8/8/8/8/8/8/K7 w - - 0 1", false, false},

		// Only one side is in check, so using the wrong king flips the answer.
		{"white in check from a rook", "4r2k/8/8/8/8/8/8/4K3 w - - 0 1", true, false},
		{"white in check from a knight", "7k/8/8/8/8/5n2/8/4K3 w - - 0 1", true, false},
		{"white in check from a queen on a diagonal", "7k/8/8/q7/8/8/8/4K3 w - - 0 1", true, false},
		{"black in check from a pawn", "8/8/8/4k3/3P4/8/8/4K3 b - - 0 1", false, true},

		// The rook's ray from e1 stops on white's own knight, so it never reaches e8.
		{"check blocked by an own piece", "4r2k/8/8/8/8/8/4N3/4K3 w - - 0 1", false, false},

		// White's own rook bears on e1 along the first rank. Asking about the wrong
		// colour - c rather than c.Opposite() - reports this as a check.
		{"own rook does not check its king", "7k/8/8/8/8/8/8/R3K3 w - - 0 1", false, false},

		// Illegal in a real game, but move generation depends on this answer to stop
		// the kings from moving next to each other.
		{"adjacent kings check each other", "8/8/8/8/8/8/4k3/4K3 w - - 0 1", true, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pos, err := ParseFEN(test.fen)
			if err != nil {
				t.Fatalf("Failed to parse FEN: %v", err)
			}

			if got := pos.InCheck(White); got != test.white {
				t.Errorf("InCheck(White) = %v, want %v", got, test.white)
			}
			if got := pos.InCheck(Black); got != test.black {
				t.Errorf("InCheck(Black) = %v, want %v", got, test.black)
			}
		})
	}
}

func TestIsDrawByRepetition(t *testing.T) {
	tests := []struct {
		name     string
		fen      string
		moves    []Move
		expected bool
	}{
		{
			name:     "no repetition",
			fen:      "7k/8/8/8/8/8/8/K7 w - - 0 1",
			moves:    []Move{},
			expected: false,
		},
		{
			// both kings walk out and back, so the start position recurs
			name: "kings shuffle back to the start",
			fen:  "7k/8/8/8/8/8/8/K7 w - - 0 1",
			moves: []Move{
				NewMove(A1, B1, QuietMove),
				NewMove(H8, G8, QuietMove),
				NewMove(B1, A1, QuietMove),
				NewMove(G8, H8, QuietMove),
			},
			expected: true,
		},
		{
			// one ply short - white is home but black is not, and it is black to move
			name: "three plies is not yet a repetition",
			fen:  "7k/8/8/8/8/8/8/K7 w - - 0 1",
			moves: []Move{
				NewMove(A1, B1, QuietMove),
				NewMove(H8, G8, QuietMove),
				NewMove(B1, A1, QuietMove),
			},
			expected: false,
		},
		{
			// the pawn push resets the halfmove clock, then the kings shuffle back
			// to the position that followed it. the clock bound must not cut the
			// scan short before reaching it
			name: "repetition of the position after an irreversible move",
			fen:  "7k/8/8/8/8/8/P7/K7 w - - 0 1",
			moves: []Move{
				NewMove(A2, A3, QuietMove),
				NewMove(H8, G8, QuietMove),
				NewMove(A1, B1, QuietMove),
				NewMove(G8, H8, QuietMove),
				NewMove(B1, A1, QuietMove),
			},
			expected: true,
		},
		{
			// same piece placement, but both sides lost castling rights on the way,
			// so the zobrist keys differ and it is not a repetition
			name: "same pieces but castling rights lost",
			fen:  "4k2r/8/8/8/8/8/8/4K2R w Kk - 0 1",
			moves: []Move{
				NewMove(H1, G1, QuietMove),
				NewMove(H8, G8, QuietMove),
				NewMove(G1, H1, QuietMove),
				NewMove(G8, H8, QuietMove),
			},
			expected: false,
		},
		{
			// halfmove clock from the fen exceeds the history length - the scan
			// bound must clamp at zero rather than indexing off the front
			name: "high halfmove clock from fen",
			fen:  "8/8/8/3k4/8/3K4/8/8 w - - 40 60",
			moves: []Move{
				NewMove(D3, C3, QuietMove),
				NewMove(D5, C5, QuietMove),
			},
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pos, err := ParseFEN(test.fen)
			if err != nil {
				t.Fatalf("Failed to parse FEN: %v", err)
			}

			for _, move := range test.moves {
				pos.MakeMove(move)
			}

			if got := pos.IsDrawByRepetition(); got != test.expected {
				t.Errorf("IsDrawByRepetition() = %v, want %v", got, test.expected)
			}
		})
	}
}
