package uci

import (
	"bufio"
	"bytes"
	"io"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"
)

type testHarness struct {
	t        *testing.T
	buf      bytes.Buffer
	received []string
	in       bytes.Buffer
	mu       *sync.Mutex
	cond     *sync.Cond
	closed   bool
}

func newTestHarness(t *testing.T) *testHarness {
	t.Helper()
	mu := &sync.Mutex{}
	cond := sync.NewCond(mu)
	return &testHarness{
		received: []string{},
		t:        t,
		mu:       mu,
		cond:     cond,
	}
}

func (t *testHarness) Close() {
	t.t.Helper()
	t.mu.Lock()
	defer t.mu.Unlock()
	t.closed = true
	t.cond.Broadcast()
}

func (t *testHarness) Send(input string) {
	t.t.Helper()
	t.mu.Lock()
	defer t.mu.Unlock()
	_, _ = t.in.WriteString(input + "\n")
	t.cond.Broadcast()
}

func (t *testHarness) receive(line string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.received = append(t.received, line)
}

func (t *testHarness) Read(p []byte) (n int, err error) {
	t.t.Helper()
	t.mu.Lock()
	defer t.mu.Unlock()
	for t.in.Len() == 0 && !t.closed {
		t.cond.Wait()
	}
	if t.closed && t.in.Len() == 0 {
		return 0, io.EOF
	}
	return t.in.Read(p)
}

func (t *testHarness) Write(p []byte) (n int, err error) {
	n, err = t.buf.Write(p)
	if err != nil {
		return n, err
	}

	s := bufio.NewScanner(bytes.NewReader(p))
	for s.Scan() {
		t.receive(s.Text())
	}

	return n, s.Err()
}

func (t *testHarness) ContainsExact(s string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	return slices.Contains(t.received, s)
}

func (t *testHarness) Contains(s string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, line := range t.received {
		if strings.Contains(line, s) {
			return true
		}
	}
	return false
}

func (t *testHarness) WaitUntilReady(d time.Duration) {
	t.t.Helper()

	index := -1
	t.mu.Lock()
	for i, msg := range t.received {
		if msg == "readyok" {
			index = i
		}
	}
	t.mu.Unlock()

	t.Send("isready")

	start := time.Now()
	for {

		t.mu.Lock()
		found := slices.Contains(t.received[index+1:], "readyok")
		t.mu.Unlock()

		if found {
			return
		}

		time.Sleep(10 * time.Millisecond)
		if time.Since(start) > d {
			t.t.Fatalf("timeout waiting for readyok")
		}
	}
}

func (t *testHarness) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.received = []string{}
}

func (t *testHarness) ContainsPrefix(s string) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, line := range t.received {
		if strings.HasPrefix(line, s) {
			return true
		}
	}
	return false
}

func (t *testHarness) WaitForExact(s string, timeout time.Duration) {
	t.t.Helper()
	start := time.Now()
	for {
		if t.ContainsExact(s) {
			return
		}
		time.Sleep(10 * time.Millisecond)
		if time.Since(start) > timeout {
			t.t.Fatalf("timeout waiting for exact string: %s", s)
		}
	}
}

func (t *testHarness) WaitForPrefix(s string, timeout time.Duration) {
	t.t.Helper()
	start := time.Now()
	for {
		if t.ContainsPrefix(s) {
			return
		}
		time.Sleep(10 * time.Millisecond)
		if time.Since(start) > timeout {
			t.t.Fatalf("timeout waiting for prefix: %s", s)
		}
	}
}

func (t *testHarness) Messages() []string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]string(nil), t.received...)
}

func (t *testHarness) Init() {
	t.t.Helper()
	t.Send("uci")
}
