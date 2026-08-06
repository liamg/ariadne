package board

import "testing"

func TestGeneratePseudoLegalMoves(t *testing.T) {
	tests := []struct {
		name          string
		fen           string
		expectedMoves []Move
	}{
		{
			name: "initial position",
			fen:  "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			expectedMoves: []Move{
				// pawn single pushes
				NewMove(A2, A3, QuietMove),
				NewMove(B2, B3, QuietMove),
				NewMove(C2, C3, QuietMove),
				NewMove(D2, D3, QuietMove),
				NewMove(E2, E3, QuietMove),
				NewMove(F2, F3, QuietMove),
				NewMove(G2, G3, QuietMove),
				NewMove(H2, H3, QuietMove),
				// pawn double pushes
				NewMove(A2, A4, DoublePawnPush),
				NewMove(B2, B4, DoublePawnPush),
				NewMove(C2, C4, DoublePawnPush),
				NewMove(D2, D4, DoublePawnPush),
				NewMove(E2, E4, DoublePawnPush),
				NewMove(F2, F4, DoublePawnPush),
				NewMove(G2, G4, DoublePawnPush),
				NewMove(H2, H4, DoublePawnPush),
				// knights
				NewMove(B1, A3, QuietMove),
				NewMove(B1, C3, QuietMove),
				NewMove(G1, F3, QuietMove),
				NewMove(G1, H3, QuietMove),
			},
		},
		{
			name: "lone knight in the centre",
			fen:  "7k/8/8/8/3N4/8/8/K7 w - - 0 1",
			expectedMoves: []Move{
				// knight on d4
				NewMove(D4, B3, QuietMove),
				NewMove(D4, B5, QuietMove),
				NewMove(D4, C2, QuietMove),
				NewMove(D4, C6, QuietMove),
				NewMove(D4, E2, QuietMove),
				NewMove(D4, E6, QuietMove),
				NewMove(D4, F3, QuietMove),
				NewMove(D4, F5, QuietMove),
				// king on a1
				NewMove(A1, A2, QuietMove),
				NewMove(A1, B1, QuietMove),
				NewMove(A1, B2, QuietMove),
			},
		},
		{
			name: "rook rays blocked by an enemy pawn",
			fen:  "7k/8/8/3p4/3R4/8/8/K7 w - - 0 1",
			expectedMoves: []Move{
				// rook on d4, north ray ends on the capture
				NewMove(D4, D5, Capture),
				NewMove(D4, D3, QuietMove),
				NewMove(D4, D2, QuietMove),
				NewMove(D4, D1, QuietMove),
				NewMove(D4, E4, QuietMove),
				NewMove(D4, F4, QuietMove),
				NewMove(D4, G4, QuietMove),
				NewMove(D4, H4, QuietMove),
				NewMove(D4, C4, QuietMove),
				NewMove(D4, B4, QuietMove),
				NewMove(D4, A4, QuietMove),
				// king on a1
				NewMove(A1, A2, QuietMove),
				NewMove(A1, B1, QuietMove),
				NewMove(A1, B2, QuietMove),
			},
		},
		{
			name: "bishop rays blocked by own king",
			fen:  "k7/8/8/8/3B4/8/8/K7 w - - 0 1",
			expectedMoves: []Move{
				// bishop on d4, north east
				NewMove(D4, E5, QuietMove),
				NewMove(D4, F6, QuietMove),
				NewMove(D4, G7, QuietMove),
				NewMove(D4, H8, QuietMove),
				// north west
				NewMove(D4, C5, QuietMove),
				NewMove(D4, B6, QuietMove),
				NewMove(D4, A7, QuietMove),
				// south east
				NewMove(D4, E3, QuietMove),
				NewMove(D4, F2, QuietMove),
				NewMove(D4, G1, QuietMove),
				// south west, stops short of the white king on a1
				NewMove(D4, C3, QuietMove),
				NewMove(D4, B2, QuietMove),
				// king on a1
				NewMove(A1, A2, QuietMove),
				NewMove(A1, B1, QuietMove),
				NewMove(A1, B2, QuietMove),
			},
		},
		{
			name: "promotions, quiet and capturing both ways",
			fen:  "r1n5/1P6/8/7k/8/8/8/4K3 w - - 0 1",
			expectedMoves: []Move{
				// quiet promotions onto the empty b8
				NewMove(B7, B8, KnightPromotion),
				NewMove(B7, B8, BishopPromotion),
				NewMove(B7, B8, RookPromotion),
				NewMove(B7, B8, QueenPromotion),
				// capturing the rook on a8
				NewMove(B7, A8, KnightPromotionCapture),
				NewMove(B7, A8, BishopPromotionCapture),
				NewMove(B7, A8, RookPromotionCapture),
				NewMove(B7, A8, QueenPromotionCapture),
				// capturing the knight on c8
				NewMove(B7, C8, KnightPromotionCapture),
				NewMove(B7, C8, BishopPromotionCapture),
				NewMove(B7, C8, RookPromotionCapture),
				NewMove(B7, C8, QueenPromotionCapture),
				// king on e1
				NewMove(E1, D1, QuietMove),
				NewMove(E1, D2, QuietMove),
				NewMove(E1, E2, QuietMove),
				NewMove(E1, F1, QuietMove),
				NewMove(E1, F2, QuietMove),
			},
		},
		{
			name: "en passant available",
			fen:  "7k/8/8/3pP3/8/8/8/K7 w - d6 0 1",
			expectedMoves: []Move{
				// pawn on e5
				NewMove(E5, E6, QuietMove),
				NewMove(E5, D6, EnPassantCapture),
				// king on a1
				NewMove(A1, A2, QuietMove),
				NewMove(A1, B1, QuietMove),
				NewMove(A1, B2, QuietMove),
			},
		},
		{
			name: "both castles available",
			fen:  "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			expectedMoves: []Move{
				// rook on a1, up the a file onto the black rook
				NewMove(A1, A2, QuietMove),
				NewMove(A1, A3, QuietMove),
				NewMove(A1, A4, QuietMove),
				NewMove(A1, A5, QuietMove),
				NewMove(A1, A6, QuietMove),
				NewMove(A1, A7, QuietMove),
				NewMove(A1, A8, Capture),
				// rook on a1, east until the king blocks
				NewMove(A1, B1, QuietMove),
				NewMove(A1, C1, QuietMove),
				NewMove(A1, D1, QuietMove),
				// rook on h1, up the h file onto the black rook
				NewMove(H1, H2, QuietMove),
				NewMove(H1, H3, QuietMove),
				NewMove(H1, H4, QuietMove),
				NewMove(H1, H5, QuietMove),
				NewMove(H1, H6, QuietMove),
				NewMove(H1, H7, QuietMove),
				NewMove(H1, H8, Capture),
				// rook on h1, west until the king blocks
				NewMove(H1, G1, QuietMove),
				NewMove(H1, F1, QuietMove),
				// king on e1
				NewMove(E1, D1, QuietMove),
				NewMove(E1, D2, QuietMove),
				NewMove(E1, E2, QuietMove),
				NewMove(E1, F1, QuietMove),
				NewMove(E1, F2, QuietMove),
				// castles
				NewMove(E1, G1, KingsideCastle),
				NewMove(E1, C1, QueensideCastle),
			},
		},
		{
			// the mirror of the position above, so that the black branch of the
			// castling code is exercised rather than only the white one
			name: "both castles available for black",
			fen:  "r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
			expectedMoves: []Move{
				// rook on a8, down the a file onto the white rook
				NewMove(A8, A7, QuietMove),
				NewMove(A8, A6, QuietMove),
				NewMove(A8, A5, QuietMove),
				NewMove(A8, A4, QuietMove),
				NewMove(A8, A3, QuietMove),
				NewMove(A8, A2, QuietMove),
				NewMove(A8, A1, Capture),
				// rook on a8, east until the king blocks
				NewMove(A8, B8, QuietMove),
				NewMove(A8, C8, QuietMove),
				NewMove(A8, D8, QuietMove),
				// rook on h8, down the h file onto the white rook
				NewMove(H8, H7, QuietMove),
				NewMove(H8, H6, QuietMove),
				NewMove(H8, H5, QuietMove),
				NewMove(H8, H4, QuietMove),
				NewMove(H8, H3, QuietMove),
				NewMove(H8, H2, QuietMove),
				NewMove(H8, H1, Capture),
				// rook on h8, west until the king blocks
				NewMove(H8, G8, QuietMove),
				NewMove(H8, F8, QuietMove),
				// king on e8
				NewMove(E8, D8, QuietMove),
				NewMove(E8, D7, QuietMove),
				NewMove(E8, E7, QuietMove),
				NewMove(E8, F8, QuietMove),
				NewMove(E8, F7, QuietMove),
				// castles
				NewMove(E8, G8, KingsideCastle),
				NewMove(E8, C8, QueensideCastle),
			},
		},
		{
			// two of the knight's eight targets hold enemy pawns, so each must
			// appear exactly once, as a capture and not also as a quiet move
			name: "knight with captures among its targets",
			fen:  "7k/8/2p1p3/8/3N4/8/8/K7 w - - 0 1",
			expectedMoves: []Move{
				// knight on d4, captures
				NewMove(D4, C6, Capture),
				NewMove(D4, E6, Capture),
				// knight on d4, quiet
				NewMove(D4, B5, QuietMove),
				NewMove(D4, F5, QuietMove),
				NewMove(D4, B3, QuietMove),
				NewMove(D4, F3, QuietMove),
				NewMove(D4, C2, QuietMove),
				NewMove(D4, E2, QuietMove),
				// king on a1
				NewMove(A1, A2, QuietMove),
				NewMove(A1, B1, QuietMove),
				NewMove(A1, B2, QuietMove),
			},
		},
		{
			// the king's capture loop is otherwise never exercised
			name: "king with captures among its targets",
			fen:  "7k/8/8/2n1n3/3K4/8/8/8 w - - 0 1",
			expectedMoves: []Move{
				NewMove(D4, C5, Capture),
				NewMove(D4, E5, Capture),
				NewMove(D4, C3, QuietMove),
				NewMove(D4, C4, QuietMove),
				NewMove(D4, D3, QuietMove),
				NewMove(D4, D5, QuietMove),
				NewMove(D4, E3, QuietMove),
				NewMove(D4, E4, QuietMove),
			},
		},
		{
			// a corner knight has only two targets and must not wrap
			name: "knight in the corner",
			fen:  "7k/8/8/8/4K3/8/8/N7 w - - 0 1",
			expectedMoves: []Move{
				NewMove(A1, B3, QuietMove),
				NewMove(A1, C2, QuietMove),
				NewMove(E4, D3, QuietMove),
				NewMove(E4, D4, QuietMove),
				NewMove(E4, D5, QuietMove),
				NewMove(E4, E3, QuietMove),
				NewMove(E4, E5, QuietMove),
				NewMove(E4, F3, QuietMove),
				NewMove(E4, F4, QuietMove),
				NewMove(E4, F5, QuietMove),
			},
		},
		{
			// in check: neither castle may be generated
			name: "no castling out of check",
			fen:  "4r3/8/2k5/8/8/8/8/R3K2R w KQ - 0 1",
			expectedMoves: []Move{
				NewMove(A1, A2, QuietMove),
				NewMove(A1, A3, QuietMove),
				NewMove(A1, A4, QuietMove),
				NewMove(A1, A5, QuietMove),
				NewMove(A1, A6, QuietMove),
				NewMove(A1, A7, QuietMove),
				NewMove(A1, A8, QuietMove),
				NewMove(A1, B1, QuietMove),
				NewMove(A1, C1, QuietMove),
				NewMove(A1, D1, QuietMove),
				NewMove(H1, H2, QuietMove),
				NewMove(H1, H3, QuietMove),
				NewMove(H1, H4, QuietMove),
				NewMove(H1, H5, QuietMove),
				NewMove(H1, H6, QuietMove),
				NewMove(H1, H7, QuietMove),
				NewMove(H1, H8, QuietMove),
				NewMove(H1, G1, QuietMove),
				NewMove(H1, F1, QuietMove),
				NewMove(E1, D1, QuietMove),
				NewMove(E1, D2, QuietMove),
				NewMove(E1, E2, QuietMove),
				NewMove(E1, F1, QuietMove),
				NewMove(E1, F2, QuietMove),
			},
		},
		{
			// f1 is attacked, so only the queenside castle survives
			name: "no castling through an attacked square",
			fen:  "2k2r2/8/8/8/8/8/8/R3K2R w KQ - 0 1",
			expectedMoves: []Move{
				NewMove(A1, A2, QuietMove),
				NewMove(A1, A3, QuietMove),
				NewMove(A1, A4, QuietMove),
				NewMove(A1, A5, QuietMove),
				NewMove(A1, A6, QuietMove),
				NewMove(A1, A7, QuietMove),
				NewMove(A1, A8, QuietMove),
				NewMove(A1, B1, QuietMove),
				NewMove(A1, C1, QuietMove),
				NewMove(A1, D1, QuietMove),
				NewMove(H1, H2, QuietMove),
				NewMove(H1, H3, QuietMove),
				NewMove(H1, H4, QuietMove),
				NewMove(H1, H5, QuietMove),
				NewMove(H1, H6, QuietMove),
				NewMove(H1, H7, QuietMove),
				NewMove(H1, H8, QuietMove),
				NewMove(H1, G1, QuietMove),
				NewMove(H1, F1, QuietMove),
				NewMove(E1, D1, QuietMove),
				NewMove(E1, D2, QuietMove),
				NewMove(E1, E2, QuietMove),
				NewMove(E1, F1, QuietMove),
				NewMove(E1, F2, QuietMove),
				NewMove(E1, C1, QueensideCastle),
			},
		},
		{
			// b1 is attacked but empty; the king never crosses it, so both
			// castles remain legal. Requiring b1 unattacked drops this to 25.
			name: "castling legal with b1 attacked",
			fen:  "1r2k3/8/8/8/8/8/8/R3K2R w KQ - 0 1",
			expectedMoves: []Move{
				NewMove(A1, A2, QuietMove),
				NewMove(A1, A3, QuietMove),
				NewMove(A1, A4, QuietMove),
				NewMove(A1, A5, QuietMove),
				NewMove(A1, A6, QuietMove),
				NewMove(A1, A7, QuietMove),
				NewMove(A1, A8, QuietMove),
				NewMove(A1, B1, QuietMove),
				NewMove(A1, C1, QuietMove),
				NewMove(A1, D1, QuietMove),
				NewMove(H1, H2, QuietMove),
				NewMove(H1, H3, QuietMove),
				NewMove(H1, H4, QuietMove),
				NewMove(H1, H5, QuietMove),
				NewMove(H1, H6, QuietMove),
				NewMove(H1, H7, QuietMove),
				NewMove(H1, H8, QuietMove),
				NewMove(H1, G1, QuietMove),
				NewMove(H1, F1, QuietMove),
				NewMove(E1, D1, QuietMove),
				NewMove(E1, D2, QuietMove),
				NewMove(E1, E2, QuietMove),
				NewMove(E1, F1, QuietMove),
				NewMove(E1, F2, QuietMove),
				NewMove(E1, G1, KingsideCastle),
				NewMove(E1, C1, QueensideCastle),
			},
		},
		{
			// own bishop still on f1 blocks the kingside castle
			name: "no castling with a piece in the way",
			fen:  "4k3/8/8/8/8/8/8/R3KB1R w KQ - 0 1",
			expectedMoves: []Move{
				NewMove(A1, A2, QuietMove),
				NewMove(A1, A3, QuietMove),
				NewMove(A1, A4, QuietMove),
				NewMove(A1, A5, QuietMove),
				NewMove(A1, A6, QuietMove),
				NewMove(A1, A7, QuietMove),
				NewMove(A1, A8, QuietMove),
				NewMove(A1, B1, QuietMove),
				NewMove(A1, C1, QuietMove),
				NewMove(A1, D1, QuietMove),
				NewMove(H1, H2, QuietMove),
				NewMove(H1, H3, QuietMove),
				NewMove(H1, H4, QuietMove),
				NewMove(H1, H5, QuietMove),
				NewMove(H1, H6, QuietMove),
				NewMove(H1, H7, QuietMove),
				NewMove(H1, H8, QuietMove),
				NewMove(H1, G1, QuietMove),
				NewMove(F1, G2, QuietMove),
				NewMove(F1, H3, QuietMove),
				NewMove(F1, E2, QuietMove),
				NewMove(F1, D3, QuietMove),
				NewMove(F1, C4, QuietMove),
				NewMove(F1, B5, QuietMove),
				NewMove(F1, A6, QuietMove),
				NewMove(E1, D1, QuietMove),
				NewMove(E1, D2, QuietMove),
				NewMove(E1, E2, QuietMove),
				NewMove(E1, F2, QuietMove),
				NewMove(E1, C1, QueensideCastle),
			},
		},
		{
			// same shape, rights already spent
			name: "no castling without rights",
			fen:  "4k3/8/8/8/8/8/8/R3K2R w - - 0 1",
			expectedMoves: []Move{
				NewMove(A1, A2, QuietMove),
				NewMove(A1, A3, QuietMove),
				NewMove(A1, A4, QuietMove),
				NewMove(A1, A5, QuietMove),
				NewMove(A1, A6, QuietMove),
				NewMove(A1, A7, QuietMove),
				NewMove(A1, A8, QuietMove),
				NewMove(A1, B1, QuietMove),
				NewMove(A1, C1, QuietMove),
				NewMove(A1, D1, QuietMove),
				NewMove(H1, H2, QuietMove),
				NewMove(H1, H3, QuietMove),
				NewMove(H1, H4, QuietMove),
				NewMove(H1, H5, QuietMove),
				NewMove(H1, H6, QuietMove),
				NewMove(H1, H7, QuietMove),
				NewMove(H1, H8, QuietMove),
				NewMove(H1, G1, QuietMove),
				NewMove(H1, F1, QuietMove),
				NewMove(E1, D1, QuietMove),
				NewMove(E1, D2, QuietMove),
				NewMove(E1, E2, QuietMove),
				NewMove(E1, F1, QuietMove),
				NewMove(E1, F2, QuietMove),
			},
		},
		{
			// no other fixture contains a queen at all
			name: "queen on all eight rays",
			fen:  "k7/8/8/8/3Q4/8/8/K7 w - - 0 1",
			expectedMoves: []Move{
				// north
				NewMove(D4, D5, QuietMove),
				NewMove(D4, D6, QuietMove),
				NewMove(D4, D7, QuietMove),
				NewMove(D4, D8, QuietMove),
				// south
				NewMove(D4, D3, QuietMove),
				NewMove(D4, D2, QuietMove),
				NewMove(D4, D1, QuietMove),
				// east
				NewMove(D4, E4, QuietMove),
				NewMove(D4, F4, QuietMove),
				NewMove(D4, G4, QuietMove),
				NewMove(D4, H4, QuietMove),
				// west
				NewMove(D4, C4, QuietMove),
				NewMove(D4, B4, QuietMove),
				NewMove(D4, A4, QuietMove),
				// north east
				NewMove(D4, E5, QuietMove),
				NewMove(D4, F6, QuietMove),
				NewMove(D4, G7, QuietMove),
				NewMove(D4, H8, QuietMove),
				// north west
				NewMove(D4, C5, QuietMove),
				NewMove(D4, B6, QuietMove),
				NewMove(D4, A7, QuietMove),
				// south east
				NewMove(D4, E3, QuietMove),
				NewMove(D4, F2, QuietMove),
				NewMove(D4, G1, QuietMove),
				// south west, stops short of the friendly king on a1
				NewMove(D4, C3, QuietMove),
				NewMove(D4, B2, QuietMove),
				// king on a1
				NewMove(A1, A2, QuietMove),
				NewMove(A1, B1, QuietMove),
				NewMove(A1, B2, QuietMove),
			},
		},
		{
			// b2 cannot push at all; c2 can push but its double is blocked on c4
			name: "blocked pawn pushes",
			fen:  "7k/8/8/8/2n5/1n6/PPP5/K7 w - - 0 1",
			expectedMoves: []Move{
				NewMove(A2, A3, QuietMove),
				NewMove(A2, A4, DoublePawnPush),
				NewMove(A2, B3, Capture),
				NewMove(C2, C3, QuietMove),
				NewMove(C2, B3, Capture),
				NewMove(A1, B1, QuietMove),
			},
		},
		{
			// the mirror, so black pawn direction is exercised
			name: "blocked pawn pushes for black",
			fen:  "k7/ppp5/1N6/2N5/8/8/8/7K b - - 0 1",
			expectedMoves: []Move{
				NewMove(A7, A6, QuietMove),
				NewMove(A7, A5, DoublePawnPush),
				NewMove(A7, B6, Capture),
				NewMove(C7, C6, QuietMove),
				NewMove(C7, B6, Capture),
				NewMove(A8, B8, QuietMove),
			},
		},
		{
			// an a file pawn must not capture "west" onto the h file
			name: "pawn capture does not wrap the board",
			fen:  "7k/8/8/8/8/7n/P7/K7 w - - 0 1",
			expectedMoves: []Move{
				NewMove(A2, A3, QuietMove),
				NewMove(A2, A4, DoublePawnPush),
				NewMove(A1, B1, QuietMove),
				NewMove(A1, B2, QuietMove),
			},
		},
		{
			// black capturing en passant, the mirror of the white case above
			name: "en passant available for black",
			fen:  "7k/8/8/8/3Pp3/8/8/K7 b - d3 0 1",
			expectedMoves: []Move{
				NewMove(E4, E3, QuietMove),
				NewMove(E4, D3, EnPassantCapture),
				NewMove(H8, G8, QuietMove),
				NewMove(H8, G7, QuietMove),
				NewMove(H8, H7, QuietMove),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pos, err := ParseFEN(test.fen)
			if err != nil {
				t.Fatalf("Failed to parse FEN: %v", err)
			}

			moves := make([]Move, 0, 256)
			moves = pos.GeneratePseudoLegalMoves(moves)
			if len(moves) != len(test.expectedMoves) {
				t.Errorf("Expected %d moves, got %d", len(test.expectedMoves), len(moves))
			}

			expectedMap := make(map[Move]bool, len(test.expectedMoves))
			for _, expectedMove := range test.expectedMoves {
				expectedMap[expectedMove] = true
			}

			moveMap := make(map[Move]bool)
			for _, move := range moves {
				if _, ok := expectedMap[move]; !ok {
					t.Errorf("Unexpected move generated: %v", move)
				}
				moveMap[move] = true
			}

			for _, expectedMove := range test.expectedMoves {
				if !moveMap[expectedMove] {
					t.Errorf("Expected move %v not found in generated moves", expectedMove)
				}
				if !pos.IsPseudoLegalMove(expectedMove) {
					t.Errorf("Expected move %v is not pseudo-legal", expectedMove)
				}
			}
		})
	}
}

