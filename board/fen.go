package board

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	fenEmpty    = `8/8/8/8/8/8/8/8 w - - 0 1`
	fenStarting = `rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1`
	fenKiwiPete = `r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq - 0 1`
)

func ParseFEN(fen string) (*Position, error) {
	pos := EmptyPosition()

	fields := strings.Fields(fen)
	if len(fields) < 4 {
		return nil, fmt.Errorf("invalid FEN: expected at least fields, got %d", len(fields))
	}

	placement := fields[0]
	ranks := strings.Split(placement, "/")
	if len(ranks) != 8 {
		return nil, fmt.Errorf("invalid FEN: expected 8 ranks, got %d", len(ranks))
	}

	for rank := Rank8; rank >= Rank1; rank-- {
		rankStr := ranks[7-rank]
		file := FileA
		for _, r := range rankStr {
			if r >= '1' && r <= '8' {
				file += File(1 + r - '1')
				continue
			}
			var colour Colour
			var pieceType PieceType
			switch r {
			case 'P':
				colour = White
				pieceType = Pawn
			case 'N':
				colour = White
				pieceType = Knight
			case 'B':
				colour = White
				pieceType = Bishop
			case 'R':
				colour = White
				pieceType = Rook
			case 'Q':
				colour = White
				pieceType = Queen
			case 'K':
				colour = White
				pieceType = King
			case 'p':
				colour = Black
				pieceType = Pawn
			case 'n':
				colour = Black
				pieceType = Knight
			case 'b':
				colour = Black
				pieceType = Bishop
			case 'r':
				colour = Black
				pieceType = Rook
			case 'q':
				colour = Black
				pieceType = Queen
			case 'k':
				colour = Black
				pieceType = King
			default:
				return nil, fmt.Errorf("invalid FEN: invalid piece character '%c'", r)
			}
			sq, err := SquareFromFileAndRank(file, rank)
			if err != nil {
				return nil, fmt.Errorf("invalid FEN: invalid square for piece '%c' within %q: %v", r, rankStr, err)
			}
			pos.addPiece(sq, NewPiece(colour, pieceType))
			file++
		}
		if file > FileH+1 {
			return nil, fmt.Errorf("invalid FEN: rank %d has too many pieces", rank+1)
		} else if file < FileH+1 {
			return nil, fmt.Errorf("invalid FEN: rank %d has too few pieces", rank+1)
		}
	}

	sideToMove := fields[1]

	switch sideToMove {
	case "w", "W":
		pos.sideToMove = White
	case "b", "B":
		pos.sideToMove = Black
	default:
		return nil, fmt.Errorf("invalid FEN: side to move must be 'w' or 'b', got '%s'", sideToMove)
	}

	castling := fields[2]

	for _, r := range castling {
		switch r {
		case 'K':
			pos.state.castlingRights |= CastlingWhiteKingside
		case 'Q':
			pos.state.castlingRights |= CastlingWhiteQueenside
		case 'k':
			pos.state.castlingRights |= CastlingBlackKingside
		case 'q':
			pos.state.castlingRights |= CastlingBlackQueenside
		case '-':
			if len(castling) > 1 {
				return nil, fmt.Errorf("invalid FEN: castling rights '-' must be alone")
			}
		default:
			return nil, fmt.Errorf("invalid FEN: invalid castling right '%c'", r)
		}
	}

	enPassant := fields[3]
	if enPassant != "-" {
		sq, err := ParseSquare(enPassant)
		if err != nil {
			return nil, fmt.Errorf("invalid FEN: en passant square is invalid: %v", err)
		}
		pos.state.enPassantSquare = sq
	}

	if len(fields) >= 5 {

		if fields[4][0] < '0' || fields[4][0] > '9' {
			// extednded fen!
			return nil, fmt.Errorf("extended FEN is not supported yet - cannot parse field %q", fields[4])
		}

		halfmoveClock := fields[4]

		hm, err := strconv.ParseInt(halfmoveClock, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid FEN: halfmove clock is not a valid integer: %v", err)
		}
		if hm < 0 {
			return nil, fmt.Errorf("invalid FEN: halfmove clock cannot be negative")
		}
		if hm > 255 {
			return nil, fmt.Errorf("invalid FEN: halfmove clock cannot be greater than 255")
		}
		pos.state.halfMoveClock = uint8(hm)

		if len(fields) >= 6 {
			fullmoveNumber := fields[5]
			fm, err := strconv.ParseInt(fullmoveNumber, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid FEN: fullmove number is not a valid integer: %v", err)
			}
			if fm <= 0 {
				return nil, fmt.Errorf("invalid FEN: fullmove number must be positive")
			}
			pos.fullMoveNumber = int(fm)
		}
	}

	pos.tidy()

	if err := pos.Validate(); err != nil {
		return nil, fmt.Errorf("invalid FEN: %v", err)
	}

	return pos, nil
}

func GenerateFEN(pos *Position) string {
	var sb strings.Builder

	// Piece placement
	for rank := Rank8; rank >= Rank1; rank-- {
		emptyCount := 0
		for file := FileA; file <= FileH; file++ {
			sq, _ := SquareFromFileAndRank(file, rank)
			piece := pos.PieceAt(sq)
			if piece == NoPiece {
				emptyCount++
			} else {
				if emptyCount > 0 {
					sb.WriteString(strconv.Itoa(emptyCount))
					emptyCount = 0
				}
				sb.WriteString(piece.String())
			}
		}
		if emptyCount > 0 {
			sb.WriteString(strconv.Itoa(emptyCount))
		}
		if rank != Rank1 {
			sb.WriteString("/")
		}
	}

	sb.WriteString(" ")

	if pos.sideToMove == White {
		sb.WriteString("w")
	} else {
		sb.WriteString("b")
	}

	sb.WriteString(" ")

	castling := ""
	if pos.state.castlingRights&CastlingWhiteKingside != 0 {
		castling += "K"
	}
	if pos.state.castlingRights&CastlingWhiteQueenside != 0 {
		castling += "Q"
	}
	if pos.state.castlingRights&CastlingBlackKingside != 0 {
		castling += "k"
	}
	if pos.state.castlingRights&CastlingBlackQueenside != 0 {
		castling += "q"
	}
	if castling == "" {
		castling = "-"
	}
	sb.WriteString(castling)

	sb.WriteString(" ")

	if pos.state.enPassantSquare == NoSquare {
		sb.WriteString("-")
	} else {
		sb.WriteString(pos.state.enPassantSquare.String())
	}

	sb.WriteString(" ")

	sb.WriteString(strconv.Itoa(int(pos.state.halfMoveClock)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(pos.fullMoveNumber))

	return sb.String()
}
