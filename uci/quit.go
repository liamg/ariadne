package uci

import (
	"context"
	"io"
)

func init() {
	register("quit", func(_ context.Context, e *engine, r responder, _ string) error {
		r.debug("quit command received, exiting")
		e.stopSearch()
		return io.EOF
	})
}
