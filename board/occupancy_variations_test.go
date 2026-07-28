package board

import (
	"slices"
	"testing"
)

func TestComputeOccupancyVariations(t *testing.T) {
	tests := []struct {
		name              string
		relevantOccupancy Bitboard
		expected          []Bitboard
	}{
		{
			name:              "3 bits contiguous",
			relevantOccupancy: 0b000111000,
			expected: []Bitboard{
				0b000000000,
				0b000001000,
				0b000010000,
				0b000011000,
				0b000100000,
				0b000101000,
				0b000110000,
				0b000111000,
			},
		},
		{
			name:              "3 bits non-contiguous",
			relevantOccupancy: 0b001010100,
			expected: []Bitboard{
				0b000000000,
				0b000000100,
				0b000010000,
				0b000010100,
				0b001000000,
				0b001000100,
				0b001010000,
				0b001010100,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := computeOccupancyVariations(test.relevantOccupancy)
			if !slices.Equal(actual, test.expected) {
				t.Errorf("computeOccupancyVariations(%b) = %v; want %v", test.relevantOccupancy, actual, test.expected)
			}
		})
	}
}
