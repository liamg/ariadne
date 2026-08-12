package search

import (
	"github.com/liamg/ariadne/board"
	"github.com/liamg/ariadne/eval"
)

var (
	noKillers = [2]board.Move{board.NullMove, board.NullMove}
	noHistory = &[16][64]int32{}
)

func (s *Searcher) quiescence(pos *board.Position, ply int, alpha, beta eval.Score) eval.Score {
	s.state.NodeCount++

	if ply >= MaxPly {
		return eval.Evaluate(pos)
	}

	if ply > s.state.maxPly {
		s.state.maxPly = ply
	}

	var quietScore eval.Score

	inCheck := pos.InCheck(pos.SideToMove())
	if inCheck {
		quietScore = -eval.Infinity
	} else {
		quietScore = eval.Evaluate(pos)
	}

	if quietScore >= beta {
		return quietScore
	}

	if quietScore > alpha {
		alpha = quietScore
	}

	genType := moveGenTypeCapturesAndPromotions
	if inCheck {
		// TODO: need a GenerateEvasions() function to generate only evasions when in check, otherwise we will generate illegal moves
		genType = moveGenTypeAllMoves // should be evasions specifically
	}

	picker := newMovePicker(genType, s.plyBuffers[ply][:0], s.scoreBuffers[ply][:0], pos, noKillers, noHistory)

	bestScore := quietScore

	var legalMoveCount int

	for {
		move, ok := picker.next()
		if !ok {
			break
		}
		undo := pos.MakeMove(move)
		if pos.IsLastMoveIllegal() {
			pos.UnmakeMove(undo)
			continue
		}
		score := -s.quiescence(pos, ply+1, -beta, -alpha)

		legalMoveCount++

		if score > bestScore {
			bestScore = score
			if bestScore >= beta {
				pos.UnmakeMove(undo)
				return bestScore
			}
			if bestScore > alpha {
				alpha = bestScore
			}
		}

		pos.UnmakeMove(undo)
	}

	// TODO: store in TT?

	if inCheck && legalMoveCount == 0 {
		return -eval.Mate + eval.Score(ply)
	}

	return bestScore
}
