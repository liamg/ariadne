package search

import (
	"github.com/liamg/chess/board"
	"github.com/liamg/chess/eval"
)

func (s *Searcher) minimax(pos *board.Position, depth int8, ply int) eval.Score {
	if depth <= 0 {
		score := s.quiescence(pos, ply, -eval.Infinity, eval.Infinity)
		if pos.SideToMove() == board.Black {
			score = -score
		}
		return score
	}

	s.state.NodeCount++

	var moves []board.Move
	if len(s.plyBuffers) <= ply {
		moves = make([]board.Move, 0, MaxMoves)
		s.plyBuffers = append(s.plyBuffers, moves)
	} else {
		moves = s.plyBuffers[ply][:0]
	}

	moves = pos.GeneratePseudoLegalMoves(moves)

	best := -eval.Infinity
	if pos.SideToMove() == board.Black {
		best = eval.Infinity
	}

	isWhite := pos.SideToMove() == board.White

	var legalMoves int
	for _, move := range moves {
		undo := pos.MakeMove(move)
		if !pos.IsLastMoveIllegal() {
			legalMoves++
			score := s.minimax(pos, depth-1, ply+1)
			if isWhite {
				if score > best {
					best = score
				}
			} else if score < best {
				best = score
			}
		}
		pos.UnmakeMove(undo)
	}

	if legalMoves == 0 {
		if pos.InCheck(pos.SideToMove()) {
			if isWhite {
				return -(eval.Mate - eval.Score(ply))
			}
			return eval.Mate - eval.Score(ply)
		}
		return eval.Draw
	}

	return best
}
