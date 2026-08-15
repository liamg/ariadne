package uci

import (
	"context"
)

func init() {
	register("ponderhit", func(ctx context.Context, e *engine, r responder, input string) error {
		e.ponderHit()
		return nil
	})
}
