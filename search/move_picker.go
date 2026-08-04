package search

import (
	"github.com/liamg/chess/board"
)

type movePicker struct {
	ttMove    board.Move
	ttUsed    bool
	genType   moveGenType
	generated bool
	moves     []board.Move
	scores    []orderScore
	pos       *board.Position
	index     int
}

type moveGenType uint8

const (
	moveGenTypeAllMoves moveGenType = iota
	moveGenTypeCapturesAndPromotions
	moveGenTypeEvasions
)

func newMovePicker(genType moveGenType, moves []board.Move, scores []orderScore, pos *board.Position) movePicker {
	return movePicker{
		moves:   moves,
		genType: genType,
		pos:     pos,
		scores:  scores,
	}
}

func (mp *movePicker) setTTMove(move board.Move) {
	if mp.ttUsed || mp.generated {
		panic("cannot set tt move after moves have been generated")
	}
	mp.ttMove = move
}

func (mp *movePicker) next() (board.Move, bool) {
	// if we have a tt move, return it, and remove it
	if mp.ttMove != board.NullMove && !mp.ttUsed {
		mp.ttUsed = true
		if mp.pos.IsPseudoLegalMove(mp.ttMove) {
			return mp.ttMove, true
		}
	}

	if !mp.generated {
		mp.generated = true
		switch mp.genType {
		case moveGenTypeAllMoves:
			mp.moves = mp.pos.GeneratePseudoLegalMoves(mp.moves[:0])
		case moveGenTypeCapturesAndPromotions:
			mp.moves = mp.pos.GeneratePseudoLegalCapturesAndPromotions(mp.moves[:0])
		default: // TODO: support evasion generation
			panic("unsupported move generation type")
		}

		// remove a previously used tt move from the list of moves
		if mp.ttUsed {
			for i, move := range mp.moves {
				if move == mp.ttMove {
					mp.moves[0], mp.moves[i] = mp.moves[i], mp.moves[0]
					mp.index++
					break
				}
			}
		}

		mp.scores = mp.scores[:len(mp.moves)]
		for i, move := range mp.moves[mp.index:] {
			s := scoreMove(mp.pos, move)
			mp.scores[i+mp.index] = s
		}
	}

	if len(mp.moves[mp.index:]) == 0 {
		return board.NullMove, false
	}

	var highScore orderScore
	var highIndex int

	for i, score := range mp.scores[mp.index:] {
		if i == 0 || score > highScore {
			highScore = score
			highIndex = i + mp.index
		}
	}

	mp.scores[mp.index], mp.scores[highIndex] = mp.scores[highIndex], mp.scores[mp.index]
	mp.moves[mp.index], mp.moves[highIndex] = mp.moves[highIndex], mp.moves[mp.index]

	move := mp.moves[mp.index]
	mp.index++
	return move, true
}
