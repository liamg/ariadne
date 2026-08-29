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
	scratchBuffers   [][]board.Move // tried quiet moves
	ttSizeMB         int
	tt               *transpositionTable
	age              uint8
	progressCallback func(Progress)
	history          [16][64]int32
	evaluator        *eval.Evaluator
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
	scratchBuffers := make([][]board.Move, MaxPly)
	for i := range len(plyBuffers) {
		plyBuffers[i] = make([]board.Move, 0, MaxMoves)
		scoreBuffers[i] = make([]orderScore, 0, MaxMoves)
		scratchBuffers[i] = make([]board.Move, 0, MaxMoves)
	}

	s := &Searcher{
		state:          &State{},
		plyBuffers:     plyBuffers,
		scoreBuffers:   scoreBuffers,
		scratchBuffers: scratchBuffers,
		ttSizeMB:       64, // default to 64MB
		evaluator:      eval.New(),
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
	Infinite                 bool
	Ponder                   bool
}

type Result struct {
	BestMove  board.Move
	Score     eval.Score
	NodeCount uint64
	PV        []board.Move
}

// State is the per-search state - reset between searches
type State struct {
	NodeCount uint64
	NodeLimit uint64
	Stop      atomic.Bool
	maxPly    int
	killers   [MaxPly][2]board.Move
	pv        [MaxPly + 1][MaxPly + 1]board.Move
	pvLengths [MaxPly + 1]int
	rootDepth int
}

func (s *Searcher) Reset() {
	s.history = [16][64]int32{}
	s.tt.reset()
}

type Progress struct {
	Depth     int
	SelDepth  int
	Score     int64
	NodeCount int64
	ElapsedMS int64
	PV        []board.Move
}

// Search finds the best move for the given position, using the given search limits.
// It returns a Result containing the best move and its score.. The score is higher
// for better positions for the side to move.
func (s *Searcher) Search(ctx context.Context, pos *board.Position, limits Limits, ponderHit <-chan struct{}) Result {
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

	var result Result
	var score eval.Score
	var clockOffset atomic.Int64

	if limits.Ponder {
		clockOffset.Store(-1)
	} else {
		clockOffset.Store(0)
	}

	cancelChan := make(chan struct{})
	defer close(cancelChan)

	// grab a pointer to the _current_ state,
	// so we don't overwrite the state if the searcher is reused
	state := s.state
	go func() {
		if limits.Ponder {
			select {
			case <-ctx.Done():
				state.Stop.Store(true)
				return
			case <-cancelChan:
				return
			case <-ponderHit:
				// start the clock, ponder hit!
				clockOffset.Store(int64(time.Since(start)))
			}
		}

		if budgets.Unlimited {
			select {
			case <-ctx.Done():
				state.Stop.Store(true)
				return
			case <-cancelChan:
				return
			}
		} else {
			timer := time.NewTimer(time.Duration(budgets.Hard) * time.Millisecond)
			defer timer.Stop()

			select {
			case <-ctx.Done():
				state.Stop.Store(true)
				return
			case <-cancelChan:
				return
			case <-timer.C:
				state.Stop.Store(true)
				return
			}
		}
	}()

	// adoptRootPV copies the root PV out of the search state into result. Only
	// valid when pvLengths[0] > 0, i.e. at least one root move raised alpha.
	adoptRootPV := func(score eval.Score) {
		result.Score = score
		result.BestMove = s.state.pv[0][0]
		pv := s.state.pv[0][:s.state.pvLengths[0]]
		result.PV = make([]board.Move, len(pv))
		copy(result.PV, pv)
	}

	// report always sends what is in result, never a freshly computed score, so
	// the PV a GUI last saw can never disagree with the bestmove that follows it.
	// fastchess warns on exactly that mismatch, and a warning invalidates a run.
	report := func(depth int) {
		if s.progressCallback == nil {
			return
		}
		s.progressCallback(Progress{
			Depth:     depth,
			SelDepth:  s.state.maxPly,
			Score:     int64(result.Score),
			NodeCount: int64(s.state.NodeCount),
			ElapsedMS: time.Since(start).Milliseconds(),
			PV:        result.PV,
		})
	}

	for depth := 1; depth <= limits.Depth; depth++ {

		localClockOffset := clockOffset.Load()

		if depth > 1 && !budgets.Unlimited && localClockOffset != -1 && (time.Since(start)-time.Duration(localClockOffset)).Milliseconds() > budgets.Soft {
			break
		}

		s.state.rootDepth = depth

		// calculate best score
		score = s.negamax(pos, depth, 0, -eval.Infinity, eval.Infinity, true)
		if s.state.Stop.Load() {
			// The iteration was cut short, so score is only a lower bound - it
			// covers however many root moves were searched before the stop, and
			// pv[0][0] is merely the first of those. Root alpha starts at
			// -Infinity, so the very first root move searched always writes the
			// PV, whether or not it is any good.
			//
			// Taking it unconditionally lets a move nothing has been compared
			// against displace one that was fully searched a ply shallower, so it
			// is only adopted when it actually beat the last completed depth. The
			// NullMove case is depth 1 being cut short, where a partial answer
			// still beats having no move to play.
			if score != -eval.Infinity && s.state.pvLengths[0] > 0 &&
				(result.BestMove == board.NullMove || score > result.Score) {
				adoptRootPV(score)
				report(depth)
			}
			break
		}

		result.Score = score
		if s.state.pvLengths[0] > 0 {
			adoptRootPV(score)
		}

		report(depth)
	}

	if limits.Ponder {
		select {
		case <-ponderHit:
		case <-ctx.Done():
		}
	}

	s.age++

	result.NodeCount = s.state.NodeCount
	return result
}
