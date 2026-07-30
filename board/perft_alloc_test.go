//go:build !race

package board

import "testing"

func TestPerftIsZeroAlloc(t *testing.T) {
	for _, test := range perftCases {
		t.Run(test.name, func(t *testing.T) {
			p, err := ParseFEN(test.fen)
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", test.fen, err)
			}
			allocs := testing.AllocsPerRun(2, func() {
				p.Perft(4)
			})
			if allocs != 0 {
				t.Errorf("Perft(%q, 4) = %f allocations; want 0", test.fen, allocs)
			}
		})
	}
}
