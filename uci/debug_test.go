package uci

import (
	"testing"
	"time"
)

func TestDebug(t *testing.T) {
	h := newTestHarness(t)
	defer h.Close()

	u := New(h, h)
	go func() { _ = u.Run(t.Context()) }()

	h.Send("uci")
	h.WaitForExact("uciok", time.Second)

	h.Send("nonsense")
	h.WaitUntilReady(time.Second)
	if h.Contains("nonsense") {
		t.Errorf("unexpected info message when debug is off by default")
	}

	h.Send("debug on")
	h.Send("nonsense")
	h.WaitUntilReady(time.Second)
	if !h.Contains("nonsense") {
		t.Errorf("expected info message when debug is on")
	}
	h.Clear()

	h.Send("debug off")
	h.Send("nonsense")
	h.WaitUntilReady(time.Second)
	if h.Contains("nonsense") {
		t.Errorf("unexpected info message when debug is off")
	}
}
