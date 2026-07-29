package board

const (
	whiteKingsideEmptyMask  = (1 << F1) | (1 << G1)
	blackKingsideEmptyMask  = (1 << F8) | (1 << G8)
	whiteQueensideEmptyMask = (1 << B1) | (1 << C1) | (1 << D1)
	blackQueensideEmptyMask = (1 << B8) | (1 << C8) | (1 << D8)
)

// GeneratePseudoLegalMoves fills the passed buffer (growing if needed) with all pseudo-legal moves for the position
// this does NOT check for checks, so the caller must filter out moves that leave the king in check
func (p *Position) GeneratePseudoLegalMoves(moves []Move) []Move {
	occ := p.Occupancy()
	moves = p.generateKingMoves(moves, occ)
	moves = p.generateKnightMoves(moves, occ)
	moves = p.generatePawnMoves(moves, occ)
	return p.generateSliderMoves(moves, occ)
}

func (p *Position) generateKingMoves(moves []Move, occ Bitboard) []Move {
	king := p.kingSquare[p.sideToMove]
	kingMoves := kingAttacks[king] &^ occ

	var sq Square
	for kingMoves != 0 {
		sq, kingMoves = kingMoves.PopSquare()
		moves = append(moves, NewMove(king, sq, QuietMove))
	}

	kingAttackBB := kingAttacks[king] & p.byColour[p.sideToMove.Opposite()]
	for kingAttackBB != 0 {
		sq, kingAttackBB = kingAttackBB.PopSquare()
		moves = append(moves, NewMove(king, sq, Capture))
	}

	switch p.sideToMove {
	case White:
		if p.state.castlingRights&CastlingWhiteKingside != 0 {
			if occ&whiteKingsideEmptyMask == 0 {
				if !p.isSquareAttacked(E1, Black) && !p.isSquareAttacked(F1, Black) && !p.isSquareAttacked(G1, Black) {
					moves = append(moves, NewMove(E1, G1, KingsideCastle))
				}
			}
		}
		if p.state.castlingRights&CastlingWhiteQueenside != 0 {
			if occ&whiteQueensideEmptyMask == 0 {
				if !p.isSquareAttacked(E1, Black) && !p.isSquareAttacked(D1, Black) && !p.isSquareAttacked(C1, Black) {
					moves = append(moves, NewMove(E1, C1, QueensideCastle))
				}
			}
		}
	case Black:
		if p.state.castlingRights&CastlingBlackKingside != 0 {
			if occ&blackKingsideEmptyMask == 0 {
				if !p.isSquareAttacked(E8, White) && !p.isSquareAttacked(F8, White) && !p.isSquareAttacked(G8, White) {
					moves = append(moves, NewMove(E8, G8, KingsideCastle))
				}
			}
		}
		if p.state.castlingRights&CastlingBlackQueenside != 0 {
			if occ&blackQueensideEmptyMask == 0 {
				if !p.isSquareAttacked(E8, White) && !p.isSquareAttacked(D8, White) && !p.isSquareAttacked(C8, White) {
					moves = append(moves, NewMove(E8, C8, QueensideCastle))
				}
			}
		}
	}

	return moves
}

func (p *Position) generateKnightMoves(moves []Move, occ Bitboard) []Move {
	knights := p.byType[Knight] & p.byColour[p.sideToMove]

	var from Square
	for knights != 0 {
		from, knights = knights.PopSquare()
		knightMoves := knightAttacks[from] &^ occ

		var sq Square
		for knightMoves != 0 {
			sq, knightMoves = knightMoves.PopSquare()
			moves = append(moves, NewMove(from, sq, QuietMove))
		}

		knightAttackBB := knightAttacks[from] & p.byColour[p.sideToMove.Opposite()]
		for knightAttackBB != 0 {
			sq, knightAttackBB = knightAttackBB.PopSquare()
			moves = append(moves, NewMove(from, sq, Capture))
		}
	}

	return moves
}

