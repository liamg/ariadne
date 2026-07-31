package board

import (
	"fmt"
	"slices"
	"testing"
)

var perftCases = []struct {
	name                     string
	fen                      string
	expectedDepths           []uint64
	expectedDepthThreeCounts []PerftCount
}{
	{
		name:           "Initial Position",
		fen:            "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		expectedDepths: []uint64{1, 20, 400, 8902, 197281, 4865609, 119060324},
		expectedDepthThreeCounts: []PerftCount{
			{NewMove(A2, A3, QuietMove), 380},
			{NewMove(A2, A4, DoublePawnPush), 420},
			{NewMove(B1, A3, QuietMove), 400},
			{NewMove(B1, C3, QuietMove), 440},
			{NewMove(B2, B3, QuietMove), 420},
			{NewMove(B2, B4, DoublePawnPush), 421},
			{NewMove(C2, C3, QuietMove), 420},
			{NewMove(C2, C4, DoublePawnPush), 441},
			{NewMove(D2, D3, QuietMove), 539},
			{NewMove(D2, D4, DoublePawnPush), 560},
			{NewMove(E2, E3, QuietMove), 599},
			{NewMove(E2, E4, DoublePawnPush), 600},
			{NewMove(F2, F3, QuietMove), 380},
			{NewMove(F2, F4, DoublePawnPush), 401},
			{NewMove(G1, F3, QuietMove), 440},
			{NewMove(G1, H3, QuietMove), 400},
			{NewMove(G2, G3, QuietMove), 420},
			{NewMove(G2, G4, DoublePawnPush), 421},
			{NewMove(H2, H3, QuietMove), 380},
			{NewMove(H2, H4, DoublePawnPush), 420},
		},
	},
	{
		name:           "Kiwipete",
		fen:            FenKiwiPete,
		expectedDepths: []uint64{1, 48, 2039, 97862, 4085603, 193690690},
		expectedDepthThreeCounts: []PerftCount{
			{NewMove(A1, B1, QuietMove), 1969},
			{NewMove(A1, C1, QuietMove), 1968},
			{NewMove(A1, D1, QuietMove), 1885},
			{NewMove(A2, A3, QuietMove), 2186},
			{NewMove(A2, A4, DoublePawnPush), 2149},
			{NewMove(B2, B3, QuietMove), 1964},
			{NewMove(C3, A4, QuietMove), 2203},
			{NewMove(C3, B1, QuietMove), 2038},
			{NewMove(C3, B5, QuietMove), 2138},
			{NewMove(C3, D1, QuietMove), 2040},
			{NewMove(D2, C1, QuietMove), 1963},
			{NewMove(D2, E3, QuietMove), 2136},
			{NewMove(D2, F4, QuietMove), 2000},
			{NewMove(D2, G5, QuietMove), 2134},
			{NewMove(D2, H6, QuietMove), 2019},
			{NewMove(D5, D6, QuietMove), 1991},
			{NewMove(D5, E6, Capture), 2241},
			{NewMove(E1, C1, QueensideCastle), 1887},
			{NewMove(E1, D1, QuietMove), 1894},
			{NewMove(E1, F1, QuietMove), 1855},
			{NewMove(E1, G1, KingsideCastle), 2059},
			{NewMove(E2, A6, Capture), 1907},
			{NewMove(E2, B5, QuietMove), 2057},
			{NewMove(E2, C4, QuietMove), 2082},
			{NewMove(E2, D1, QuietMove), 1733},
			{NewMove(E2, D3, QuietMove), 2050},
			{NewMove(E2, F1, QuietMove), 2060},
			{NewMove(E5, C4, QuietMove), 1880},
			{NewMove(E5, C6, QuietMove), 2027},
			{NewMove(E5, D3, QuietMove), 1803},
			{NewMove(E5, D7, Capture), 2124},
			{NewMove(E5, F7, Capture), 2080},
			{NewMove(E5, G4, QuietMove), 1878},
			{NewMove(E5, G6, Capture), 1997},
			{NewMove(F3, D3, QuietMove), 2005},
			{NewMove(F3, E3, QuietMove), 2174},
			{NewMove(F3, F4, QuietMove), 2132},
			{NewMove(F3, F5, QuietMove), 2396},
			{NewMove(F3, F6, Capture), 2111},
			{NewMove(F3, G3, QuietMove), 2214},
			{NewMove(F3, G4, QuietMove), 2169},
			{NewMove(F3, H3, Capture), 2360},
			{NewMove(F3, H5, QuietMove), 2267},
			{NewMove(G2, G3, QuietMove), 1882},
			{NewMove(G2, G4, DoublePawnPush), 1843},
			{NewMove(G2, H3, Capture), 1970},
			{NewMove(H1, F1, QuietMove), 1929},
			{NewMove(H1, G1, QuietMove), 2013},
		},
	},
	{
		name:           "Position 3",
		fen:            "8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",
		expectedDepths: []uint64{1, 14, 191, 2812, 43238, 674624, 11030083},
		expectedDepthThreeCounts: []PerftCount{
			{NewMove(A5, A4, QuietMove), 224},
			{NewMove(A5, A6, QuietMove), 240},
			{NewMove(B4, A4, QuietMove), 202},
			{NewMove(B4, B1, QuietMove), 265},
			{NewMove(B4, B2, QuietMove), 205},
			{NewMove(B4, B3, QuietMove), 248},
			{NewMove(B4, C4, QuietMove), 254},
			{NewMove(B4, D4, QuietMove), 243},
			{NewMove(B4, E4, QuietMove), 228},
			{NewMove(B4, F4, Capture), 41},
			{NewMove(E2, E3, QuietMove), 205},
			{NewMove(E2, E4, DoublePawnPush), 177},
			{NewMove(G2, G3, QuietMove), 54},
			{NewMove(G2, G4, DoublePawnPush), 226},
		},
	},
	{
		name:           "Position 4",
		fen:            "r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",
		expectedDepths: []uint64{1, 6, 264, 9467, 422333, 15833292},
		expectedDepthThreeCounts: []PerftCount{
			{NewMove(B4, C5, QuietMove), 1352},
			{NewMove(C4, C5, QuietMove), 1409},
			{NewMove(D2, D4, DoublePawnPush), 1643},
			{NewMove(F1, F2, QuietMove), 1623},
			{NewMove(F3, D4, QuietMove), 1687},
			{NewMove(G1, H1, QuietMove), 1753},
		},
	},
	{
		name:           "Position 5",
		fen:            "rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",
		expectedDepths: []uint64{1, 44, 1486, 62379, 2103487, 89941194},
		expectedDepthThreeCounts: []PerftCount{
			{NewMove(A2, A3, QuietMove), 1373},
			{NewMove(A2, A4, DoublePawnPush), 1433},
			{NewMove(B1, A3, QuietMove), 1303},
			{NewMove(B1, C3, QuietMove), 1467},
			{NewMove(B1, D2, QuietMove), 1174},
			{NewMove(B2, B3, QuietMove), 1368},
			{NewMove(B2, B4, DoublePawnPush), 1398},
			{NewMove(C1, D2, QuietMove), 1368},
			{NewMove(C1, E3, QuietMove), 1587},
			{NewMove(C1, F4, QuietMove), 1552},
			{NewMove(C1, G5, QuietMove), 1422},
			{NewMove(C1, H6, QuietMove), 1312},
			{NewMove(C2, C3, QuietMove), 1440},
			{NewMove(C4, A6, QuietMove), 1256},
			{NewMove(C4, B3, QuietMove), 1275},
			{NewMove(C4, B5, QuietMove), 1332},
			{NewMove(C4, D3, QuietMove), 1269},
			{NewMove(C4, D5, QuietMove), 1375},
			{NewMove(C4, E6, QuietMove), 1438},
			{NewMove(C4, F7, Capture), 1328},
			{NewMove(D1, D2, QuietMove), 1436},
			{NewMove(D1, D3, QuietMove), 1685},
			{NewMove(D1, D4, QuietMove), 1751},
			{NewMove(D1, D5, QuietMove), 1688},
			{NewMove(D1, D6, QuietMove), 1500},
			{NewMove(D7, C8, BishopPromotionCapture), 1668},
			{NewMove(D7, C8, KnightPromotionCapture), 1607},
			{NewMove(D7, C8, QueenPromotionCapture), 1459},
			{NewMove(D7, C8, RookPromotionCapture), 1296},
			{NewMove(E1, D2, QuietMove), 978},
			{NewMove(E1, F1, QuietMove), 1445},
			{NewMove(E1, F2, Capture), 1269},
			{NewMove(E1, G1, KingsideCastle), 1376},
			{NewMove(E2, C3, QuietMove), 1595},
			{NewMove(E2, D4, QuietMove), 1554},
			{NewMove(E2, F4, QuietMove), 1555},
			{NewMove(E2, G1, QuietMove), 1431},
			{NewMove(E2, G3, QuietMove), 1523},
			{NewMove(G2, G3, QuietMove), 1308},
			{NewMove(G2, G4, DoublePawnPush), 1337},
			{NewMove(H1, F1, QuietMove), 1364},
			{NewMove(H1, G1, QuietMove), 1311},
			{NewMove(H2, H3, QuietMove), 1371},
			{NewMove(H2, H4, DoublePawnPush), 1402},
		},
	},
	{
		name:           "Position 6",
		fen:            "r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10",
		expectedDepths: []uint64{1, 46, 2079, 89890, 3894594, 164075551},
		expectedDepthThreeCounts: []PerftCount{
			{NewMove(A1, A2, QuietMove), 1854},
			{NewMove(A1, B1, QuietMove), 1943},
			{NewMove(A1, C1, QuietMove), 1899},
			{NewMove(A1, D1, QuietMove), 1853},
			{NewMove(A1, E1, QuietMove), 1764},
			{NewMove(A3, A4, QuietMove), 2076},
			{NewMove(B2, B3, QuietMove), 1943},
			{NewMove(B2, B4, DoublePawnPush), 2027},
			{NewMove(C3, A2, QuietMove), 1897},
			{NewMove(C3, A4, QuietMove), 1944},
			{NewMove(C3, B1, QuietMove), 1719},
			{NewMove(C3, B5, QuietMove), 1986},
			{NewMove(C3, D1, QuietMove), 1674},
			{NewMove(C3, D5, QuietMove), 2030},
			{NewMove(C4, A2, QuietMove), 1944},
			{NewMove(C4, A6, Capture), 1943},
			{NewMove(C4, B3, QuietMove), 1947},
			{NewMove(C4, B5, QuietMove), 1915},
			{NewMove(C4, D5, QuietMove), 1951},
			{NewMove(C4, E6, QuietMove), 2093},
			{NewMove(C4, F7, Capture), 165},
			{NewMove(D3, D4, QuietMove), 2202},
			{NewMove(E2, D1, QuietMove), 1897},
			{NewMove(E2, D2, QuietMove), 2079},
			{NewMove(E2, E1, QuietMove), 1944},
			{NewMove(E2, E3, QuietMove), 2121},
			{NewMove(F1, B1, QuietMove), 1944},
			{NewMove(F1, C1, QuietMove), 1990},
			{NewMove(F1, D1, QuietMove), 2034},
			{NewMove(F1, E1, QuietMove), 2035},
			{NewMove(F3, D2, QuietMove), 2002},
			{NewMove(F3, D4, QuietMove), 2293},
			{NewMove(F3, E1, QuietMove), 1769},
			{NewMove(F3, E5, Capture), 2403},
			{NewMove(F3, H4, QuietMove), 2045},
			{NewMove(G1, H1, QuietMove), 2211},
			{NewMove(G2, G3, QuietMove), 2075},
			{NewMove(G5, C1, QuietMove), 1841},
			{NewMove(G5, D2, QuietMove), 2026},
			{NewMove(G5, E3, QuietMove), 2114},
			{NewMove(G5, F4, QuietMove), 2205},
			{NewMove(G5, F6, Capture), 1933},
			{NewMove(G5, H4, QuietMove), 1898},
			{NewMove(G5, H6, QuietMove), 2065},
			{NewMove(H2, H3, QuietMove), 2163},
			{NewMove(H2, H4, DoublePawnPush), 2034},
		},
	},
}

