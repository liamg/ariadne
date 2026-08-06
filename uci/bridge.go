package uci

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
)

func New(in io.Reader, out io.Writer) *Bridge {
	b := &Bridge{
		output: out,
		input:  in,
	}
	b.eng = newEngine(b)
	return b
}

func (h *Bridge) Run(ctx context.Context) error {
	scanner := bufio.NewScanner(h.input)
	for scanner.Scan() {
		line := scanner.Text()
		if err := h.handleLine(ctx, line); err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			h.debugf("error handling line %q: %v", line, err)
		}
	}
	return scanner.Err()
}

type Bridge struct {
	eng    *engine
	input  io.Reader
	output io.Writer
	mu     sync.Mutex
	dbg    bool
}

func (h *Bridge) send(msg string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	_, _ = h.output.Write([]byte(msg + "\n"))
}

func (h *Bridge) sendf(format string, args ...any) {
	h.send(fmt.Sprintf(format, args...))
}

func (h *Bridge) debugf(format string, args ...any) {
	h.debug(fmt.Sprintf(format, args...))
}

func (h *Bridge) setDebug(enabled bool) {
	h.dbg = enabled
}

func (h *Bridge) debug(msg string) {
	if !h.dbg {
		return
	}
	h.send("info string " + msg)
}

func (h *Bridge) handleLine(ctx context.Context, original string) error {
	var cmd string
	line := strings.TrimSpace(original)
	for {
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		cmd, line, _ = strings.Cut(line, " ")
		if f, ok := lookupCommand(cmd); ok {
			line = strings.TrimSpace(line)
			return f(ctx, h.eng, h, line)
		}
	}
	h.debugf("unknown command: %s", original)
	return nil
}

type responder interface {
	setDebug(enabled bool)
	debug(msg string)
	debugf(format string, args ...any)
	send(msg string)
	sendf(format string, args ...any)
}

type commandFunc func(ctx context.Context, e *engine, r responder, input string) error

var commandRegistry = map[string]commandFunc{}

func register(name string, f commandFunc) {
	if _, ok := commandRegistry[name]; ok {
		panic(fmt.Sprintf("command %s already registered", name))
	}
	commandRegistry[name] = f
}

func lookupCommand(name string) (commandFunc, bool) {
	f, ok := commandRegistry[name]
	return f, ok
}
