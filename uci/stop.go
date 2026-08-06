package uci

import "context"

func init() {
	register("stop", func(_ context.Context, e *engine, r responder, input string) error {
		e.stopSearch()
		r.debug("stop command received, stopping search")
		return nil
	})
}