// perftEdgeCases isolate individual move generation rules that the six standard
// positions can mask, since two compensating errors there can still sum to the
// right total. Each has a single published count. All are sub-second.
var perftEdgeCases = []struct {
	name          string
	fen           string
	depth         int
	expectedNodes uint64
}{
	// En passant legality: the capture is unavailable because making it would
	// expose the king along the rank the captured pawn was shielding.
	{"Illegal en passant #1", "3k4/3p4/8/K1P4r/8/8/8/8 b - - 0 1", 6, 1134888},
	{"Illegal en passant #2", "8/8/4k3/8/2p5/8/B2P2K1/8 w - - 0 1", 6, 1015133},
	// An en passant capture that delivers check to the opponent.
	{"En passant capture gives check", "8/8/1k6/2b5/2pP4/8/5K2/8 b - d3 0 1", 6, 1440467},
	// Castling where the rook's destination attacks the enemy king.
	{"Short castling gives check", "5k2/8/8/8/8/8/8/4K2R w K - 0 1", 6, 661072},
	{"Long castling gives check", "3k4/8/8/8/8/8/8/R3K3 w Q - 0 1", 6, 803711},
	// Castling rights being lost, and castling forbidden through attacked squares.
	{"Castle rights", "r3k2r/1b4bq/8/8/8/8/7B/R3K2R w KQkq - 0 1", 4, 1274206},
	{"Castling prevented", "r3k2r/8/3Q4/8/8/5q2/8/R3K2R b KQkq - 0 1", 4, 1720476},
	// Promotion as the only legal check evasion.
	{"Promote out of check", "2K2r2/4P3/8/8/8/8/8/3k4 w - - 0 1", 6, 3821001},
	// Moving a piece that was blocking an attack on its own king.
	{"Discovered check", "8/8/1P2K3/8/2n5/1q6/8/5k2 b - - 0 1", 5, 1004658},
	// Promotion delivering check, including the three under-promotions.
	{"Promote to give check", "4k3/1P6/8/8/8/8/K7/8 w - - 0 1", 6, 217342},
	{"Underpromote to give check", "8/P1k5/K7/8/8/8/8/8 w - - 0 1", 6, 92683},
	// Terminal nodes: positions with no legal moves must count as zero
	// continuations rather than being skipped or double counted.
	{"Self stalemate", "K1k5/8/P7/8/8/8/8/8 w - - 0 1", 6, 2217},
	{"Stalemate and checkmate #1", "8/k1P5/8/1K6/8/8/8/8 w - - 0 1", 7, 567584},
	{"Stalemate and checkmate #2", "8/8/2k5/5q2/5n2/8/5K2/8 b - - 0 1", 4, 23527},
}

