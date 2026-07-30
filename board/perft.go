package board

// Perft returns the number of leaf nodes of all legal moves from
// the position, generated to the specified depth.
// It is used to test move generation.
func (p *Position) Perft(depth int) uint64 {
	if depth <= 0 {
		return 1
	}

	depth--

	var nodes uint64

	moves := make([]Move, 0, 256)
	moves = p.GeneratePseudoLegalMoves(moves)
	for _, move := range moves {
		undo := p.MakeMove(move)
		if !p.IsLastMoveIllegal() {
			nodes += p.Perft(depth)
		}
		p.UnmakeMove(undo)
	}

	return nodes
}

type PerftCount struct {
	Move  Move
	Count uint64
}

// PerftDivide returns a slice of PerftCount, which contains the number of leaf nodes
// for each legal move from the position, generated to the specified depth.
func (p *Position) PerftDivide(depth int) []PerftCount {
	if depth <= 0 {
		return nil
	}

	depth--

	nodes := make([]PerftCount, 0, 256)

	moves := make([]Move, 0, 256)
	moves = p.GeneratePseudoLegalMoves(moves)
	for _, move := range moves {
		undo := p.MakeMove(move)
		if !p.IsLastMoveIllegal() {
			nodes = append(nodes, PerftCount{
				Move:  move,
				Count: p.Perft(depth),
			})
		}
		p.UnmakeMove(undo)
	}

	return nodes
}
