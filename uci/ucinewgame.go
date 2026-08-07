package uci

import "context"

func init() {
	register("ucinewgame", func(_ context.Context, e *engine, r responder, input string) error {
		r.debug("ucinewgame command received, resetting engine state")
		e.newGame()
		return nil
	})
}