func TestGenerateLegalMoves(t *testing.T) {
	// The pseudo-legal count is asserted alongside the legal moves so that a failure
	// distinguishes "the filter removed the wrong moves" from "the generator never
	// produced them in the first place".
	tests := []struct {
		name          string
		fen           string
		pseudoLegal   int
		expectedMoves []Move
	}{
		{
			// the knight on e2 is pinned against its own king by the rook on e8. a
			// knight can never move along the line it is pinned on, so all six of its
			// pseudo-legal moves are illegal and only the king's survive
			name:        "pinned knight has no legal moves",
			fen:         "4r2k/8/8/8/8/8/4N3/4K3 w - - 0 1",
			pseudoLegal: 10,
			expectedMoves: []Move{
				NewMove(E1, D1, QuietMove),
				NewMove(E1, D2, QuietMove),
				NewMove(E1, F1, QuietMove),
				NewMove(E1, F2, QuietMove),
			},
		},
		{
			// white is in check from the rook on e8. the king may step off the e file
			// but not along it, and of the rook's fourteen moves only Re2 blocks
			name:        "in check, only evasions survive",
			fen:         "4r2k/8/8/8/8/8/R7/4K3 w - - 0 1",
			pseudoLegal: 19,
			expectedMoves: []Move{
				NewMove(E1, D1, QuietMove),
				NewMove(E1, D2, QuietMove),
				NewMove(E1, F1, QuietMove),
				NewMove(E1, F2, QuietMove),
				NewMove(A2, E2, QuietMove),
			},
		},
		{
			// bxc3 en passant empties both b4 and c4, opening the fourth rank so that
			// the queen on h4 attacks the black king on a4. removing two pawns from one
			// rank is the case no pin mask catches, which is why the reference
			// implementation makes the move to find out. Kb5 is also illegal, covered
			// by the pawn on c4
			name:        "en passant capture would expose the king",
			fen:         "8/8/8/8/kpP4Q/8/8/K7 b - c3 0 1",
			pseudoLegal: 6,
			expectedMoves: []Move{
				NewMove(A4, A3, QuietMove),
				NewMove(A4, A5, QuietMove),
				NewMove(A4, B3, QuietMove),
				NewMove(B4, B3, QuietMove),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pos, err := ParseFEN(test.fen)
			if err != nil {
				t.Fatalf("Failed to parse FEN: %v", err)
			}

			pseudoLegal := pos.GeneratePseudoLegalMoves(make([]Move, 0, 256))
			if len(pseudoLegal) != test.pseudoLegal {
				t.Errorf("Expected %d pseudo-legal moves, got %d", test.pseudoLegal, len(pseudoLegal))
			}

			before := GenerateFEN(pos)
			moves := pos.GenerateLegalMoves()
			if len(moves) != len(test.expectedMoves) {
				t.Errorf("Expected %d legal moves, got %d", len(test.expectedMoves), len(moves))
			}

			expectedMap := make(map[Move]bool, len(test.expectedMoves))
			for _, expectedMove := range test.expectedMoves {
				expectedMap[expectedMove] = true
			}

			moveMap := make(map[Move]bool)
			for _, move := range moves {
				if _, ok := expectedMap[move]; !ok {
					t.Errorf("Illegal move generated: %v", move)
				}
				moveMap[move] = true
			}

			for _, expectedMove := range test.expectedMoves {
				if !moveMap[expectedMove] {
					t.Errorf("Expected move %v not found in generated moves", expectedMove)
				}
			}

			// every pseudo-legal move was made and unmade, so an asymmetry between the
			// two shows up here rather than as a perft mismatch later
			if after := GenerateFEN(pos); after != before {
				t.Errorf("Position changed by GenerateLegalMoves:\n before: %s\n after:  %s", before, after)
			}
			if err := pos.validateStateSlow(); err != nil {
				t.Errorf("Position state corrupt after GenerateLegalMoves: %v", err)
			}
		})
	}
}

