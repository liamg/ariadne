package search

import (
	"github.com/liamg/chess/board"
	"github.com/liamg/chess/eval"
)

// negamax returns a score from the perspective of the moving player - so a higher score means a better outcome for that player
func (s *Searcher) negamax(pos *board.Position, depth int8, ply int, alpha, beta eval.Score) eval.Score {
	s.state.NodeCount++

	// once we've hit the depth limit, evaluate the position and stop
	if depth <= 0 || ply >= MaxPly {
		// TODO: i guess we need quiescence here eventually
		return eval.Evaluate(pos)
	}

	// never abort early for the root node
	if ply > 0 && s.state.Stop.Load() {
		return 0
	}

	var ttMove board.Move

	if entry, ok := s.tt.probe(pos.ZobristHash(), ply); ok {
		// only cutoff if we're deep enough and non-root - otherwise we have no move
		if ply > 0 && entry.Depth >= depth {

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

	// grab the buffer for this ply
	moves := s.plyBuffers[ply][:0]

	moves = pos.GeneratePseudoLegalMoves(moves)

	if ttMove != board.NullMove {
		// move ordering - try the best move from the transposition table first
		for i, move := range moves {
			if move == ttMove {
				moves[0], moves[i] = moves[i], moves[0]
				break
			}
		}
	}

	bestScore := -eval.Infinity
	var bestMove board.Move
	alphaOrig := alpha

	var legalMoves int
	for _, move := range moves {
		undo := pos.MakeMove(move)
		if !pos.IsLastMoveIllegal() {
			legalMoves++
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
