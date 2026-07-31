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
	// TTs here
}

const MaxPly = 64

const MaxMoves = 256

func New() *Searcher {
	plyBuffers := make([][]board.Move, MaxPly)
	for i := range len(plyBuffers) {
		plyBuffers[i] = make([]board.Move, 0, MaxMoves)
	}

	return &Searcher{
		state:      &State{},
		plyBuffers: plyBuffers,
	}
}

type Limits struct {
	Depth int
}

type Result struct {
	BestMove board.Move
	Score    eval.Score
}

// State is the per-search state - reset between searches
type State struct {
	NodeCount int
	Stop      atomic.Bool
	BestMove  board.Move
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

	for depth := 1; depth <= limits.Depth; depth++ {
		// calculate best score
		score = s.negamax(pos, depth, 0, -eval.Infinity, eval.Infinity)
		if s.state.Stop.Load() {
			break
		}
		result.Score = score
		result.BestMove = s.state.BestMove
	}

	return result
}
