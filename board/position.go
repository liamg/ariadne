package board

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

type Position struct {
	byType         [7]Bitboard // Bitboards for each piece type (Pawn, Knight, Bishop, Rook, Queen, King)
	byColour       [2]Bitboard // Bitboards for each color (White, Black)
	sideToMove     Colour
	fullMoveNumber int // starts at one
	mailbox        [64]Piece
	state          PositionState
	kingSquare     [2]Square
}

type PositionState struct {
	castlingRights  CastlingRights
	enPassantSquare Square
	halfMoveClock   uint8
	// NOTE: zobrist goes here eventually...
}

type CastlingRights uint8

const (
	CastlingWhiteKingside CastlingRights = 1 << iota
	CastlingWhiteQueenside
	CastlingBlackKingside
	CastlingBlackQueenside
	CastlingAll                 = CastlingWhiteKingside | CastlingWhiteQueenside | CastlingBlackKingside | CastlingBlackQueenside
	CastlingNone CastlingRights = 0
)

func EmptyPosition() *Position {
	return &Position{
		state: PositionState{
			castlingRights:  CastlingNone,
			enPassantSquare: NoSquare,
			halfMoveClock:   0,
		},
		sideToMove:     White,
		fullMoveNumber: 1,
		byType:         [7]Bitboard{},
		byColour:       [2]Bitboard{},
		mailbox:        [64]Piece{},
		kingSquare:     [2]Square{NoSquare, NoSquare},
	}
}

func StartingPosition() *Position {
	pos := EmptyPosition()
	pos.state.castlingRights = CastlingAll
	pos.addPiece(A1, WhiteRook)
	pos.addPiece(B1, WhiteKnight)
	pos.addPiece(C1, WhiteBishop)
	pos.addPiece(D1, WhiteQueen)
	pos.addPiece(E1, WhiteKing)
	pos.addPiece(F1, WhiteBishop)
	pos.addPiece(G1, WhiteKnight)
	pos.addPiece(H1, WhiteRook)

	for sq := A2; sq <= H2; sq++ {
		pos.addPiece(sq, WhitePawn)
	}

	for sq := A7; sq <= H7; sq++ {
		pos.addPiece(sq, BlackPawn)
	}

	pos.addPiece(A8, BlackRook)
	pos.addPiece(B8, BlackKnight)
	pos.addPiece(C8, BlackBishop)
	pos.addPiece(D8, BlackQueen)
	pos.addPiece(E8, BlackKing)
	pos.addPiece(F8, BlackBishop)
	pos.addPiece(G8, BlackKnight)
	pos.addPiece(H8, BlackRook)

	return pos
}

