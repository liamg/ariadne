package search

import "github.com/liamg/ariadne/board"

type orderScore int32 // deliberately different type to eval.Score

const (
	scoreQuiet           orderScore = 0         // TODO: history score +- 32768
	scoreQueenPromoBonus orderScore = 2_000_000 // add mvv/lva to this
	scoreCaptureBonus    orderScore = 1_000_000 // add mvv/lva to this
	scoreKiller1         orderScore = 900_000
	scoreKiller2         orderScore = 800_000
	scoreCounterMove     orderScore = 700_000
	scoreUnderpromotion  orderScore = -500_000
	scoreLosingCapture   orderScore = -1_000_000 // plus mvv/lva
	scoreKnightPromotion orderScore = 50_000     // sometimes uniquely useful - above quiet but below captures
)

func scoreMove(pos *board.Position, move board.Move) orderScore {
	switch move.Kind() {
	case board.QuietMove, board.DoublePawnPush, board.KingsideCastle, board.QueensideCastle:
		return scoreQuiet // TODO: history later
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
		if pos.SideToMove() == board.White {
			victim := pos.PieceAt(pos.EnPassantSquare() - 8).Type()
			return scoreCaptureBonus + mvvLva(victim, board.Pawn)
		}
		victim := pos.PieceAt(pos.EnPassantSquare() + 8).Type()
		return scoreCaptureBonus + mvvLva(victim, board.Pawn)
	case board.Capture:
		attacker := pos.PieceAt(move.From()).Type()
		victim := pos.PieceAt(move.To()).Type()
		return scoreCaptureBonus + mvvLva(victim, attacker)
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

func mvvLva(victim, attacker board.PieceType) orderScore {
	return orderScore(int(victim)*10 - int(attacker))
}
