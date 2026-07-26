package board

import (
	"testing"
)

func TestKingAttackCount(t *testing.T) {
	var count int
	for square := A1; square <= H8; square++ {
		count += kingAttacks[square].Count()
	}
	if count != 420 {
		t.Errorf("Expected 420 king attacks, got %d", count)
	}
}

func TestKnightAttackCount(t *testing.T) {
	var count int
	for square := A1; square <= H8; square++ {
		count += knightAttacks[square].Count()
	}
	if count != 336 {
		t.Errorf("Expected 336 knight attacks, got %d", count)
	}
}

func TestKingAttackSymmetry(t *testing.T) {
	for square := A1; square <= H8; square++ {
		attacks := kingAttacks[square]
		for attack, remaining := attacks.PopSquare(); attack != NoSquare; attack, remaining = remaining.PopSquare() {
			if !kingAttacks[attack].Has(square) {
				t.Errorf("King attack symmetry failed: %s does not attack %s", attack, square)
			}
		}
	}
}

func TestKnightAttackSymmetry(t *testing.T) {
	for square := A1; square <= H8; square++ {
		attacks := knightAttacks[square]
		for attack, remaining := attacks.PopSquare(); attack != NoSquare; attack, remaining = remaining.PopSquare() {
			if !knightAttacks[attack].Has(square) {
				t.Errorf("Knight attack symmetry failed: %s does not attack %s", attack, square)
			}
		}
	}
}

func TestKingAttacks(t *testing.T) {
	for sq := A1; sq <= H8; sq++ {
		expected := EmptyBitboard
		file, rank := sq.File(), sq.Rank()
		if a, err := SquareFromFileAndRank(file-1, rank-1); err == nil {
			expected = expected.Set(a)
		}
		if a, err := SquareFromFileAndRank(file, rank-1); err == nil {
			expected = expected.Set(a)
		}
		if a, err := SquareFromFileAndRank(file+1, rank-1); err == nil {
			expected = expected.Set(a)
		}
		if a, err := SquareFromFileAndRank(file-1, rank); err == nil {
			expected = expected.Set(a)
		}
		if a, err := SquareFromFileAndRank(file+1, rank); err == nil {
			expected = expected.Set(a)
		}
		if a, err := SquareFromFileAndRank(file-1, rank+1); err == nil {
			expected = expected.Set(a)
		}
		if a, err := SquareFromFileAndRank(file, rank+1); err == nil {
			expected = expected.Set(a)
		}
		if a, err := SquareFromFileAndRank(file+1, rank+1); err == nil {
			expected = expected.Set(a)
		}
		if kingAttacks[sq] != expected {
			t.Errorf("King attacks for %s do not match expected. Got: %s, Expected: %s", sq, kingAttacks[sq], expected)
		}
	}
}

func TestKnightAttacks(t *testing.T) {
	for sq := A1; sq <= H8; sq++ {
		expected := EmptyBitboard
		file, rank := sq.File(), sq.Rank()

		if a, err := SquareFromFileAndRank(file+1, rank+2); err == nil {
			expected = expected.Set(a)
		}
		if a, err := SquareFromFileAndRank(file+1, rank-2); err == nil {
			expected = expected.Set(a)
		}
		if a, err := SquareFromFileAndRank(file-1, rank+2); err == nil {
			expected = expected.Set(a)
		}
		if a, err := SquareFromFileAndRank(file-1, rank-2); err == nil {
			expected = expected.Set(a)
		}
		if a, err := SquareFromFileAndRank(file+2, rank+1); err == nil {
			expected = expected.Set(a)
		}
		if a, err := SquareFromFileAndRank(file+2, rank-1); err == nil {
			expected = expected.Set(a)
		}
		if a, err := SquareFromFileAndRank(file-2, rank+1); err == nil {
			expected = expected.Set(a)
		}
		if a, err := SquareFromFileAndRank(file-2, rank-1); err == nil {
			expected = expected.Set(a)
		}

		if knightAttacks[sq] != expected {
			t.Errorf("Knight attacks for %s do not match expected. Got: %s, Expected: %s", sq, knightAttacks[sq], expected)
		}
	}
}