func BenchmarkGeneratePseudoLegalMoves(b *testing.B) {
	for _, test := range perftCases {
		b.Run(test.name, func(b *testing.B) {
			pos, err := ParseFEN(test.fen)
			if err != nil {
				b.Fatalf("Failed to parse FEN: %v", err)
			}
			moves := make([]Move, 0, 256)
			b.ResetTimer()
			for b.Loop() {
				moves = pos.GeneratePseudoLegalMoves(moves)
			}
			b.StopTimer()
			// use moves to stop it being optimised out by the compiler
			if len(moves) == 0 {
				b.Error("No moves generated")
			}
		})
	}
}

func TestGeneratePseudoLegalCapturesAndPromotions(t *testing.T) {
	tests := []struct {
		name          string
		fen           string
		expectedMoves []Move
	}{
		{
			name:          "initial position",
			fen:           "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
			expectedMoves: []Move{},
		},
		{
			name:          "lone knight in the centre",
			fen:           "7k/8/8/8/3N4/8/8/K7 w - - 0 1",
			expectedMoves: []Move{},
		},
		{
			name: "rook rays blocked by an enemy pawn",
			fen:  "7k/8/8/3p4/3R4/8/8/K7 w - - 0 1",
			expectedMoves: []Move{
				// rook on d4, north ray ends on the capture
				NewMove(D4, D5, Capture),
			},
		},
		{
			name:          "bishop rays blocked by own king",
			fen:           "k7/8/8/8/3B4/8/8/K7 w - - 0 1",
			expectedMoves: []Move{},
		},
		{
			name: "promotions, quiet and capturing both ways",
			fen:  "r1n5/1P6/8/7k/8/8/8/4K3 w - - 0 1",
			expectedMoves: []Move{
				// quiet promotions onto the empty b8
				NewMove(B7, B8, KnightPromotion),
				NewMove(B7, B8, BishopPromotion),
				NewMove(B7, B8, RookPromotion),
				NewMove(B7, B8, QueenPromotion),
				// capturing the rook on a8
				NewMove(B7, A8, KnightPromotionCapture),
				NewMove(B7, A8, BishopPromotionCapture),
				NewMove(B7, A8, RookPromotionCapture),
				NewMove(B7, A8, QueenPromotionCapture),
				// capturing the knight on c8
				NewMove(B7, C8, KnightPromotionCapture),
				NewMove(B7, C8, BishopPromotionCapture),
				NewMove(B7, C8, RookPromotionCapture),
				NewMove(B7, C8, QueenPromotionCapture),
			},
		},
		{
			name: "en passant available",
			fen:  "7k/8/8/3pP3/8/8/8/K7 w - d6 0 1",
			expectedMoves: []Move{
				NewMove(E5, D6, EnPassantCapture),
			},
		},
		{
			name: "both castles available",
			fen:  "r3k2r/8/8/8/8/8/8/R3K2R w KQkq - 0 1",
			expectedMoves: []Move{
				NewMove(A1, A8, Capture),
				NewMove(H1, H8, Capture),
			},
		},
		{
			// the mirror of the position above, so that the black branch of the
			// castling code is exercised rather than only the white one
			name: "both castles available for black",
			fen:  "r3k2r/8/8/8/8/8/8/R3K2R b KQkq - 0 1",
			expectedMoves: []Move{
				NewMove(A8, A1, Capture),
				NewMove(H8, H1, Capture),
			},
		},
		{
			// two of the knight's eight targets hold enemy pawns, so each must
			// appear exactly once, as a capture and not also as a quiet move
			name: "knight with captures among its targets",
			fen:  "7k/8/2p1p3/8/3N4/8/8/K7 w - - 0 1",
			expectedMoves: []Move{
				NewMove(D4, C6, Capture),
				NewMove(D4, E6, Capture),
			},
		},
		{
			// the king's capture loop is otherwise never exercised
			name: "king with captures among its targets",
			fen:  "7k/8/8/2n1n3/3K4/8/8/8 w - - 0 1",
			expectedMoves: []Move{
				NewMove(D4, C5, Capture),
				NewMove(D4, E5, Capture),
			},
		},
		{
			// b2 cannot push at all; c2 can push but its double is blocked on c4
			name: "blocked pawn pushes",
			fen:  "7k/8/8/8/2n5/1n6/PPP5/K7 w - - 0 1",
			expectedMoves: []Move{
				NewMove(A2, B3, Capture),
				NewMove(C2, B3, Capture),
			},
		},
		{
			// the mirror, so black pawn direction is exercised
			name: "blocked pawn pushes for black",
			fen:  "k7/ppp5/1N6/2N5/8/8/8/7K b - - 0 1",
			expectedMoves: []Move{
				NewMove(A7, B6, Capture),
				NewMove(C7, B6, Capture),
			},
		},
		{
			// an a file pawn must not capture "west" onto the h file
			name:          "pawn capture does not wrap the board",
			fen:           "7k/8/8/8/8/7n/P7/K7 w - - 0 1",
			expectedMoves: []Move{},
		},
		{
			// black capturing en passant, the mirror of the white case above
			name: "en passant available for black",
			fen:  "7k/8/8/8/3Pp3/8/8/K7 b - d3 0 1",
			expectedMoves: []Move{
				NewMove(E4, D3, EnPassantCapture),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pos, err := ParseFEN(test.fen)
			if err != nil {
				t.Fatalf("Failed to parse FEN: %v", err)
			}

			moves := make([]Move, 0, 256)
			moves = pos.GeneratePseudoLegalCapturesAndPromotions(moves)
			if len(moves) != len(test.expectedMoves) {
				t.Errorf("Expected %d captures/promos, got %d", len(test.expectedMoves), len(moves))
			}

			expectedMap := make(map[Move]bool, len(test.expectedMoves))
			for _, expectedMove := range test.expectedMoves {
				expectedMap[expectedMove] = true
			}

			moveMap := make(map[Move]bool)
			for _, move := range moves {
				if _, ok := expectedMap[move]; !ok {
					t.Errorf("Unexpected move generated: %v", move)
				}
				moveMap[move] = true
			}

			for _, expectedMove := range test.expectedMoves {
				if !moveMap[expectedMove] {
					t.Errorf("Expected move %v not found in generated moves", expectedMove)
				}
			}
		})
	}
}
