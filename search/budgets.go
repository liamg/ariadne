package search

import "github.com/liamg/ariadne/board"

type Budgets struct {
	Soft      int64
	Hard      int64
	Unlimited bool
}

func deriveBudgets(limits Limits, sideToMove board.Colour) Budgets {
	b := Budgets{}

	if limits.MoveTimeMS > 0 {
		b.Soft = max(limits.MoveTimeMS-limits.MoveOverheadMS, 1)
		b.Hard = b.Soft
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

	base := (remaining / n) + (increment * 4 / 5)
	base = min(remaining*2/5, base)

	budget := max(base-limits.MoveOverheadMS, 1)

	b.Hard = budget
	b.Soft = budget / 2

	return b
}
