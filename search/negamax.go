package search

import (
	"github.com/liamg/ariadne/board"
	"github.com/liamg/ariadne/eval"
)

// negamax returns a score from the perspective of the moving player - so a higher score means a better outcome for that player
func (s *Searcher) negamax(pos *board.Position, depth int, ply int, alpha, beta eval.Score) eval.Score {
	if ply > 0 && pos.HalfMoveClock() >= 100 {
		s.state.NodeCount++
		return eval.Draw
	}

	// once we've hit the depth limit, evaluate the position and stop
	if depth <= 0 || ply >= MaxPly {
		return s.quiescence(pos, ply, alpha, beta)
	}

	s.state.NodeCount++

	if ply > 0 {
		if s.state.NodeLimit > 0 && s.state.NodeCount >= s.state.NodeLimit {
			s.state.Stop.Store(true)
			return 0
		}

		// never abort early for the root node
		if s.state.Stop.Load() {
			return 0
		}

		if pos.IsDrawByRepetition() {
			return eval.Draw
		}
	}

	var ttMove board.Move

	if entry, ok := s.tt.probe(pos.ZobristHash(), ply); ok {
		// only cutoff if we're deep enough and non-root - otherwise we have no move
		if ply > 0 && int(entry.Depth) >= depth {

			score := eval.Score(entry.Score)

			// TODO: cutoff needs to be aware of 50-move counter and repetition stalemate

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
		ttMove = entry.Move
	}

	picker := newMovePicker(moveGenTypeAllMoves, s.plyBuffers[ply][:0], s.scoreBuffers[ply][:0], pos)
	if ttMove != board.NullMove {
		picker.setTTMove(ttMove)
	}

	bestScore := -eval.Infinity
	var bestMove board.Move
	alphaOrig := alpha

	var legalMoves int
	for {
		move, ok := picker.next()
		if !ok {
			break
		}
		undo := pos.MakeMove(move)
		if !pos.IsLastMoveIllegal() {
			legalMoves++
			// NOTE: consider not decrementing depth when the enemy is in check here, as it's a promising branch...
			score := -s.negamax(pos, depth-1, ply+1, -beta, -alpha)
			if depth > 1 && s.state.Stop.Load() {
				pos.UnmakeMove(undo)
				if ply == 0 { // if this is the root, just take the best we have
					return bestScore
				}
				return 0
			}
			if score > bestScore {
				bestScore = score
				bestMove = move
				if ply == 0 {
					s.state.BestMove = move
				}
				if bestScore >= beta {
					pos.UnmakeMove(undo)
					s.tt.store(pos.ZobristHash(), int16(bestScore), bestMove, depth, lowerBound, s.age, ply)
					return bestScore
				}
				if bestScore > alpha {
					alpha = bestScore
				}
			}
		}
		pos.UnmakeMove(undo)
	}

	if legalMoves == 0 {
		if pos.InCheck(pos.SideToMove()) {
			// always negative if the mover is checkmated
			// but we make the score less "bad" depending on the number of plies taken to get there
			// i.e. losing mate in 1 is worse than losing mate in 10
			return -(eval.Mate - eval.Score(ply))
		}
		return eval.Draw
	}

	// lower bound is already handled in the loop above
	bound := upperBound
	if bestScore > alphaOrig {
		bound = exact
	}
	s.tt.store(pos.ZobristHash(), int16(bestScore), bestMove, depth, bound, s.age, ply)

	return bestScore
}
