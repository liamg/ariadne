package bench

import (
	"context"
	"fmt"
	"time"

	"github.com/liamg/ariadne/board"
	"github.com/liamg/ariadne/search"
)

// fens is the frozen bench position set. Adding, removing or reordering entries
// invalidates every historical bench number, so treat it as immutable.
var fens = []string{
	// openings
	"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",          // start position
	"r1bqkbnr/pppp1ppp/2n5/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R b KQkq - 3 3", // ruy lopez
	"rnbqkb1r/1p2pppp/p2p1n2/8/3NP3/2N5/PPP2PPP/R1BQKB1R w KQkq - 0 6",  // sicilian najdorf
	"rnbqkb1r/p3pppp/1p6/2ppP3/3N4/2P5/PP3PPP/R1BQKBNR w KQkq - 0 1",    // BK.04

	// middlegames
	"r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1",     // kiwipete
	"r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1",         // perft 4
	"rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8",                // perft 5
	"r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10", // perft 6
	"1k1r4/pp1b1R2/3q2pp/4p3/2B5/4Q3/PPP2B2/2K5 b - - 0 1",                     // BK.01
	"3r1k2/4npp1/1ppr3p/p6P/P2PPPP1/1NR5/5K2/2R5 w - - 0 1",                    // BK.02
	"2q1rr1k/3bbnnp/p2p1pp1/2pPp3/PpP1P1P1/1P2BNNP/2BQ1PRK/7R b - - 0 1",       // BK.03
	"r1b2rk1/2q1b1pp/p2ppn2/1p6/3QP3/1BN1B3/PPP3PP/R4RK1 w - - 0 1",            // BK.05
	"1nk1r1r1/pp2p1pp/2p2p2/3N1n2/2P1P3/1P4P1/P3PPBP/R2QK2R w KQ - 0 1",        // BK.07
	"r3r1k1/ppqb1ppp/8/4p1NQ/8/2P5/PP3PPP/R3R1K1 b - - 0 1",                    // BK.12

	// endgames
	"8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1",           // perft 3
	"8/k7/3p4/p2P1p2/P2P1P2/8/8/K7 w - - 0 1",             // lasker-reichhelm (fine #70), deep zugzwang
	"2r3k1/pppR1pp1/4p3/4P1P1/5P2/1P4K1/P1P5/8 w - - 0 1", // BK.06
	"4b3/p3kp2/6p1/3pP2p/2pP1P2/4K1P1/P3N2P/8 w - - 0 1",  // BK.08
	"1K6/1P1k4/8/8/8/8/r7/4R3 w - - 0 1",                  // lucena
	"8/8/1p1r1k2/p1pPN1p1/P3KnP1/1P6/8/3R4 w - - 0 1",     // rook and knight endgame
}

const (
	bDepth = 7
	ttMB   = 64
)

func Run(ctx context.Context) error {
	var total uint64

	realStart := time.Now()

	for _, fen := range fens {
		pos, err := board.ParseFEN(fen)
		if err != nil {
			return err
		}
		start := time.Now()
		nodes := benchPosition(ctx, pos, bDepth)
		ns := time.Since(start).Nanoseconds()
		fmt.Println("Position              :", fen)
		fmt.Println("Depth                 :", bDepth)
		fmt.Println("Nodes searched        :", nodes)
		fmt.Println("Nodes/second          :", 1000000000*nodes/ns)
		fmt.Println("Time taken (ms)       :", ns/1000000)
		fmt.Println()
		total += uint64(nodes)
	}

	fmt.Println("Total nodes searched  :", total)
	fmt.Println("Total time taken (ms) :", time.Since(realStart).Milliseconds())
	return nil
}

func benchPosition(ctx context.Context, pos *board.Position, depth int) int64 {
	s := search.New(
		search.WithTranspositionTableSizeInMegaBytes(ttMB),
	)
	result := s.Search(ctx, pos, search.Limits{Depth: depth}, nil)
	return int64(result.NodeCount)
}
