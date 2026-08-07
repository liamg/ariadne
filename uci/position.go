package uci

import (
	"context"
	"fmt"
	"strings"

	"github.com/liamg/ariadne/board"
)

func init() {
	register("position", func(_ context.Context, e *engine, r responder, input string) error {
		pos, remainder, _ := strings.Cut(input, " ")
		remainder = strings.TrimSpace(remainder)

		var position *board.Position
		var moveStr string

		switch pos {
		case "startpos":
			position = board.StartingPosition()
			_, moveStr, _ = strings.Cut(remainder, "moves ")
		case "fen":
			fen, moves, _ := strings.Cut(remainder, " moves ")
			fen = strings.TrimSpace(fen)
			moveStr = moves

			var err error
			position, err = board.ParseFEN(fen)
			if err != nil {
				return fmt.Errorf("invalid FEN: %s", fen)
			}

		default:
			return fmt.Errorf("invalid position command: %s", input)
		}

		moveStr = strings.TrimSpace(moveStr)

		for move := range strings.FieldsSeq(moveStr) {
			availableMoves := position.GenerateLegalMoves()
			var found bool
			for _, m := range availableMoves {
				if m.String() == move {
					position.MakeMove(m)
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("invalid move: %s", move)
			}
		}

		e.setPosition(position)
		return nil
	})
}
