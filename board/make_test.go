package board

import "testing"

var makeMoveTests = []struct {
	name        string
	fen         string
	move        Move
	expectedFEN string
}{
	{
		name:        "double push with no adjacent pawn leaves ep unset",
		fen:         "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		move:        NewMove(E2, E4, DoublePawnPush),
		expectedFEN: "rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1",
	},
	{
		name:        "knight quiet move",
		fen:         "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		move:        NewMove(G1, F3, QuietMove),
		expectedFEN: "rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R b KQkq - 1 1",
	},
	{
		name:        "pawn quiet move resets clock",
		fen:         "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 5 3",
		move:        NewMove(E2, E3, QuietMove),
		expectedFEN: "rnbqkbnr/pppppppp/8/8/8/4P3/PPPP1PPP/RNBQKBNR b KQkq - 0 3",
	},
	{
		name:        "double push sets ep square",
		fen:         "rnbqkbnr/ppp1pppp/8/8/3p4/8/PPPPPPPP/RNBQKBNR w KQkq - 0 3",
		move:        NewMove(E2, E4, DoublePawnPush),
		expectedFEN: "rnbqkbnr/ppp1pppp/8/8/3pP3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 3",
	},
	{
		name:        "capture by knight resets clock",
		fen:         "rnbqkbnr/pppp1ppp/8/4p3/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 3 3",
		move:        NewMove(F3, E5, Capture),
		expectedFEN: "rnbqkbnr/pppp1ppp/8/4N3/8/8/PPPPPPPP/RNBQKB1R b KQkq - 0 3",
	},
	{
		name:        "capture by pawn",
		fen:         "rnbqkbnr/ppp1pppp/8/3p4/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 3",
		move:        NewMove(E4, D5, Capture),
		expectedFEN: "rnbqkbnr/ppp1pppp/8/3P4/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 3",
	},
	{
		name:        "en passant capture",
		fen:         "rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 3",
		move:        NewMove(E5, D6, EnPassantCapture),
		expectedFEN: "rnbqkbnr/ppp1pppp/3P4/8/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 3",
	},
	{
		name:        "queen promotion",
		fen:         "8/1P6/8/8/8/8/8/K6k w - - 0 1",
		move:        NewMove(B7, B8, QueenPromotion),
		expectedFEN: "1Q6/8/8/8/8/8/8/K6k b - - 0 1",
	},
	{
		name:        "knight underpromotion",
		fen:         "8/1P6/8/8/8/8/8/K6k w - - 0 1",
		move:        NewMove(B7, B8, KnightPromotion),
		expectedFEN: "1N6/8/8/8/8/8/8/K6k b - - 0 1",
	},
	{
		name:        "promotion capture takes rook and its castling right",
		fen:         "r3k2r/1P6/8/8/8/8/8/R3K2R w KQkq - 0 1",
		move:        NewMove(B7, A8, QueenPromotionCapture),
		expectedFEN: "Q3k2r/8/8/8/8/8/8/R3K2R b KQk - 0 1",
	},
	{
		name:        "white kingside castle",
		fen:         "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
		move:        NewMove(E1, G1, KingsideCastle),
		expectedFEN: "r3k2r/8/8/8/8/8/8/R4RK1 b kq - 1 1",
	},
	{
		name:        "white queenside castle",
		fen:         "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
		move:        NewMove(E1, C1, QueensideCastle),
		expectedFEN: "r3k2r/8/8/8/8/8/8/2KR3R b kq - 1 1",
	},
	{
		name:        "black kingside castle",
		fen:         "r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
		move:        NewMove(E8, G8, KingsideCastle),
		expectedFEN: "r4rk1/8/8/8/8/8/8/R3K2R w KQ - 1 2",
	},
	{
		name:        "black queenside castle",
		fen:         "r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
		move:        NewMove(E8, C8, QueensideCastle),
		expectedFEN: "2kr3r/8/8/8/8/8/8/R3K2R w KQ - 1 2",
	},
	{
		name:        "rook move loses one castling right",
		fen:         "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
		move:        NewMove(H1, H5, QuietMove),
		expectedFEN: "r3k2r/8/8/7R/8/8/8/R3K3 b Qkq - 1 1",
	},
	{
		name:        "king move loses both castling rights",
		fen:         "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
		move:        NewMove(E1, E2, QuietMove),
		expectedFEN: "r3k2r/8/8/8/8/8/4K3/R6R b kq - 1 1",
	},
	{
		name:        "black en passant capture",
		fen:         "rnbqkbnr/ppp1pppp/8/8/3pP3/8/PPPP1PPP/RNBQKBNR b KQkq e3 0 3",
		move:        NewMove(D4, E3, EnPassantCapture),
		expectedFEN: "rnbqkbnr/ppp1pppp/8/8/8/4p3/PPPP1PPP/RNBQKBNR w KQkq - 0 4",
	},
	{
		name:        "black double push sets ep square",
		fen:         "rnbqkbnr/pppppppp/8/4P3/8/8/PPPP1PPP/RNBQKBNR b KQkq - 0 3",
		move:        NewMove(D7, D5, DoublePawnPush),
		expectedFEN: "rnbqkbnr/ppp1pppp/8/3pP3/8/8/PPPP1PPP/RNBQKBNR w KQkq d6 0 4",
	},
	{
		name:        "black queen promotion",
		fen:         "K6k/8/8/8/8/8/1p6/8 b - - 0 1",
		move:        NewMove(B2, B1, QueenPromotion),
		expectedFEN: "K6k/8/8/8/8/8/8/1q6 w - - 0 2",
	},
	{
		name:        "rook takes rook on its original square",
		fen:         "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
		move:        NewMove(A1, A8, Capture),
		expectedFEN: "R3k2r/8/8/8/8/8/8/4K2R b Kk - 0 1",
	},
	{
		name:        "black rook move loses one castling right",
		fen:         "r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
		move:        NewMove(H8, H5, QuietMove),
		expectedFEN: "r3k3/8/8/7r/8/8/8/R3K2R w KQq - 1 2",
	},
	{
		name:        "black king move loses both castling rights",
		fen:         "r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
		move:        NewMove(E8, E7, QuietMove),
		expectedFEN: "r6r/4k3/8/8/8/8/8/R3K2R w KQ - 1 2",
	},
	{
		name:        "rook promotion capture onto h8",
		fen:         "r3k2r/6P1/8/8/8/8/8/R3K2R w KQkq - 0 1",
		move:        NewMove(G7, H8, RookPromotionCapture),
		expectedFEN: "r3k2R/8/8/8/8/8/8/R3K2R b KQq - 0 1",
	},
	{
		name:        "black promotion capture takes rook and its castling right",
		fen:         "r3k2r/8/8/8/8/8/1p6/R3K2R b KQkq - 0 1",
		move:        NewMove(B2, A1, QueenPromotionCapture),
		expectedFEN: "r3k2r/8/8/8/8/8/8/q3K2R w Kkq - 0 2",
	},
	{
		name:        "knight promotion capture",
		fen:         "r3k2r/1P6/8/8/8/8/8/R3K2R w KQkq - 0 1",
		move:        NewMove(B7, A8, KnightPromotionCapture),
		expectedFEN: "N3k2r/8/8/8/8/8/8/R3K2R b KQk - 0 1",
	},
	{
		name:        "bishop promotion",
		fen:         "8/1P6/8/8/8/8/8/K6k w - - 0 1",
		move:        NewMove(B7, B8, BishopPromotion),
		expectedFEN: "1B6/8/8/8/8/8/8/K6k b - - 0 1",
	},
	{
		name:        "rook promotion",
		fen:         "8/1P6/8/8/8/8/8/K6k w - - 0 1",
		move:        NewMove(B7, B8, RookPromotion),
		expectedFEN: "1R6/8/8/8/8/8/8/K6k b - - 0 1",
	},
	{
		name:        "capture a queen",
		fen:         "rnb1kbnr/pppppppp/8/6q1/8/5N2/PPPPPPPP/RNBQKB1R w KQkq - 4 4",
		move:        NewMove(F3, G5, Capture),
		expectedFEN: "rnb1kbnr/pppppppp/8/6N1/8/8/PPPPPPPP/RNBQKB1R b KQkq - 0 4",
	},
}

