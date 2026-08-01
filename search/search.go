package search

import (
	"context"
	"sync/atomic"

	"github.com/liamg/chess/board"
	"github.com/liamg/chess/eval"
)

type Searcher struct {
	state      *State
	plyBuffers [][]board.Move
	ttSizeMB   int
	tt         *transpositionTable
	age        uint8
}

const MaxPly = 128

const MaxMoves = 256

type Option func(*Searcher)

func WithTranspositionTableSizeInMegaBytes(size int) Option {
	return func(s *Searcher) {
		s.ttSizeMB = size
	}
}

func New(options ...Option) *Searcher {
	plyBuffers := make([][]board.Move, MaxPly)
	for i := range len(plyBuffers) {
		plyBuffers[i] = make([]board.Move, 0, MaxMoves)
	}

	s := &Searcher{
		state:      &State{},
		plyBuffers: plyBuffers,
		ttSizeMB:   64, // default to 64MB
	}

	for _, opt := range options {
		opt(s)
	}

	s.tt = newTranspositionTable(s.ttSizeMB)

	return s
}

type Limits struct {
	Depth int8
}

type Result struct {
	BestMove  board.Move
	Score     eval.Score
	NodeCount uint64
}

// State is the per-search state - reset between searches
type State struct {
	NodeCount uint64
	Stop      atomic.Bool
	BestMove  board.Move
}

func (s *Searcher) Reset() {
	s.tt.reset()
}

// Search finds the best move for the given position, using the given search limits.
// It returns a Result containing the best move and its score.. The score is higher
// for better positions for the side to move.
func (s *Searcher) Search(ctx context.Context, pos *board.Position, limits Limits) Result {
	// reset state but leave TTs in place
	s.state = &State{}

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

	for depth := int8(1); depth <= limits.Depth; depth++ {
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
	}

	s.age++

	result.NodeCount = s.state.NodeCount
	return result
}
