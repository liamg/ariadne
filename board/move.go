package board

import "strings"

type Move uint16

type MoveKind uint16

const (
	QuietMove MoveKind = iota
	DoublePawnPush
	KingsideCastle
	QueensideCastle
	Capture
	EnPassantCapture
	_
	_
	KnightPromotion
	BishopPromotion
	RookPromotion
	QueenPromotion
	KnightPromotionCapture
	BishopPromotionCapture
	RookPromotionCapture
	QueenPromotionCapture
)

const NullMove = Move(0)

func NewMove(from Square, to Square, kind MoveKind) Move {
	return Move(
		(uint16(from) << 10) |
			(uint16(to) << 4) |
			uint16(kind))
}

func (m Move) From() Square {
	return Square((m >> 10) & 0x3F)
}

func (m Move) To() Square {
	return Square((m >> 4) & 0x3F)
}

func (m Move) Kind() MoveKind {
	return MoveKind(m & 0xF)
}

func (m Move) String() string {
	if m.IsNull() {
		return "0000"
	}
	if m.IsPromotion() {
		return m.From().String() + m.To().String() + strings.ToLower(m.PromotionPieceType().String())
	}

	return m.From().String() + m.To().String()
}

func (m Move) IsPromotion() bool {
	return m&8 != 0
}

func (m Move) IsCapture() bool {
	return m&4 != 0
}

func (m Move) IsNull() bool {
	return m == NullMove
}

func (m Move) IsCastle() bool {
	k := m.Kind()
	return k == KingsideCastle || k == QueensideCastle
}

func (m Move) PromotionPieceType() PieceType {
	if !m.IsPromotion() {
		return NoPieceType
	}
	return PieceType(m.Kind()&3) + Knight
}
