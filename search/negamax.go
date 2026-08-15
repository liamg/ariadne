package search

import (
	"github.com/liamg/ariadne/board"
	"github.com/liamg/ariadne/eval"
)

// negamax returns a score from the perspective of the moving player - so a higher score means a better outcome for that player
func (s *Searcher) negamax(pos *board.Position, depth int, ply int, alpha, beta eval.Score, canNull bool) eval.Score {
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

	// null move pruning
	if canNull && ply > 0 && depth >= 3 && beta <= eval.Mate-eval.Score(MaxPly) && !pos.InCheck(pos.SideToMove()) && pos.HasNonPawnMaterial(pos.SideToMove()) {
		// only null move if we have non-pawn material, otherwise it can lead to zugzwang
		undo := pos.MakeMove(board.NullMove)
		score := -s.negamax(pos, depth-3, ply+1, -beta, -beta+1, false)
		pos.UnmakeMove(undo)
		if score >= beta {
			// don't return mate scores as they're not "real" mate
			if score >= eval.Mate-MaxPly {
				return beta
			}
			return score
		}
	}

	picker := newMovePicker(
		moveGenTypeAllMoves,
		s.plyBuffers[ply][:0],
		s.scoreBuffers[ply][:0],
		pos,
		s.state.killers[ply],
		&s.history,
	)
	if ttMove != board.NullMove {
		picker.setTTMove(ttMove)
	}

	bestScore := -eval.Infinity
	var bestMove board.Move
	alphaOrig := alpha

	// reset scratch buffers for this ply
	s.scratchBuffers[ply] = s.scratchBuffers[ply][:0]

	var legalMoves int
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
		legalMoves++
		// NOTE: consider not decrementing depth when the enemy is in check here, as it's a promising branch...

		// PVS: for non-first move, search with a null window first, then re-search if it fails
		// for first move, use the full window straight away
		var score eval.Score
		if legalMoves == 1 {
			score = -s.negamax(pos, depth-1, ply+1, -beta, -alpha, true)
		} else {
			score = -s.negamax(pos, depth-1, ply+1, -alpha-1, -alpha, true)
			if score > alpha && score < beta {
				score = -s.negamax(pos, depth-1, ply+1, -beta, -alpha, true)
			}
		}

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
				// store killer moves for quiet moves, double pawn pushes, and castling
				if move.IsQuietish() {
					if move != s.state.killers[ply][0] {
						s.state.killers[ply][1] = s.state.killers[ply][0]
						s.state.killers[ply][0] = move
					}
					bonus := min(int32(depth*depth), scoreHistoryBonusCap)
					piece := pos.PieceAt(move.From())
					cur := s.history[piece][move.To()]
					s.history[piece][move.To()] = cur + bonus - (cur*bonus)/scoreHistoryMax
					for _, m := range s.scratchBuffers[ply] {
						piece := pos.PieceAt(m.From())
						cur = s.history[piece][m.To()]
						s.history[piece][m.To()] = cur - bonus - (cur*bonus)/scoreHistoryMax
					}
				}
				s.tt.store(pos.ZobristHash(), int16(bestScore), bestMove, depth, lowerBound, s.age, ply)
				return bestScore
			}
			if bestScore > alpha {
				alpha = bestScore
			}
		}
		if move.IsQuietish() {
			s.scratchBuffers[ply] = append(s.scratchBuffers[ply], move)
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
