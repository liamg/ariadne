package search

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/liamg/ariadne/board"
	"github.com/liamg/ariadne/eval"
)

type Searcher struct {
	state            *State
	plyBuffers       [][]board.Move
	scoreBuffers     [][]orderScore
	ttSizeMB         int
	tt               *transpositionTable
	age              uint8
	progressCallback func(Progress)
}

const MaxPly = 128

const MaxMoves = 256

type Option func(*Searcher)

func WithTranspositionTableSizeInMegaBytes(size int) Option {
	return func(s *Searcher) {
		s.ttSizeMB = size
	}
}

func WithProgressCallback(callback func(Progress)) Option {
	return func(s *Searcher) {
		s.progressCallback = callback
	}
}

func New(options ...Option) *Searcher {
	plyBuffers := make([][]board.Move, MaxPly)
	scoreBuffers := make([][]orderScore, MaxPly)
	for i := range len(plyBuffers) {
		plyBuffers[i] = make([]board.Move, 0, MaxMoves)
		scoreBuffers[i] = make([]orderScore, 0, MaxMoves)
	}

	s := &Searcher{
		state:        &State{},
		plyBuffers:   plyBuffers,
		scoreBuffers: scoreBuffers,
		ttSizeMB:     64, // default to 64MB
	}

	for _, opt := range options {
		opt(s)
	}

	s.tt = newTranspositionTable(s.ttSizeMB)

	return s
}

type Limits struct {
	Depth                    int
	Nodes                    int64
	WhiteTimeMS              int64
	BlackTimeMS              int64
	MoveTimeMS               int64
	WhiteIncrementMS         int64
	BlackIncrementMS         int64
	MovesBeforeControlSwitch int64
	MoveOverheadMS           int64
	HasTimeControl           bool
}

type Result struct {
	BestMove  board.Move
	Score     eval.Score
	NodeCount uint64
}

// State is the per-search state - reset between searches
type State struct {
	NodeCount uint64
	NodeLimit uint64
	Stop      atomic.Bool
	BestMove  board.Move
	maxPly    int
}

func (s *Searcher) Reset() {
	s.tt.reset()
}

type Progress struct {
	Depth     int
	SelDepth  int
	Score     int64
	NodeCount int64
	ElapsedMS int64
}

// Search finds the best move for the given position, using the given search limits.
// It returns a Result containing the best move and its score.. The score is higher
// for better positions for the side to move.
func (s *Searcher) Search(ctx context.Context, pos *board.Position, limits Limits) Result {
	start := time.Now()

	// reset state but leave TTs in place
	s.state = &State{
		NodeLimit: uint64(limits.Nodes),
	}

	limits.Depth = min(limits.Depth, MaxPly)
	if limits.Depth <= 0 {
		limits.Depth = MaxPly
	}

	budgets := deriveBudgets(limits, pos.SideToMove())

	if !budgets.Unlimited {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(budgets.Hard)*time.Millisecond)
		defer cancel()
	}

	var result Result
	var score eval.Score

	cancelChan := make(chan struct{})
	defer close(cancelChan)

	go func() {
		// grab a pointer to the _current_ state,
		// so we don't overwrite the state if the searcher is reused
		state := s.state
		select {
		case <-ctx.Done():
			state.Stop.Store(true)
		case <-cancelChan:
		}
	}()

	for depth := 1; depth <= limits.Depth; depth++ {

		if depth > 1 && !budgets.Unlimited && time.Since(start) > time.Duration(budgets.Soft)*time.Millisecond {
			break
		}

		// calculate best score
		score = s.negamax(pos, depth, 0, -eval.Infinity, eval.Infinity)
		if s.state.Stop.Load() {
			if score != -eval.Infinity {
				result.Score = score
				result.BestMove = s.state.BestMove
			}
			break
		}
		result.Score = score
		result.BestMove = s.state.BestMove

		if s.progressCallback != nil {
			s.progressCallback(Progress{
				Depth:     depth,
				SelDepth:  s.state.maxPly,
				Score:     int64(score),
				NodeCount: int64(s.state.NodeCount),
				ElapsedMS: time.Since(start).Milliseconds(),
			})
		}
	}

	s.age++

	result.NodeCount = s.state.NodeCount
	return result
}
