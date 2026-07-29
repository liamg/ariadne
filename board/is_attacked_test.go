package board

import "testing"

func TestIsSquareAttacked(t *testing.T) {
	// pawns face each other on the e file, kings tucked away on e1 and e8
	const facingPawns = "4k3/8/8/4p3/4P3/8/8/4K3 w - - 0 1"

	tests := []struct {
		name     string
		fen      string
		square   Square
		by       Colour
		expected bool
	}{
		// Pawn attacks are not symmetric, so the lookup has to radiate away from
		// the target square in the direction the attacker came from.
		{"white pawn attacks up and left", facingPawns, D5, White, true},
		{"white pawn attacks up and right", facingPawns, F5, White, true},
		{"white pawn does not attack straight ahead", facingPawns, E5, White, false},
		{"white pawn does not attack behind itself", facingPawns, D3, White, false},
		{"black pawn attacks down and left", facingPawns, D4, Black, true},
		{"black pawn attacks down and right", facingPawns, F4, Black, true},
		{"black pawn does not attack straight ahead", facingPawns, E4, Black, false},
		{"black pawn does not attack behind itself", facingPawns, D6, Black, false},
		{"square attacked by white is not attacked by black", facingPawns, D5, Black, false},

		// A pawn on the a file must not wrap around onto the h file.
		{"edge pawn attacks inward", "4k3/8/8/8/P7/8/8/4K3 w - - 0 1", B5, White, true},
		{"edge pawn does not wrap", "4k3/8/8/8/P7/8/8/4K3 w - - 0 1", H5, White, false},

		// Knight attacks are symmetric, so the square can be probed directly.
		{"knight attacks two up one across", "4k3/8/8/8/3N4/8/8/4K3 w - - 0 1", C6, White, true},
		{"knight attacks two down one across", "4k3/8/8/8/3N4/8/8/4K3 w - - 0 1", B3, White, true},
		{"knight does not attack adjacent square", "4k3/8/8/8/3N4/8/8/4K3 w - - 0 1", D5, White, false},

		// Kings, likewise symmetric.
		{"king attacks diagonally", "8/8/8/8/8/8/8/4K2k w - - 0 1", D2, White, true},
		{"king does not attack two ranks away", "8/8/8/8/8/8/8/4K2k w - - 0 1", D3, White, false},
		{"black king attacks its neighbours", "8/8/8/8/8/8/8/4K2k w - - 0 1", G2, Black, true},
		{"black king does not attack two files away", "8/8/8/8/8/8/8/4K2k w - - 0 1", F1, Black, false},

		// Sliding attacks depend on occupancy: a ray reaches the first occupied
		// square and stops, and a friendly piece blocks just as an enemy does.
		{"rook attacks along a clear file", "4k3/8/8/8/8/8/8/R3K3 w - - 0 1", A8, White, true},
		{"rook attacks along a clear rank", "4k3/8/8/8/8/8/8/R3K3 w - - 0 1", C1, White, true},
		{"rook ray stops at its own king", "4k3/8/8/8/8/8/8/R3K3 w - - 0 1", G1, White, false},
		{"rook attacks the square of the blocker", "4k3/8/8/8/p7/8/8/R3K3 w - - 0 1", A4, White, true},
		{"rook ray stops at an enemy blocker", "4k3/8/8/8/p7/8/8/R3K3 w - - 0 1", A8, White, false},
		{"bishop attacks along a clear diagonal", "4k3/8/8/8/8/8/8/2B1K3 w - - 0 1", H6, White, true},
		{"bishop attacks the short diagonal", "4k3/8/8/8/8/8/8/2B1K3 w - - 0 1", A3, White, true},
		{"bishop ray stops at a blocker", "4k3/8/8/8/5P2/8/8/2B1K3 w - - 0 1", H6, White, false},
		// g5 is beyond the blocked bishop ray, but the blocking pawn attacks it
		{"blocker itself still attacks", "4k3/8/8/8/5P2/8/8/2B1K3 w - - 0 1", G5, White, true},

		// The queen shares the bishop and rook ray lookups.
		{"queen attacks along a file", "4k3/8/8/3Q4/8/8/8/4K3 w - - 0 1", D8, White, true},
		{"queen attacks along a rank", "4k3/8/8/8/8/8/8/4K3 w - - 0 1", A5, White, false},
		{"queen attacks along a rank with queen present", "4k3/8/8/3Q4/8/8/8/4K3 w - - 0 1", A5, White, true},
		{"queen attacks down a diagonal", "4k3/8/8/3Q4/8/8/8/4K3 w - - 0 1", H1, White, true},
		{"queen attacks up a diagonal", "4k3/8/8/3Q4/8/8/8/4K3 w - - 0 1", A8, White, true},
		{"queen does not attack off its lines", "4k3/8/8/3Q4/8/8/8/4K3 w - - 0 1", C3, White, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pos, err := ParseFEN(test.fen)
			if err != nil {
				t.Fatalf("Failed to parse FEN: %v", err)
			}

			if got := pos.isSquareAttacked(test.square, test.by); got != test.expected {
				t.Errorf("isSquareAttacked(%s, %s) = %v, want %v",
					test.square, test.by, got, test.expected)
			}
		})
	}
}