// RandomPosition generates a random chess position with a realistic distribution of pieces.
// It does NOT guarantee legality!
func RandomPosition(rnd *rand.Rand) *Position {
	if rnd == nil {
		rnd = rand.New(rand.NewPCG(1, 2))
	}

	randomPieces := func(count int, free Bitboard) (Bitboard, Bitboard) {
		if free == 0 {
			return EmptyBitboard, free
		}
		n := free.Count()
		if count > n {
			return free, EmptyBitboard
		}

		selected := EmptyBitboard
		for range count {
			n = free.Count()
			k := rnd.IntN(n) + 1
			var sq Square
			pool := free
			for range k {
				sq, pool = pool.PopSquare()
			}
			if sq != NoSquare {
				selected |= sq.Bitboard()
				free &= ^sq.Bitboard()
			}
		}

		return selected, free
	}

	pos := EmptyPosition()
	free := FullBitboard
	whitePawns, free := randomPieces(rnd.IntN(9), free&^(Rank1Mask|Rank8Mask)) // always realistic pawn density
	blackPawns, free := randomPieces(rnd.IntN(9), free&^(Rank1Mask|Rank8Mask))
	free |= (Rank1Mask | Rank8Mask) // allow pieces on first and last rank
	whiteKing, free := randomPieces(1, free)
	blackKing, free := randomPieces(1, free)
	whiteKnights, free := randomPieces(rnd.IntN(3), free)
	blackKnights, free := randomPieces(rnd.IntN(3), free)
	whiteBishops, free := randomPieces(rnd.IntN(3), free)
	blackBishops, free := randomPieces(rnd.IntN(3), free)
	whiteRooks, free := randomPieces(rnd.IntN(3), free)
	blackRooks, free := randomPieces(rnd.IntN(3), free)
	whiteQueens, free := randomPieces(1, free)
	blackQueens, _ := randomPieces(1, free)

	pos.byColour[White] = whitePawns | whiteKnights | whiteBishops | whiteRooks | whiteQueens | whiteKing
	pos.byColour[Black] = blackPawns | blackKnights | blackBishops | blackRooks | blackQueens | blackKing
	pos.byType[Pawn] = whitePawns | blackPawns
	pos.byType[Knight] = whiteKnights | blackKnights
	pos.byType[Bishop] = whiteBishops | blackBishops
	pos.byType[Rook] = whiteRooks | blackRooks
	pos.byType[Queen] = whiteQueens | blackQueens
	pos.byType[King] = whiteKing | blackKing

	whitePawns.EachSquareSlow(func(sq Square) {
		pos.mailbox[sq] = WhitePawn
	})
	blackPawns.EachSquareSlow(func(sq Square) {
		pos.mailbox[sq] = BlackPawn
	})
	whiteKnights.EachSquareSlow(func(sq Square) {
		pos.mailbox[sq] = WhiteKnight
	})
	blackKnights.EachSquareSlow(func(sq Square) {
		pos.mailbox[sq] = BlackKnight
	})
	whiteBishops.EachSquareSlow(func(sq Square) {
		pos.mailbox[sq] = WhiteBishop
	})
	blackBishops.EachSquareSlow(func(sq Square) {
		pos.mailbox[sq] = BlackBishop
	})
	whiteRooks.EachSquareSlow(func(sq Square) {
		pos.mailbox[sq] = WhiteRook
	})
	blackRooks.EachSquareSlow(func(sq Square) {
		pos.mailbox[sq] = BlackRook
	})
	whiteQueens.EachSquareSlow(func(sq Square) {
		pos.mailbox[sq] = WhiteQueen
	})
	blackQueens.EachSquareSlow(func(sq Square) {
		pos.mailbox[sq] = BlackQueen
	})
	whiteKing.EachSquareSlow(func(sq Square) {
		pos.mailbox[sq] = WhiteKing
		pos.kingSquare[White] = sq
	})
	blackKing.EachSquareSlow(func(sq Square) {
		pos.mailbox[sq] = BlackKing
		pos.kingSquare[Black] = sq
	})

	pos.sideToMove = Colour(rnd.IntN(2))
	pos.state.castlingRights = CastlingRights(rnd.IntN(16))
	pos.fullMoveNumber = rnd.IntN(100) + 1
	pos.state.halfMoveClock = uint8(min(rnd.IntN(100)+1, pos.fullMoveNumber*2-1)) // half-moves cannot exceed full moves * 2 - 1e

	if rnd.IntN(10) < 8 { // 80% chance of no en passant square
		pos.state.enPassantSquare = NoSquare
	} else {

		opp := pos.sideToMove.Opposite()
		enemyPawns := pos.byColour[opp] & pos.byType[Pawn]

		switch opp {
		case Black:
			doublePushedPawns := enemyPawns & Rank5Mask
			if doublePushedPawns > 0 {
				sq, _ := doublePushedPawns.PopSquare()
				pos.state.enPassantSquare, _ = sq.Bitboard().North().PopSquare()
			}
		case White:
			doublePushedPawns := enemyPawns & Rank4Mask
			if doublePushedPawns > 0 {
				sq, _ := doublePushedPawns.PopSquare()
				pos.state.enPassantSquare, _ = sq.Bitboard().South().PopSquare()
			}
		}
	}

	return pos
}

func (p *Position) addPiece(square Square, piece Piece) *Position {
	colour := piece.Colour()
	pt := piece.Type()
	p.byType[pt] |= square.Bitboard()
	p.byColour[colour] |= square.Bitboard()
	p.mailbox[square] = piece
	if pt == King {
		p.kingSquare[colour] = square
	}
	return p
}

func (p *Position) movePiece(piece Piece, from, to Square) {
	colour := piece.Colour()
	pt := piece.Type()
	p.byType[pt] ^= from.Bitboard() | to.Bitboard()
	p.byColour[colour] ^= from.Bitboard() | to.Bitboard()
	p.mailbox[from] = NoPiece
	p.mailbox[to] = piece
	if pt == King {
		p.kingSquare[colour] = to
	}
}