func TestMake(t *testing.T) {
	for _, test := range makeMoveTests {
		t.Run(test.name, func(t *testing.T) {
			pos, err := ParseFEN(test.fen)
			if err != nil {
				t.Fatalf("Failed to parse FEN: %v", err)
			}

			pos.MakeMove(test.move)
			actualFEN := GenerateFEN(pos)
			if test.expectedFEN != actualFEN {
				t.Errorf("Expected FEN: %s, got: %s", test.expectedFEN, actualFEN)
			}
		})
	}
}

// TestMakeUnmakeRoundTrip asserts that unmaking a move restores the position to
// exactly its prior state - every bitboard, the mailbox, and all of the state
// fields. It reuses the cases from makeMoveTests so that every move kind is
// covered in both directions.
func TestMakeUnmakeRoundTrip(t *testing.T) {
	for _, test := range makeMoveTests {
		t.Run(test.name, func(t *testing.T) {
			pos, err := ParseFEN(test.fen)
			if err != nil {
				t.Fatalf("Failed to parse FEN: %v", err)
			}
			before := *pos

			undo := pos.MakeMove(test.move)
			if err := pos.validateStateSlow(); err != nil {
				t.Fatalf("Position corrupt after MakeMove: %v", err)
			}

			pos.UnmakeMove(undo)
			if err := pos.validateStateSlow(); err != nil {
				t.Fatalf("Position corrupt after UnmakeMove: %v", err)
			}

			if *pos != before {
				t.Errorf("Position not restored.\nExpected FEN: %s\nGot FEN:      %s",
					GenerateFEN(&before), GenerateFEN(pos))
			}
		})
	}
}

