package search

import "github.com/liamg/ariadne/board"

var seeValues = [7]int32{
	board.NoPieceType: 0,
	board.Pawn:        100,
	board.Knight:      320,
	board.Bishop:      330,
	board.Rook:        500,
	board.Queen:       900,
	board.King:        10000,
}

func see(pos *board.Position, move board.Move) int32 {
	from := move.From()
	to := move.To()

	var gain [32]int32 // 32 is the most pieces that can be on the board

	side := pos.SideToMove()

	occ := pos.Occupancy() &^ from.Bitboard() // remove the moving piece from the occupancy

	if move.Kind() == board.EnPassantCapture {
		gain[0] = seeValues[board.Pawn]
		// remove the capture pawn
		if side == board.White {
			occ &^= (to - 8).Bitboard()
		} else {
			occ &^= (to + 8).Bitboard()
		}
	} else {
		gain[0] = seeValues[pos.PieceAt(to).Type()]
	}

	attackers := pos.AllAttackersForSquare(to, occ)

	attacker := pos.PieceAt(from).Type()

	side = side.Opposite()

	var d int
	for d = 1; d < len(gain); d++ {
		gain[d] = seeValues[attacker] - gain[d-1]
		sideAttackers := attackers & pos.PiecesByColour(side)
		if sideAttackers == 0 {
			break
		}

		var typedAttackers board.Bitboard
		var attackerType board.PieceType
		for pt := board.Pawn; pt <= board.King; pt++ {
			typedAttackers = sideAttackers & pos.PiecesByType(pt)
			if typedAttackers != 0 {
				attackerType = pt
				break
			}
		}

		if attackerType == board.King && attackers&pos.PiecesByColour(side.Opposite()) != 0 {
			// king can't capture into attacked square - cannot move into check!
			break
		}

		sq, _ := typedAttackers.PopSquare()
		occ &^= sq.Bitboard()
		attackers &^= sq.Bitboard()

		attackers |= board.BishopLookup(to, occ) & (pos.PiecesByType(board.Bishop) | pos.PiecesByType(board.Queen))
		attackers |= board.RookLookup(to, occ) & (pos.PiecesByType(board.Rook) | pos.PiecesByType(board.Queen))
		attackers &= occ

		attacker = attackerType
		side = side.Opposite()
	}

	for i := d - 1; i > 0; i-- {
		gain[i-1] = -max(-gain[i-1], gain[i])
	}

	return gain[0]
}
