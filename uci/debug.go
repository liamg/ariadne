package uci

import (
	"context"
)

func init() {
	register("debug", func(ctx context.Context, e *engine, r responder, input string) error {
		// engine details
		e.options.debug = input == "on"
		r.setDebug(e.options.debug)
		if e.options.debug {
			r.send("info string Debug mode enabled")
		}
		return nil
	})
}
