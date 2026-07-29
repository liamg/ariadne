package board

import "testing"

func TestMoveRoundTrip(t *testing.T) {
	for sqA := A1; sqA <= H8; sqA++ {
		for sqB := A1; sqB <= H8; sqB++ {
			for kind := QuietMove; kind <= QueenPromotionCapture; kind++ {
				move := NewMove(sqA, sqB, kind)

				if move.From() != sqA {
					t.Fatalf("Expected From() to return %v, got %v", sqA, move.From())
				}
				if move.To() != sqB {
					t.Fatalf("Expected To() to return %v, got %v", sqB, move.To())
				}
				if move.Kind() != kind {
					t.Fatalf("Expected Kind() to return %v, got %v", kind, move.Kind())
				}
				if kind == QuietMove && sqA == A1 && sqB == A1 {
					if !move.IsNull() {
						t.Fatalf("Expected NullMove to be null, got %v", move)
					}
				} else if move.IsNull() {
					t.Fatalf("Expected IsNull() to return false, got true")
				}
			}
		}
	}
}

func TestMoveIsPromotion(t *testing.T) {
	tests := []struct {
		kind          MoveKind
		want          bool
		wantPieceType PieceType
	}{
		{QuietMove, false, NoPieceType},
		{DoublePawnPush, false, NoPieceType},
		{KingsideCastle, false, NoPieceType},
		{QueensideCastle, false, NoPieceType},
		{Capture, false, NoPieceType},
		{EnPassantCapture, false, NoPieceType},
		{KnightPromotion, true, Knight},
		{BishopPromotion, true, Bishop},
		{RookPromotion, true, Rook},
		{QueenPromotion, true, Queen},
		{KnightPromotionCapture, true, Knight},
		{BishopPromotionCapture, true, Bishop},
		{RookPromotionCapture, true, Rook},
		{QueenPromotionCapture, true, Queen},
	}

	for _, tt := range tests {
		move := NewMove(A1, A2, tt.kind)
		if got := move.IsPromotion(); got != tt.want {
			t.Errorf("MoveKind %v: IsPromotion() = %v, want %v", tt.kind, got, tt.want)
		}
		pt := move.PromotionPieceType()
		if pt != tt.wantPieceType {
			t.Errorf("MoveKind %v: PromotionPieceType() = %v, want %v", tt.kind, pt, tt.wantPieceType)
		}
	}
}

func TestMoveIsCapture(t *testing.T) {
	tests := []struct {
		kind MoveKind
		want bool
	}{
		{QuietMove, false},
		{DoublePawnPush, false},
		{KingsideCastle, false},
		{QueensideCastle, false},
		{Capture, true},
		{EnPassantCapture, true},
		{KnightPromotion, false},
		{BishopPromotion, false},
		{RookPromotion, false},
		{QueenPromotion, false},
		{KnightPromotionCapture, true},
		{BishopPromotionCapture, true},
		{RookPromotionCapture, true},
		{QueenPromotionCapture, true},
	}

	for _, tt := range tests {
		move := NewMove(A1, A2, tt.kind)
		if got := move.IsCapture(); got != tt.want {
			t.Errorf("MoveKind %v: IsCapture() = %v, want %v", tt.kind, got, tt.want)
		}
	}
}

func TestMoveIsCastle(t *testing.T) {
	tests := []struct {
		kind MoveKind
		want bool
	}{
		{QuietMove, false},
		{DoublePawnPush, false},
		{KingsideCastle, true},
		{QueensideCastle, true},
		{Capture, false},
		{EnPassantCapture, false},
		{KnightPromotion, false},
		{BishopPromotion, false},
		{RookPromotion, false},
		{QueenPromotion, false},
		{KnightPromotionCapture, false},
		{BishopPromotionCapture, false},
		{RookPromotionCapture, false},
		{QueenPromotionCapture, false},
	}

	for _, tt := range tests {
		move := NewMove(A1, A2, tt.kind)
		if got := move.IsCastle(); got != tt.want {
			t.Errorf("MoveKind %v: IsCastle() = %v, want %v", tt.kind, got, tt.want)
		}
	}
}

func TestNullMove(t *testing.T) {
	if !NullMove.IsNull() {
		t.Errorf("Expected NullMove.IsNull() to return true, got false")
	}
	if NullMove.From() != A1 || NullMove.To() != A1 || NullMove.Kind() != QuietMove {
		t.Errorf("Expected NullMove to have From=A1, To=A1, Kind=QuietMove, got From=%v, To=%v, Kind=%v", NullMove.From(), NullMove.To(), NullMove.Kind())
	}
}

func TestMoveString(t *testing.T) {
	tests := []struct {
		move Move
		want string
	}{
		{NewMove(A2, A4, DoublePawnPush), "a2a4"},
		{NewMove(E7, E8, QueenPromotion), "e7e8q"},
		{NewMove(G7, H8, KnightPromotionCapture), "g7h8n"},
		{NewMove(E1, G1, KingsideCastle), "e1g1"},
		{NewMove(E1, C1, QueensideCastle), "e1c1"},
		{NullMove, "0000"},
	}

	for _, tt := range tests {
		if got := tt.move.String(); got != tt.want {
			t.Errorf("Move %v: String() = %v, want %v", tt.move, got, tt.want)
		}
	}
}
