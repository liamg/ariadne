package uci

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/liamg/ariadne/board"
	"github.com/liamg/ariadne/eval"
	"github.com/liamg/ariadne/search"
)

type engine struct {
	mu            sync.Mutex
	options       engineOptions
	searcher      *search.Searcher
	position      *board.Position
	ctxCancel     func()
	searching     bool
	ponderHitChan chan struct{}
}

func newEngine(r responder) *engine {
	e := &engine{
		options:  defaultEngineOptions,
		position: board.StartingPosition(),
	}
	e.searcher = e.createSearcher(r)
	return e
}

func (e *engine) createSearcher(r responder) *search.Searcher {
	return search.New(
		search.WithTranspositionTableSizeInMegaBytes(e.options.transpositionTableMemoryMB),
		search.WithProgressCallback(e.reportProgress(r)),
	)
}

func (e *engine) reportProgress(r responder) func(progress search.Progress) {
	var scoreStr string
	return func(progress search.Progress) {
		if int(progress.Score) > int(eval.Mate)-search.MaxPly {
			plies := int(eval.Mate) - int(progress.Score)
			scoreStr = fmt.Sprintf("mate %d", (plies+1)/2)
		} else if int(progress.Score) < -(int(eval.Mate) - search.MaxPly) {
			plies := (int(eval.Mate) + int(progress.Score))
			scoreStr = fmt.Sprintf("mate -%d", (plies+1)/2)
		} else {
			scoreStr = fmt.Sprintf("cp %d", progress.Score)
		}

		var pv string
		if len(progress.PV) > 0 {
			pvMoves := make([]string, len(progress.PV))
			for i, move := range progress.PV {
				pvMoves[i] = move.String()
			}
			pv = fmt.Sprintf(" pv %s", strings.Join(pvMoves, " "))
		}

		r.sendf("info depth %d seldepth %d score %s nodes %d nps %d time %d%s",
			progress.Depth,
			progress.SelDepth,
			scoreStr,
			progress.NodeCount,
			1000*progress.NodeCount/max(1, progress.ElapsedMS),
			progress.ElapsedMS,
			pv,
		)
	}
}

func (e *engine) clearTranspositionTable() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.searcher == nil {
		return
	}
	e.searcher.Reset()
}

func (e *engine) newGame() {
	e.clearTranspositionTable()
	e.mu.Lock()
	defer e.mu.Unlock()
	e.position = board.StartingPosition()
}

func (e *engine) setPosition(pos *board.Position) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.position = pos
}

var errAlreadySearching = fmt.Errorf("already searching")

func (e *engine) goSearch(ctx context.Context, limits search.Limits) (search.Result, error) {
	e.mu.Lock()
	if e.searching {
		e.mu.Unlock()
		return search.Result{}, errAlreadySearching
	}
	ctx, cancel := context.WithCancel(ctx)
	e.ctxCancel = cancel
	if e.searcher == nil || e.position == nil {
		e.mu.Unlock()
		return search.Result{}, fmt.Errorf("searcher or position is nil")
	}
	e.searching = true
	e.ponderHitChan = make(chan struct{})
	phc := e.ponderHitChan
	pos := e.position
	searcher := e.searcher
	e.mu.Unlock()
	result := searcher.Search(ctx, pos, limits, phc)
	e.stopSearch()
	e.mu.Lock()
	e.searching = false
	e.ponderHitChan = nil
	e.mu.Unlock()
	return result, nil
}

func (e *engine) stopSearch() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.ctxCancel != nil {
		e.ctxCancel()
		e.ctxCancel = nil
	}
}

func (e *engine) ponderHit() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.ponderHitChan != nil {
		close(e.ponderHitChan)
		e.ponderHitChan = nil
	}
}