func (p *Position) generatePawnMoves(moves []Move, occ Bitboard) []Move {
	pawns := p.byType[Pawn] & p.byColour[p.sideToMove]
	singles, doubles := findPawnPushes(p.sideToMove, pawns, occ)

	var sq Square
	var from Square
	var promoMask Bitboard
	var pushOffset Square
	var doubleOffset Square
	var westOffset, eastOffset Square
	switch p.sideToMove {
	case White:
		promoMask = Rank8Mask
		pushOffset = -8
		doubleOffset = -16
		westOffset = -7
		eastOffset = -9
	case Black:
		promoMask = Rank1Mask
		pushOffset = 8
		doubleOffset = 16
		westOffset = 9
		eastOffset = 7
	}
	singlesPromo := singles & promoMask
	singles &^= singlesPromo
	for singles != 0 {
		sq, singles = singles.PopSquare()
		from = sq + pushOffset
		moves = append(moves, NewMove(from, sq, QuietMove))
	}
	for singlesPromo != 0 {
		sq, singlesPromo = singlesPromo.PopSquare()
		from = sq + pushOffset
		moves = append(moves, NewMove(from, sq, QueenPromotion))
		moves = append(moves, NewMove(from, sq, RookPromotion))
		moves = append(moves, NewMove(from, sq, BishopPromotion))
		moves = append(moves, NewMove(from, sq, KnightPromotion))
	}
	for doubles != 0 {
		sq, doubles = doubles.PopSquare()
		from = sq + doubleOffset
		moves = append(moves, NewMove(from, sq, DoublePawnPush))
	}

	enemyColour := p.sideToMove.Opposite()
	enemy := p.byColour[enemyColour]

	if epSquare := p.state.enPassantSquare; epSquare != NoSquare {
		attackers := findEnPassantAttackers(p.sideToMove, pawns, epSquare)
		for attackers != 0 {
			from, attackers = attackers.PopSquare()
			moves = append(moves, NewMove(from, epSquare, EnPassantCapture))
		}
	}

	west, east := findPawnCaptures(p.sideToMove, pawns, enemy)
	westPromo := west & promoMask
	west &^= westPromo
	eastPromo := east & promoMask
	east &^= eastPromo

	for west != 0 {
		sq, west = west.PopSquare()
		moves = append(moves, NewMove(sq+westOffset, sq, Capture))
	}
	for westPromo != 0 {
		sq, westPromo = westPromo.PopSquare()
		from = sq + westOffset
		moves = append(moves, NewMove(from, sq, QueenPromotionCapture))
		moves = append(moves, NewMove(from, sq, RookPromotionCapture))
		moves = append(moves, NewMove(from, sq, BishopPromotionCapture))
		moves = append(moves, NewMove(from, sq, KnightPromotionCapture))
	}

	for east != 0 {
		sq, east = east.PopSquare()
		moves = append(moves, NewMove(sq+eastOffset, sq, Capture))
	}
	for eastPromo != 0 {
		sq, eastPromo = eastPromo.PopSquare()
		from = sq + eastOffset
		moves = append(moves, NewMove(from, sq, QueenPromotionCapture))
		moves = append(moves, NewMove(from, sq, RookPromotionCapture))
		moves = append(moves, NewMove(from, sq, BishopPromotionCapture))
		moves = append(moves, NewMove(from, sq, KnightPromotionCapture))
	}

	return moves
}

func (p *Position) generateSliderMoves(moves []Move, occ Bitboard) []Move {
	queens := p.byType[Queen]
	bishopsAndQueens := (queens | p.byType[Bishop]) & p.byColour[p.sideToMove]
	rooksAndQueens := (queens | p.byType[Rook]) & p.byColour[p.sideToMove]

	enemyOcc := p.byColour[p.sideToMove.Opposite()]

	var sq Square
	var to Square
	for bishopsAndQueens != 0 {
		sq, bishopsAndQueens = bishopsAndQueens.PopSquare()
		moveset := bishopLookup(sq, occ)
		quiet := moveset &^ occ
		attacks := moveset & enemyOcc
		for quiet != 0 {
			to, quiet = quiet.PopSquare()
			moves = append(moves, NewMove(sq, to, QuietMove))
		}
		for attacks != 0 {
			to, attacks = attacks.PopSquare()
			moves = append(moves, NewMove(sq, to, Capture))
		}
	}

	for rooksAndQueens != 0 {
		sq, rooksAndQueens = rooksAndQueens.PopSquare()
		moveset := rookLookup(sq, occ)
		quiet := moveset &^ occ
		attacks := moveset & enemyOcc
		for quiet != 0 {
			to, quiet = quiet.PopSquare()
			moves = append(moves, NewMove(sq, to, QuietMove))
		}
		for attacks != 0 {
			to, attacks = attacks.PopSquare()
			moves = append(moves, NewMove(sq, to, Capture))
		}
	}

	return moves
}

// GenerateLegalMoves generates all legal moves for the current position,
// filtering out any pseudo-legal moves that would leave the king in check.
// It is slow, and should not be used in search, but is useful for testing,
// and as a reference implementation/oracle. May also be useful for a
// GUI for human play etc.
func (p *Position) GenerateLegalMoves() []Move {
	moves := make([]Move, 0, 256)
	moves = p.GeneratePseudoLegalMoves(moves)
	legalMoves := make([]Move, 0, len(moves))
	for _, move := range moves {
		undo := p.MakeMove(move)
		if !p.isLastMoveIllegal() {
			legalMoves = append(legalMoves, move)
		}
		p.UnmakeMove(undo)
	}
	return legalMoves
}