func TestPerft(t *testing.T) {
	for _, tt := range perftCases {
		for depth, expected := range tt.expectedDepths {
			t.Run(fmt.Sprintf("%s depth %d", tt.name, depth), func(t *testing.T) {
				if depth > 4 && testing.Short() {
					t.Skipf("Skipping depth %d for %s in short mode", depth, tt.name)
				}
				pos, err := ParseFEN(tt.fen)
				if err != nil {
					t.Fatalf("Failed to create position from FEN: %v", err)
				}

				nodes := pos.Perft(depth)
				if nodes != expected {
					t.Errorf("Perft(%d) = %d; want %d", depth, nodes, expected)
				}
			})
		}
	}
}

func TestPerftEdges(t *testing.T) {
	for _, tt := range perftEdgeCases {
		t.Run(fmt.Sprintf("%s depth %d", tt.name, tt.depth), func(t *testing.T) {
			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("Failed to create position from FEN: %v", err)
			}

			nodes := pos.Perft(tt.depth)
			if nodes != tt.expectedNodes {
				t.Errorf("Perft(%d) = %d; want %d", tt.depth, nodes, tt.expectedNodes)
			}
		})
	}
}

func TestPerftDivide(t *testing.T) {
	depth := 3
	for _, tt := range perftCases {
		t.Run(fmt.Sprintf("%s depth %d", tt.name, depth), func(t *testing.T) {
			pos, err := ParseFEN(tt.fen)
			if err != nil {
				t.Fatalf("Failed to create position from FEN: %v", err)
			}

			nodes := pos.PerftDivide(depth)
			assertPerfCountsMatch(t, tt.expectedDepthThreeCounts, nodes)
		})
	}
}

func assertPerfCountsMatch(t *testing.T, expected, actual []PerftCount) {
	if len(expected) != len(actual) {
		t.Fatalf("Expected %d counts, got %d", len(expected), len(actual))
	}

	f := func(a, b PerftCount) int {
		if a.Move < b.Move {
			return -1
		}
		if a.Move > b.Move {
			return 1
		}
		return 0
	}

	slices.SortFunc(expected, f)
	slices.SortFunc(actual, f)

	if !slices.Equal(expected, actual) {
		t.Errorf("Expected counts %v, got %v", expected, actual)
	}
}

func BenchmarkPerft(b *testing.B) {
	pos, err := ParseFEN(FenKiwiPete)
	if err != nil {
		b.Fatalf("Failed to create position from FEN: %v", err)
	}
	b.ResetTimer()
	for b.Loop() {
		pos.Perft(5)
	}
}
