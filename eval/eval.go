package eval

import "github.com/liamg/chess/board"

func Evaluate(p *board.Position) Score {
	// dumb material counting for now

	var score Score

	score += Score(p.Pieces(board.White, board.Pawn).Count() * 100)
	score += Score(p.Pieces(board.White, board.Knight).Count() * 320)
	score += Score(p.Pieces(board.White, board.Bishop).Count() * 330)
	score += Score(p.Pieces(board.White, board.Rook).Count() * 500)
	score += Score(p.Pieces(board.White, board.Queen).Count() * 900)

	score -= Score(p.Pieces(board.Black, board.Pawn).Count() * 100)
	score -= Score(p.Pieces(board.Black, board.Knight).Count() * 320)
	score -= Score(p.Pieces(board.Black, board.Bishop).Count() * 330)
	score -= Score(p.Pieces(board.Black, board.Rook).Count() * 500)
	score -= Score(p.Pieces(board.Black, board.Queen).Count() * 900)

	if p.SideToMove() == board.Black {
		score = -score
	}
	return score
}
