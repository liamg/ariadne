package eval

type Score int32

const (
	Mate     Score = 32000
	Draw     Score = 0
	Infinity Score = Mate + 1
)
