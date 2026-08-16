package search

import "github.com/liamg/ariadne/board"

type orderScore int32 // deliberately different type to eval.Score

const (
	scoreQueenPromoBonus orderScore = 2_000_000 // add mvv/lva to this
	scoreCaptureBonus    orderScore = 1_000_000 // add mvv/lva to this
	scoreKiller1         orderScore = 900_000
	scoreKiller2         orderScore = 800_000
	scoreCounterMove     orderScore = 700_000
	scoreUnderpromotion  orderScore = -500_000
	scoreLosingCapture   orderScore = -1_000_000 // plus mvv/lva
	scoreKnightPromotion orderScore = 50_000     // sometimes uniquely useful - above quiet but below captures
	scoreHistoryMax                 = 16384
	scoreHistoryBonusCap            = 2000
	scoreTT              orderScore = 1_000_000
	scoreEqualCapture    orderScore = 600_000
)

func (mp *movePicker) scoreMove(move board.Move) orderScore {
	pos := mp.pos
	switch move.Kind() {
	case board.QuietMove, board.DoublePawnPush, board.KingsideCastle, board.QueensideCastle:

		if move == mp.killers[0] {
			return scoreKiller1
		}

		if move == mp.killers[1] {
			return scoreKiller2
		}

		return orderScore(mp.history[pos.PieceAt(move.From())][move.To()])
	case board.KnightPromotion:
		return scoreKnightPromotion
	case board.KnightPromotionCapture:
		victim := pos.PieceAt(move.To()).Type()
		return scoreKnightPromotion + mvvLva(victim, board.Pawn)
	case board.QueenPromotion:
		return scoreQueenPromoBonus
	case board.QueenPromotionCapture:
		victim := pos.PieceAt(move.To()).Type()
		return scoreQueenPromoBonus + mvvLva(victim, board.Pawn)
	case board.EnPassantCapture:
		// en passant is always pawn takes pawn, so both sides of the mvv/lva are known
		return captureScore(see(pos, move), board.Pawn, board.Pawn)
	case board.Capture:
		attacker := pos.PieceAt(move.From()).Type()
		victim := pos.PieceAt(move.To()).Type()
		return captureScore(see(pos, move), victim, attacker)
	case board.BishopPromotionCapture, board.RookPromotionCapture:
		attacker := pos.PieceAt(move.From()).Type()
		victim := pos.PieceAt(move.To()).Type()
		return scoreUnderpromotion + mvvLva(victim, attacker)
	case board.BishopPromotion, board.RookPromotion:
		return scoreUnderpromotion
	default:
		panic("unknown move kind")
	}
}

// captureScore buckets a capture by its static exchange evaluation, using mvv/lva
// only to order moves within a bucket.
func captureScore(seeScore int32, victim, attacker board.PieceType) orderScore {
	switch {
	case seeScore > 0:
		return scoreCaptureBonus + mvvLva(victim, attacker)
	case seeScore == 0:
		return scoreEqualCapture + mvvLva(victim, attacker)
	default:
		return scoreLosingCapture + mvvLva(victim, attacker)
	}
}

func mvvLva(victim, attacker board.PieceType) orderScore {
	return orderScore(int(victim)*10 - int(attacker))
}
