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
	moves = moves[:0]
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
				if !p.IsSquareAttacked(E1, Black) && !p.IsSquareAttacked(F1, Black) && !p.IsSquareAttacked(G1, Black) {
					moves = append(moves, NewMove(E1, G1, KingsideCastle))
				}
			}
		}
		if p.state.castlingRights&CastlingWhiteQueenside != 0 {
			if occ&whiteQueensideEmptyMask == 0 {
				if !p.IsSquareAttacked(E1, Black) && !p.IsSquareAttacked(D1, Black) && !p.IsSquareAttacked(C1, Black) {
					moves = append(moves, NewMove(E1, C1, QueensideCastle))
				}
			}
		}
	case Black:
		if p.state.castlingRights&CastlingBlackKingside != 0 {
			if occ&blackKingsideEmptyMask == 0 {
				if !p.IsSquareAttacked(E8, White) && !p.IsSquareAttacked(F8, White) && !p.IsSquareAttacked(G8, White) {
					moves = append(moves, NewMove(E8, G8, KingsideCastle))
				}
			}
		}
		if p.state.castlingRights&CastlingBlackQueenside != 0 {
			if occ&blackQueensideEmptyMask == 0 {
				if !p.IsSquareAttacked(E8, White) && !p.IsSquareAttacked(D8, White) && !p.IsSquareAttacked(C8, White) {
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

	epSquare := p.state.enPassantSquare
	attackers := findEnPassantAttackers(p.sideToMove, pawns, epSquare)
	for attackers != 0 {
		from, attackers = attackers.PopSquare()
		moves = append(moves, NewMove(from, epSquare, EnPassantCapture))
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
		moveset := BishopLookup(sq, occ)
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
		moveset := RookLookup(sq, occ)
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
		if !p.IsLastMoveIllegal() {
			legalMoves = append(legalMoves, move)
		}
		p.UnmakeMove(undo)
	}
	return legalMoves
}

// GeneratePseudoLegalCapturesAndPromotions fills the passed buffer (growing if needed) with all pseudo-legal CAPTURES for the position
// this does NOT check for checks, so the caller must filter out moves that leave the king in check
func (p *Position) GeneratePseudoLegalCapturesAndPromotions(moves []Move) []Move {
	moves = moves[:0]
	occ := p.Occupancy()
	moves = p.generateKingCaptures(moves)
	moves = p.generateKnightCaptures(moves)
	moves = p.generatePawnCapturesAndPromotions(moves, occ)
	return p.generateSliderCaptures(moves, occ)
}

func (p *Position) generateKingCaptures(moves []Move) []Move {
	king := p.kingSquare[p.sideToMove]

	var sq Square
	kingAttackBB := kingAttacks[king] & p.byColour[p.sideToMove.Opposite()]
	for kingAttackBB != 0 {
		sq, kingAttackBB = kingAttackBB.PopSquare()
		moves = append(moves, NewMove(king, sq, Capture))
	}

	return moves
}

func (p *Position) generateKnightCaptures(moves []Move) []Move {
	knights := p.byType[Knight] & p.byColour[p.sideToMove]

	var from Square
	for knights != 0 {
		from, knights = knights.PopSquare()
		var sq Square
		knightAttackBB := knightAttacks[from] & p.byColour[p.sideToMove.Opposite()]
		for knightAttackBB != 0 {
			sq, knightAttackBB = knightAttackBB.PopSquare()
			moves = append(moves, NewMove(from, sq, Capture))
		}
	}

	return moves
}

func (p *Position) generatePawnCapturesAndPromotions(moves []Move, occ Bitboard) []Move {
	pawns := p.byType[Pawn] & p.byColour[p.sideToMove]

	var sq Square
	var from Square
	var pushOffset Square
	var westOffset, eastOffset Square
	var promoCandidates Bitboard
	var promoMask Bitboard
	switch p.sideToMove {
	case White:
		promoMask = Rank8Mask
		promoCandidates = pawns & Rank7Mask
		pushOffset = -8
		westOffset = -7
		eastOffset = -9
	case Black:
		promoMask = Rank1Mask
		promoCandidates = pawns & Rank2Mask
		pushOffset = 8
		westOffset = 9
		eastOffset = 7
	}

	singlesPromo, _ := findPawnPushes(p.sideToMove, promoCandidates, occ)

	for singlesPromo != 0 {
		sq, singlesPromo = singlesPromo.PopSquare()
		from = sq + pushOffset
		moves = append(moves, NewMove(from, sq, QueenPromotion))
		moves = append(moves, NewMove(from, sq, RookPromotion))
		moves = append(moves, NewMove(from, sq, BishopPromotion))
		moves = append(moves, NewMove(from, sq, KnightPromotion))
	}

	enemyColour := p.sideToMove.Opposite()
	enemy := p.byColour[enemyColour]

	epSquare := p.state.enPassantSquare
	attackers := findEnPassantAttackers(p.sideToMove, pawns, epSquare)
	for attackers != 0 {
		from, attackers = attackers.PopSquare()
		moves = append(moves, NewMove(from, epSquare, EnPassantCapture))
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

func (p *Position) generateSliderCaptures(moves []Move, occ Bitboard) []Move {
	queens := p.byType[Queen]
	bishopsAndQueens := (queens | p.byType[Bishop]) & p.byColour[p.sideToMove]
	rooksAndQueens := (queens | p.byType[Rook]) & p.byColour[p.sideToMove]

	enemyOcc := p.byColour[p.sideToMove.Opposite()]

	var sq Square
	var to Square
	for bishopsAndQueens != 0 {
		sq, bishopsAndQueens = bishopsAndQueens.PopSquare()
		moveset := BishopLookup(sq, occ)
		attacks := moveset & enemyOcc
		for attacks != 0 {
			to, attacks = attacks.PopSquare()
			moves = append(moves, NewMove(sq, to, Capture))
		}
	}

	for rooksAndQueens != 0 {
		sq, rooksAndQueens = rooksAndQueens.PopSquare()
		moveset := RookLookup(sq, occ)
		attacks := moveset & enemyOcc
		for attacks != 0 {
			to, attacks = attacks.PopSquare()
			moves = append(moves, NewMove(sq, to, Capture))
		}
	}

	return moves
}

// IsPseudoLegalMove checks if a move is pseudo-legal in the current position.
// It assumes the move was legal for the position it was generated in, e.g. cannot be a double pawn push from A1
func (p *Position) IsPseudoLegalMove(move Move) bool {
	if move == NullMove {
		return false
	}

	piece := p.mailbox[move.From()]

	// check the moving piece is the right colour
	if piece.Colour() != p.sideToMove {
		return false
	}

	// piece cannot move on top of an ally
	if p.byColour[p.sideToMove]&move.To().Bitboard() != 0 {
		return false
	}

	switch piece.Type() {
	case NoPieceType:
		return false
	case Pawn:
		switch move.Kind() {
		case DoublePawnPush:
			// cannot move on top of enemy - we checked for ally above
			if p.byColour[p.sideToMove.Opposite()]&move.To().Bitboard() != 0 {
				return false
			}
			if p.sideToMove == White {
				if p.Occupancy()&move.From().Bitboard().North() != 0 {
					// cannot double move through piece
					return false
				}
			} else if p.Occupancy()&move.From().Bitboard().South() != 0 {
				// cannot double move through piece
				return false
			}
			return true
		case QuietMove, KnightPromotion, BishopPromotion, RookPromotion, QueenPromotion:
			// cannot move on top of enemy - we checked for ally above
			if p.byColour[p.sideToMove.Opposite()]&move.To().Bitboard() != 0 {
				return false
			}
			return true
		case Capture, KnightPromotionCapture, BishopPromotionCapture, RookPromotionCapture, QueenPromotionCapture:
			// if the pawn cannot attack a piece at one of it's attack locations, then it is not a valid capture
			if pawnAttacks[p.sideToMove][move.From()]&move.To().Bitboard()&p.byColour[p.sideToMove.Opposite()] == 0 {
				return false
			}
			return true
		case EnPassantCapture:
			if pawnAttacks[p.sideToMove][move.From()]&move.To().Bitboard()&p.state.enPassantSquare.Bitboard() == 0 {
				return false
			}
			return true
		default:
			return false
		}
	case King:
		// quiet moves and captures fall through to the shared geometry check below
		switch move.Kind() {
		case KingsideCastle:
			if p.sideToMove == White {
				if move.To() != G1 {
					return false
				}
				if p.state.castlingRights&CastlingWhiteKingside == 0 {
					return false
				}
				if p.Occupancy()&whiteKingsideEmptyMask != 0 {
					return false
				}
				if p.IsSquareAttacked(E1, Black) || p.IsSquareAttacked(F1, Black) || p.IsSquareAttacked(G1, Black) {
					return false
				}
				return true
			}
			if move.To() != G8 {
				return false
			}
			if p.state.castlingRights&CastlingBlackKingside == 0 {
				return false
			}
			if p.Occupancy()&blackKingsideEmptyMask != 0 {
				return false
			}
			if p.IsSquareAttacked(E8, White) || p.IsSquareAttacked(F8, White) || p.IsSquareAttacked(G8, White) {
				return false
			}
			return true
		case QueensideCastle:
			if p.sideToMove == White {
				if move.To() != C1 {
					return false
				}
				if p.state.castlingRights&CastlingWhiteQueenside == 0 {
					return false
				}
				if p.Occupancy()&whiteQueensideEmptyMask != 0 {
					return false
				}
				if p.IsSquareAttacked(E1, Black) || p.IsSquareAttacked(D1, Black) || p.IsSquareAttacked(C1, Black) {
					return false
				}
				return true
			}
			if move.To() != C8 {
				return false
			}
			if p.state.castlingRights&CastlingBlackQueenside == 0 {
				return false
			}
			if p.Occupancy()&blackQueensideEmptyMask != 0 {
				return false
			}
			if p.IsSquareAttacked(E8, White) || p.IsSquareAttacked(D8, White) || p.IsSquareAttacked(C8, White) {
				return false
			}
			return true
		}
	case Knight, Bishop, Rook, Queen:
		// fall through to the shared geometry check below
	default:
		return false
	}

	// knights, sliders and non-castling king moves differ only in which attack
	// set they consult, which Attacks already knows
	switch move.Kind() {
	case QuietMove:
		if p.byColour[p.sideToMove.Opposite()]&move.To().Bitboard() != 0 {
			return false
		}
	case Capture:
		if p.byColour[p.sideToMove.Opposite()]&move.To().Bitboard() == 0 {
			return false
		}
	default:
		return false
	}

	return p.AttacksWithCustomOccupancy(piece.Type(), move.From(), p.Occupancy(), 0)&move.To().Bitboard() != 0
}