func (p *Position) removeKnownPiece(square Square, piece Piece) {
	p.byType[piece.Type()] &= ^square.Bitboard()
	p.byColour[piece.Colour()] &= ^square.Bitboard()
	p.mailbox[square] = NoPiece
}

type PositionError struct {
	Message string
	Square  Square
	Piece   Piece
}

func (e *PositionError) Error() string {
	parts := []string{fmt.Sprintf("corrupt position state: %s", e.Message)}
	parts = append(parts, fmt.Sprintf("square: %s", e.Square.String()))
	if e.Piece != NoPiece {
		parts = append(parts, fmt.Sprintf("piece: %s", e.Piece.String()))
	}
	return strings.Join(parts, ", ")
}

func (p *Position) validateStateSlow() error {
	for sq := A1; sq <= H8; sq++ {
		piece := p.mailbox[sq]
		if piece == NoPiece {
			for pieceType := Pawn; pieceType <= King; pieceType++ {
				if p.byType[pieceType].Has(sq) {
					return &PositionError{Message: "byType has piece not in mailbox", Square: sq}
				}
			}
			for colour := White; colour <= Black; colour++ {
				if p.byColour[colour].Has(sq) {
					return &PositionError{Message: "byColour has piece not in mailbox", Square: sq}
				}
			}
			continue
		}

		// Check that the piece is in the correct byType and byColour bitboards
		// Check byType
		if !p.byType[piece.Type()].Has(sq) {
			return &PositionError{Message: "mailbox has piece not present in byType", Square: sq, Piece: piece}
		}

		if !p.byColour[piece.Colour()].Has(sq) {
			return &PositionError{Message: "mailbox has piece not present in byColour", Square: sq, Piece: piece}
		}

		hasPieceByType := false
		for pt := Pawn; pt <= King; pt++ {
			if p.byType[pt].Has(sq) {
				if hasPieceByType {
					return &PositionError{Message: "multiple pieces in byType for square", Square: sq}
				}
				hasPieceByType = true
			}
		}
		if p.byColour[White].Has(sq) && p.byColour[Black].Has(sq) {
			return &PositionError{Message: "square has pieces of both colours", Square: sq}
		}
	}

	if p.state.castlingRights&^CastlingAll != 0 {
		return &PositionError{Message: "invalid castling rights"}
	}

	if p.byType[NoPieceType] != 0 {
		return &PositionError{Message: "byType has NoPieceType set"}
	}

	var kingErr error
	var whiteKingCount int
	(p.byType[King] & p.byColour[White]).EachSquareSlow(func(sq Square) {
		whiteKingCount++
		if p.kingSquare[White] != sq {
			kingErr = &PositionError{Message: "white king square mismatch", Square: sq}
		}
	})
	if kingErr != nil {
		return kingErr
	}

	if whiteKingCount == 0 && p.kingSquare[White] != NoSquare {
		return &PositionError{Message: "white king square is not NoSquare but there is no white king on the board"}
	}

	var blackKingCount int
	(p.byType[King] & p.byColour[Black]).EachSquareSlow(func(sq Square) {
		blackKingCount++
		if p.kingSquare[Black] != sq {
			kingErr = &PositionError{Message: "black king square mismatch", Square: sq}
		}
	})
	if kingErr != nil {
		return kingErr
	}

	if blackKingCount == 0 && p.kingSquare[Black] != NoSquare {
		return &PositionError{Message: "black king square is not NoSquare but there is no black king on the board"}
	}

	return nil
}

