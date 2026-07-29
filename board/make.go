package board

import "fmt"

type Undo struct {
	previousState     PositionState
	move              Move
	capturedPieceType PieceType
	// maybe add zobrist here?
}

var castlingRightsMask = [64]CastlingRights{}

func init() {
	castlingRightsMask[E1] = CastlingWhiteKingside | CastlingWhiteQueenside
	castlingRightsMask[A1] = CastlingWhiteQueenside
	castlingRightsMask[H1] = CastlingWhiteKingside
	castlingRightsMask[E8] = CastlingBlackKingside | CastlingBlackQueenside
	castlingRightsMask[A8] = CastlingBlackQueenside
	castlingRightsMask[H8] = CastlingBlackKingside
}

// MakeMove applies the given move to the position and returns an Undo object that can be used to revert the move.
// The move MUST be valid, and should be one which has come from the move generator. If the move is invalid, the behavior is undefined.
func (p *Position) MakeMove(move Move) Undo {
	undo := Undo{
		previousState: p.state,
		move:          move,
	}

	to := move.To()
	from := move.From()

	p.state.castlingRights &^= castlingRightsMask[from] | castlingRightsMask[to]

	// always reset ep square
	p.state.enPassantSquare = NoSquare
	piece := p.mailbox[from]

	kind := move.Kind()

	switch kind {
	case QuietMove:
		if piece.Type() == Pawn {
			p.state.halfMoveClock = 0
		} else {
			p.state.halfMoveClock++
		}
		p.movePiece(piece, from, to)
	case Capture:
		p.state.halfMoveClock = 0
		target := p.mailbox[to]
		undo.capturedPieceType = target.Type()
		p.removeKnownPiece(to, target) // remove captured piece
		p.movePiece(piece, from, to)
	case EnPassantCapture:
		p.state.halfMoveClock = 0

		switch p.sideToMove {
		case White:
			capturedPawnSquare, _ := to.Bitboard().South().PopSquare()
			undo.capturedPieceType = Pawn
			target := p.mailbox[capturedPawnSquare]
			p.removeKnownPiece(capturedPawnSquare, target)
		case Black:
			capturedPawnSquare, _ := to.Bitboard().North().PopSquare()
			undo.capturedPieceType = Pawn
			target := p.mailbox[capturedPawnSquare]
			p.removeKnownPiece(capturedPawnSquare, target)
		}

		p.movePiece(piece, from, to)
	case DoublePawnPush:
		p.state.halfMoveClock = 0
		p.movePiece(piece, from, to)
		if p.byType[Pawn]&p.byColour[p.sideToMove.Opposite()]&(to.Bitboard().West()|to.Bitboard().East()) != 0 {
			switch p.sideToMove {
			case White:
				p.state.enPassantSquare, _ = to.Bitboard().South().PopSquare()
			case Black:
				p.state.enPassantSquare, _ = to.Bitboard().North().PopSquare()
			}
		}
	case BishopPromotionCapture, KnightPromotionCapture, RookPromotionCapture, QueenPromotionCapture:
		target := p.mailbox[to]
		undo.capturedPieceType = target.Type()
		p.removeKnownPiece(to, target) // remove captured piece
		fallthrough
	case BishopPromotion, KnightPromotion, RookPromotion, QueenPromotion:
		p.state.halfMoveClock = 0
		p.removeKnownPiece(from, piece)
		p.addPiece(to, NewPiece(p.sideToMove, move.PromotionPieceType()))
	case KingsideCastle:
		p.state.halfMoveClock++
		// the piece being moved is the KING, the ROOK moved by side-effect
		switch p.sideToMove {
		case White:
			p.movePiece(piece, E1, G1)
			p.movePiece(NewPiece(White, Rook), H1, F1)
		case Black:
			p.movePiece(piece, E8, G8)
			p.movePiece(NewPiece(Black, Rook), H8, F8)
		}
	case QueensideCastle:
		// the piece being moved is the KING, the ROOK moved by side-effect
		p.state.halfMoveClock++
		switch p.sideToMove {
		case White:
			p.movePiece(piece, E1, C1)
			p.movePiece(NewPiece(White, Rook), A1, D1)
		case Black:
			p.movePiece(piece, E8, C8)
			p.movePiece(NewPiece(Black, Rook), A8, D8)
		}
	default:
		panic(fmt.Sprintf("cannot make invalid move 0b%b", move))
	}

	if p.sideToMove == Black {
		p.fullMoveNumber++
		p.sideToMove = White
	} else {
		p.sideToMove = Black
	}

	return undo
}

// UnmakeMove reverts the last move made on the position using the provided Undo object.
// The Undo object must be the one returned by the corresponding MakeMove call. If the Undo object is invalid
// or does not match the last move made, the behavior is undefined.
func (p *Position) UnmakeMove(undo Undo) {
	p.state = undo.previousState

	movedSide := p.sideToMove.Opposite()
	if movedSide == Black {
		p.fullMoveNumber--
	}

	to := undo.move.To()
	from := undo.move.From()
	piece := p.mailbox[to]

	kind := undo.move.Kind()

	switch kind {
	case QuietMove, DoublePawnPush:
		p.movePiece(piece, to, from)
	case Capture:
		capturedPiece := NewPiece(p.sideToMove, undo.capturedPieceType)
		p.movePiece(piece, to, from)
		p.addPiece(to, capturedPiece)
	case EnPassantCapture:
		p.movePiece(piece, to, from)
		capturedPiece := NewPiece(p.sideToMove, undo.capturedPieceType)
		switch movedSide {
		case White:
			capturedPawnSquare, _ := to.Bitboard().South().PopSquare()
			p.addPiece(capturedPawnSquare, capturedPiece)
		case Black:
			capturedPawnSquare, _ := to.Bitboard().North().PopSquare()
			p.addPiece(capturedPawnSquare, capturedPiece)
		}
	case BishopPromotionCapture, KnightPromotionCapture, RookPromotionCapture, QueenPromotionCapture:
		capturedPiece := NewPiece(p.sideToMove, undo.capturedPieceType)
		p.removeKnownPiece(to, piece)
		p.addPiece(from, NewPiece(movedSide, Pawn))
		p.addPiece(to, capturedPiece)
	case BishopPromotion, KnightPromotion, RookPromotion, QueenPromotion:
		p.removeKnownPiece(to, piece)
		p.addPiece(from, NewPiece(movedSide, Pawn))
	case KingsideCastle:
		switch movedSide {
		case White:
			p.movePiece(piece, G1, E1)
			p.movePiece(NewPiece(White, Rook), F1, H1)
		case Black:
			p.movePiece(piece, G8, E8)
			p.movePiece(NewPiece(Black, Rook), F8, H8)
		}
	case QueensideCastle:
		switch movedSide {
		case White:
			p.movePiece(piece, C1, E1)
			p.movePiece(NewPiece(White, Rook), D1, A1)
		case Black:
			p.movePiece(piece, C8, E8)
			p.movePiece(NewPiece(Black, Rook), D8, A8)
		}
	default:
		panic(fmt.Sprintf("cannot undo invalid move 0b%b", undo.move))
	}

	p.sideToMove = movedSide
}