// TestMakeUnmakeSequence plays a series of moves and then unmakes them all in
// reverse order, asserting the original position is recovered. A single-move
// round trip cannot catch state that is restored to a plausible-but-wrong value,
// nor errors that only appear once several plies of history are outstanding.
func TestMakeUnmakeSequence(t *testing.T) {
	tests := []struct {
		name  string
		fen   string
		moves []Move
	}{
		{
			name: "italian game into kingside castle",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			moves: []Move{
				NewMove(E2, E4, DoublePawnPush),
				NewMove(E7, E5, DoublePawnPush),
				NewMove(G1, F3, QuietMove),
				NewMove(B8, C6, QuietMove),
				NewMove(F1, C4, QuietMove),
				NewMove(F8, C5, QuietMove),
				NewMove(E1, G1, KingsideCastle),
			},
		},
		{
			name: "pawn advance into en passant capture",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			moves: []Move{
				NewMove(E2, E4, DoublePawnPush),
				NewMove(A7, A6, QuietMove),
				NewMove(E4, E5, QuietMove),
				NewMove(D7, D5, DoublePawnPush),
				NewMove(E5, D6, EnPassantCapture),
			},
		},
		{
			name: "both sides castle queenside",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			moves: []Move{
				NewMove(D2, D4, DoublePawnPush),
				NewMove(D7, D5, DoublePawnPush),
				NewMove(C1, F4, QuietMove),
				NewMove(C8, F5, QuietMove),
				NewMove(B1, C3, QuietMove),
				NewMove(B8, C6, QuietMove),
				NewMove(D1, D3, QuietMove),
				NewMove(D8, D7, QuietMove),
				NewMove(E1, C1, QueensideCastle),
				NewMove(E8, C8, QueensideCastle),
			},
		},
		{
			name: "promotion capture with further moves on top",
			fen:  "r3k2r/1P6/8/8/8/8/8/R3K2R w KQkq - 0 1",
			moves: []Move{
				NewMove(B7, A8, QueenPromotionCapture),
				NewMove(E8, G8, KingsideCastle),
				NewMove(A8, A5, QuietMove),
				NewMove(F8, E8, QuietMove),
				NewMove(E1, G1, KingsideCastle),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pos, err := ParseFEN(test.fen)
			if err != nil {
				t.Fatalf("Failed to parse FEN: %v", err)
			}
			before := *pos

			undos := make([]Undo, 0, len(test.moves))
			for i, move := range test.moves {
				undos = append(undos, pos.MakeMove(move))
				if err := pos.validateStateSlow(); err != nil {
					t.Fatalf("Position corrupt after move %d (%s): %v", i+1, move, err)
				}
			}

			for i := len(undos) - 1; i >= 0; i-- {
				pos.UnmakeMove(undos[i])
				if err := pos.validateStateSlow(); err != nil {
					t.Fatalf("Position corrupt after unmaking move %d (%s): %v", i+1, test.moves[i], err)
				}
			}

			if *pos != before {
				t.Errorf("Position not restored after %d moves.\nExpected FEN: %s\nGot FEN:      %s",
					len(test.moves), GenerateFEN(&before), GenerateFEN(pos))
			}
		})
	}
}
