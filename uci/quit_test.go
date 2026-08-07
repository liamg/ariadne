package uci

import (
	"testing"
	"time"
)

func TestQuit(t *testing.T) {
	h := newTestHarness(t)
	defer h.Close()

	u := New(h, h)
	quitChan := make(chan struct{})
	var err error
	go func() {
		err = u.Run(t.Context())
		close(quitChan)
	}()

	h.Send("quit")
	select {
	case <-quitChan:
		// Quit successfully
		if err != nil {
			t.Errorf("engine returned an error on quit: %v", err)
		}
	case <-time.After(time.Second):
		t.Errorf("engine did not quit within the expected time")
	}
}
