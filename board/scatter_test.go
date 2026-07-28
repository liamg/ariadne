package board

import (
	"math/bits"
	"testing"
)

func TestScatterSamples(t *testing.T) {
	tests := []struct {
		name     string
		index    int
		mask     Bitboard
		expected Bitboard
	}{
		{
			name:     "scatter 0b101 into 0b101010000",
			index:    0b101,
			mask:     0b101010000,
			expected: 0b100010000,
		},
		{
			name:     "scatter 0b111 into 0b111000",
			index:    0b111,
			mask:     0b111000,
			expected: 0b111000,
		},
		{
			name:     "scatter 0b110 into 0b101010",
			index:    0b110,
			mask:     0b101010,
			expected: 0b101000,
		},
		{
			name:     "index 0",
			index:    0b0,
			mask:     0b1111,
			expected: 0b0,
		},
		{
			name:     "empty mask",
			index:    0b1111,
			mask:     0b0,
			expected: 0b0,
		},
		{
			name:     "bits above mask",
			index:    0b1111111,
			mask:     0b0111,
			expected: 0b0111,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := scatter(tt.index, tt.mask)
			if actual != tt.expected {
				t.Errorf("scatter(%b, %b) = %b; want %b", tt.index, tt.mask, actual, tt.expected)
			}
		})
	}
}

func TestScatterEnumeration(t *testing.T) {
	pieces := []struct {
		name  string
		masks [64]Bitboard
	}{
		{
			"rook",
			rookRelevantOccupancy,
		},
		{
			"bishop",
			bishopRelevantOccupancy,
		},
	}

	for _, piece := range pieces {
		t.Run(piece.name, func(t *testing.T) {
			for sq := A1; sq <= H8; sq++ {
				mask := piece.masks[sq]
				numBits := mask.Count()
				numCombinations := 1 << numBits

				results := make(map[Bitboard]struct{})
				for index := range numCombinations {
					scattered := scatter(index, mask)
					if scattered.Count() != bits.OnesCount(uint(index)) {
						t.Errorf("scatter(%b, %b) = %b; expected popcount %d, got %d", index, mask, scattered, bits.OnesCount(uint(index)), scattered.Count())
					}
					if scattered&^mask != 0 {
						t.Errorf("scatter(%b, %b) = %b; scattered bits outside of mask", index, mask, scattered)
					}
					results[scattered] = struct{}{}
				}
				if len(results) != numCombinations {
					t.Errorf("For square %v, expected %d unique scattered results, got %d", sq, numCombinations, len(results))
				}
			}
		})
	}
}
