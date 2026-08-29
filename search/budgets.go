package search

import "github.com/liamg/ariadne/board"

type Budgets struct {
	// Soft is the amount of time that the search should aim to finish by. Fine to exceed.
	Soft int64
	// Hard is the maximum amount of time that the search should take. Exceeding this risks losing on time.
	Hard int64
	// Unlimited indicates that the search has no time limits.
	Unlimited      bool
	IterationStart int64
}

func deriveBudgets(limits Limits, sideToMove board.Colour) Budgets {
	b := Budgets{}

	if limits.Infinite {
		b.Unlimited = true
		return b
	}

	if limits.MoveTimeMS > 0 {
		b.Soft = max(limits.MoveTimeMS-limits.MoveOverheadMS, 1)
		b.Hard = b.Soft
		b.IterationStart = b.Soft
		return b
	}

	if !limits.HasTimeControl {
		b.Unlimited = true
		return b
	}

	var remaining, increment int64

	if sideToMove == board.White {
		remaining = max(limits.WhiteTimeMS, 0)
		increment = max(limits.WhiteIncrementMS, 0)
	} else {
		remaining = max(limits.BlackTimeMS, 0)
		increment = max(limits.BlackIncrementMS, 0)
	}

	n := limits.MovesBeforeControlSwitch
	if n <= 0 {
		n = 25
	}

	optimum := (remaining / n) + (increment * 4 / 5)

	b.Hard = max(min(optimum*3, remaining*2/5)-limits.MoveOverheadMS, 1)
	b.Soft = min(optimum, b.Hard)
	b.IterationStart = max(b.Soft/2, 1)

	return b
}
