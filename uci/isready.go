package uci

import "context"

func init() {
	register("isready", func(_ context.Context, _ *engine, r responder, input string) error {
		r.send("readyok")
		return nil
	})
}
