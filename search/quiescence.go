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
		return s.evaluator.Evaluate(pos)
	}

	if ply > s.state.maxPly {
		s.state.maxPly = ply
	}

	// tt probe
	if entry, ok := s.tt.probe(pos.ZobristHash(), ply); ok {
		score := eval.Score(entry.Score)
		switch entry.Bound {
		case exact:
			return score
		case lowerBound:
			if score >= beta {
				return score
			}
		case upperBound:
			if score <= alpha {
				return score
			}
		}
	}

	var quietScore eval.Score

	inCheck := pos.InCheck(pos.SideToMove())
	if inCheck {
		quietScore = -eval.Infinity
	} else {
		quietScore = s.evaluator.Evaluate(pos)
	}

	alphaOrig := alpha

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
	var bestMove board.Move

	var legalMoveCount int

	for {
		move, os, ok := picker.next()
		if !ok {
			break
		}
		if !inCheck && os < scoreUnderpromotion { // band above losing capture as we add mvvLVA to losing capture score
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
			bestMove = move
			if bestScore >= beta {
				pos.UnmakeMove(undo)
				s.tt.store(pos.ZobristHash(), int16(bestScore), bestMove, 0, lowerBound, s.age, ply)
				return bestScore
			}
			if bestScore > alpha {
				alpha = bestScore
			}
		}

		pos.UnmakeMove(undo)
	}

	if inCheck && legalMoveCount == 0 {
		score := -eval.Mate + eval.Score(ply)
		s.tt.store(pos.ZobristHash(), int16(score), bestMove, 0, exact, s.age, ply)
		return score
	}

	if bestScore > alphaOrig {
		s.tt.store(pos.ZobristHash(), int16(bestScore), bestMove, 0, exact, s.age, ply)
	} else {
		s.tt.store(pos.ZobristHash(), int16(bestScore), bestMove, 0, upperBound, s.age, ply)
	}

	return bestScore
}
