package search

import (
	"github.com/liamg/chess/board"
	"github.com/liamg/chess/eval"
)

// negamax returns a score from the perspective of the moving player - so a higher score means a better outcome for that player
func (s *Searcher) negamax(pos *board.Position, depth, ply int, alpha, beta eval.Score) eval.Score {
	s.state.NodeCount++

	if depth <= 0 {
		return eval.Evaluate(pos)
	}

	var moves []board.Move
	if len(s.plyBuffers) <= ply {
		moves = make([]board.Move, 0, MaxMoves)
		s.plyBuffers = append(s.plyBuffers, moves)
	} else {
		moves = s.plyBuffers[ply][:0]
	}
	moves = pos.GeneratePseudoLegalMoves(moves)

	if ply == 0 && s.state.BestMove != board.NullMove {
		// move ordering - try the best move from the previous search first
		// this is the payoff of iterative deepening
		// we know from the previous iteration what the best move was,
		// so we try it first in the hope that it will be the best move again
		// and we can cut off the rest of the search
		for i, move := range moves {
			if move == s.state.BestMove {
				moves[0], moves[i] = moves[i], moves[0]
				break
			}
		}
	}

	best := -eval.Infinity

	var legalMoves int
	for _, move := range moves {
		undo := pos.MakeMove(move)
		if !pos.IsLastMoveIllegal() {
			legalMoves++
			score := -s.negamax(pos, depth-1, ply+1, -beta, -alpha)
			if score > best {
				best = score
				if ply == 0 {
					s.state.BestMove = move
				}
				if best >= beta {
					pos.UnmakeMove(undo)
					return best
				}
				if best > alpha {
					alpha = best
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

	return best
}