func (p *Position) String() string {
	var sb strings.Builder
	sb.WriteString("\n")
	for rank := Rank8; rank >= Rank1; rank-- {
		sb.WriteString(rank.String())
		sb.WriteString(" ")
		for file := FileA; file <= FileH; file++ {
			sq, _ := SquareFromFileAndRank(file, rank)
			piece := p.mailbox[sq]
			if piece.Type() == NoPieceType {
				sb.WriteString(".")
			} else {
				sb.WriteString(piece.String())
			}
			sb.WriteString(" ")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("  a b c d e f g h \n")
	return sb.String()
}

func (p *Position) PiecesByColour(colour Colour) Bitboard {
	return p.byColour[colour]
}

func (p *Position) Pieces(colour Colour, pieceType PieceType) Bitboard {
	return p.byType[pieceType] & p.byColour[colour]
}

func (p *Position) Occupancy() Bitboard {
	return p.byColour[White] | p.byColour[Black]
}

func (p *Position) PieceAt(square Square) Piece {
	return p.mailbox[square]
}

func (p *Position) Validate() error {
	if p.state.castlingRights&CastlingWhiteQueenside != 0 {
		if p.PieceAt(A1) != WhiteRook {
			return fmt.Errorf("invalid castling rights: White queenside castling right is set, but A1 does not have a white rook")
		}
		if p.PieceAt(E1) != WhiteKing {
			return fmt.Errorf("invalid castling rights: White queenside castling right is set, but E1 does not have a white king")
		}
	}

	if p.state.castlingRights&CastlingWhiteKingside != 0 {
		if p.PieceAt(H1) != WhiteRook {
			return fmt.Errorf("invalid castling rights: White kingside castling right is set, but H1 does not have a white rook")
		}
		if p.PieceAt(E1) != WhiteKing {
			return fmt.Errorf("invalid castling rights: White kingside castling right is set, but E1 does not have a white king")
		}
	}

	if p.state.castlingRights&CastlingBlackQueenside != 0 {
		if p.PieceAt(A8) != BlackRook {
			return fmt.Errorf("invalid castling rights: Black queenside castling right is set, but A8 does not have a black rook")
		}
		if p.PieceAt(E8) != BlackKing {
			return fmt.Errorf("invalid castling rights: Black queenside castling right is set, but E8 does not have a black king")
		}
	}

	if p.state.castlingRights&CastlingBlackKingside != 0 {
		if p.PieceAt(H8) != BlackRook {
			return fmt.Errorf("invalid castling rights: Black kingside castling right is set, but H8 does not have a black rook")
		}
		if p.PieceAt(E8) != BlackKing {
			return fmt.Errorf("invalid castling rights: Black kingside castling right is set, but E8 does not have a black king")
		}
	}

	if p.state.enPassantSquare != NoSquare {
		epSquare := p.state.enPassantSquare
		if epSquare.Rank() != Rank3 && epSquare.Rank() != Rank6 {
			return fmt.Errorf("invalid en passant square: %s is not on rank 3 or 6", epSquare.String())
		}

		var pawnSquare Square
		if epSquare.Rank() == Rank3 {
			pawnSquare, _ = epSquare.Bitboard().North().PopSquare()
			if p.PieceAt(pawnSquare) != WhitePawn {
				return fmt.Errorf("invalid en passant square: %s is set, but there is no white pawn on %s", epSquare.String(), pawnSquare.String())
			}
		}

		if epSquare.Rank() == Rank6 {
			pawnSquare, _ = epSquare.Bitboard().South().PopSquare()
			if p.PieceAt(pawnSquare) != BlackPawn {
				return fmt.Errorf("invalid en passant square: %s is set, but there is no black pawn on %s", epSquare.String(), pawnSquare.String())
			}
		}

		if p.mailbox[epSquare] != NoPiece {
			return fmt.Errorf("invalid en passant square: %s is set, but there is a piece on that square", epSquare.String())
		}

		switch p.sideToMove {
		case White:
			if epSquare.Rank() != Rank6 {
				return fmt.Errorf("invalid en passant square: %s is set, but it is not on rank 6 when it is white's turn to move", epSquare.String())
			}
			file := epSquare.File()
			sq, _ := SquareFromFileAndRank(file, Rank7)
			if p.mailbox[sq] != NoPiece {
				return fmt.Errorf("invalid en passant square: %s is set, but there is a piece on %s", epSquare.String(), sq.String())
			}
		case Black:
			if epSquare.Rank() != Rank3 {
				return fmt.Errorf("invalid en passant square: %s is set, but it is not on rank 3 when it is black's turn to move", epSquare.String())
			}
			file := epSquare.File()
			sq, _ := SquareFromFileAndRank(file, Rank2)
			if p.mailbox[sq] != NoPiece {
				return fmt.Errorf("invalid en passant square: %s is set, but there is a piece on %s", epSquare.String(), sq.String())
			}
		}
	}

	return nil
}

// InCheck returns true if the king of the given colour is in check.
// The position MUST have kings on the board, otherwise it will panic.
func (p *Position) InCheck(c Colour) bool {
	king := p.kingSquare[c]
	return p.isSquareAttacked(king, c.Opposite())
}

// isLastMoveIllegal returns true if the last move made was illegal, i.e. it left the side that moved in check.
func (p *Position) isLastMoveIllegal() bool {
	return p.InCheck(p.sideToMove.Opposite())
}
