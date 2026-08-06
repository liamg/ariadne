package uci

import (
	"fmt"
	"testing"
	"time"
)

func TestSetOption(t *testing.T) {
	h := newTestHarness(t)
	defer h.Close()
	u := New(h, h)
	go func() { _ = u.Run(t.Context()) }()

	h.Send("uci")
	h.WaitForExact("uciok", time.Second)

	expected := u.eng.options.transpositionTableMemoryMB * 2

	h.Send(fmt.Sprintf("setoption name Hash value %d", expected))
	h.WaitUntilReady(time.Second)

	if u.eng.options.transpositionTableMemoryMB != expected {
		t.Errorf("expected Hash option to be set to %d, got %d", expected, u.eng.options.transpositionTableMemoryMB)
	}
}
