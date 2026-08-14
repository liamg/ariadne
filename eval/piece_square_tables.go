package eval

import "github.com/liamg/ariadne/board"

const (
	ScoreKing   Score = 0
	ScoreQueen  Score = 900
	ScoreRook   Score = 500
	ScoreBishop Score = 330
	ScoreKnight Score = 320
	ScorePawn   Score = 100
)

// TODO: switch to PeSTO later?

var scoreTable = [7]Score{
	board.NoPieceType: 0,
	board.Pawn:        ScorePawn,
	board.Knight:      ScoreKnight,
	board.Bishop:      ScoreBishop,
	board.Rook:        ScoreRook,
	board.Queen:       ScoreQueen,
	board.King:        ScoreKing,
}

var pieceSquareTablesMidGame = [2][7][64]int16{
	// the below layouts are to be viewed as from whites point of view - ignore the fact it's indexed as black - magic!
	board.Black: {
		board.NoPieceType: {},
		board.Pawn: {
			0, 0, 0, 0, 0, 0, 0, 0,
			50, 50, 50, 50, 50, 50, 50, 50,
			10, 10, 20, 30, 30, 20, 10, 10,
			5, 5, 10, 25, 25, 10, 5, 5,
			0, 0, 0, 20, 20, 0, 0, 0,
			5, -5, -10, 0, 0, -10, -5, 5,
			5, 10, 10, -20, -20, 10, 10, 5,
			0, 0, 0, 0, 0, 0, 0, 0,
		},
		board.Knight: {
			-50, -40, -30, -30, -30, -30, -40, -50,
			-40, -20, 0, 0, 0, 0, -20, -40,
			-30, 0, 10, 15, 15, 10, 0, -30,
			-30, 5, 15, 20, 20, 15, 5, -30,
			-30, 0, 15, 20, 20, 15, 0, -30,
			-30, 5, 10, 15, 15, 10, 5, -30,
			-40, -20, 0, 5, 5, 0, -20, -40,
			-50, -40, -30, -30, -30, -30, -40, -50,
		},
		board.Bishop: {
			-20, -10, -10, -10, -10, -10, -10, -20,
			-10, 0, 0, 0, 0, 0, 0, -10,
			-10, 0, 5, 10, 10, 5, 0, -10,
			-10, 5, 5, 10, 10, 5, 5, -10,
			-10, 0, 10, 10, 10, 10, 0, -10,
			-10, 10, 10, 10, 10, 10, 10, -10,
			-10, 5, 0, 0, 0, 0, 5, -10,
			-20, -10, -10, -10, -10, -10, -10, -20,
		},
		board.Rook: {
			0, 0, 0, 0, 0, 0, 0, 0,
			5, 10, 10, 10, 10, 10, 10, 5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			0, 0, 0, 5, 5, 0, 0, 0,
		},
		board.Queen: {
			-20, -10, -10, -5, -5, -10, -10, -20,
			-10, 0, 0, 0, 0, 0, 0, -10,
			-10, 0, 5, 5, 5, 5, 0, -10,
			-5, 0, 5, 5, 5, 5, 0, -5,
			0, 0, 5, 5, 5, 5, 0, -5,
			-10, 5, 5, 5, 5, 5, 0, -10,
			-10, 0, 5, 0, 0, 0, 0, -10,
			-20, -10, -10, -5, -5, -10, -10, -20,
		},
		board.King: {
			-30, -40, -40, -50, -50, -40, -40, -30,
			-30, -40, -40, -50, -50, -40, -40, -30,
			-30, -40, -40, -50, -50, -40, -40, -30,
			-30, -40, -40, -50, -50, -40, -40, -30,
			-20, -30, -30, -40, -40, -30, -30, -20,
			-10, -20, -20, -20, -20, -20, -20, -10,
			20, 20, 0, 0, 0, 0, 20, 20,
			20, 30, 10, 0, 0, 10, 30, 20,
		},
	},
}

var pieceSquareTablesEndGame = [2][7][64]int16{
	// the below layouts are to be viewed as from whites point of view - ignore the fact it's indexed as black - magic!
	board.Black: {
		board.NoPieceType: {},
		board.Pawn: {
			0, 0, 0, 0, 0, 0, 0, 0,
			50, 50, 50, 50, 50, 50, 50, 50,
			10, 10, 20, 30, 30, 20, 10, 10,
			5, 5, 10, 25, 25, 10, 5, 5,
			0, 0, 0, 20, 20, 0, 0, 0,
			5, -5, -10, 0, 0, -10, -5, 5,
			5, 10, 10, -20, -20, 10, 10, 5,
			0, 0, 0, 0, 0, 0, 0, 0,
		},
		board.Knight: {
			-50, -40, -30, -30, -30, -30, -40, -50,
			-40, -20, 0, 0, 0, 0, -20, -40,
			-30, 0, 10, 15, 15, 10, 0, -30,
			-30, 5, 15, 20, 20, 15, 5, -30,
			-30, 0, 15, 20, 20, 15, 0, -30,
			-30, 5, 10, 15, 15, 10, 5, -30,
			-40, -20, 0, 5, 5, 0, -20, -40,
			-50, -40, -30, -30, -30, -30, -40, -50,
		},
		board.Bishop: {
			-20, -10, -10, -10, -10, -10, -10, -20,
			-10, 0, 0, 0, 0, 0, 0, -10,
			-10, 0, 5, 10, 10, 5, 0, -10,
			-10, 5, 5, 10, 10, 5, 5, -10,
			-10, 0, 10, 10, 10, 10, 0, -10,
			-10, 10, 10, 10, 10, 10, 10, -10,
			-10, 5, 0, 0, 0, 0, 5, -10,
			-20, -10, -10, -10, -10, -10, -10, -20,
		},
		board.Rook: {
			0, 0, 0, 0, 0, 0, 0, 0,
			5, 10, 10, 10, 10, 10, 10, 5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			-5, 0, 0, 0, 0, 0, 0, -5,
			0, 0, 0, 5, 5, 0, 0, 0,
		},
		board.Queen: {
			-20, -10, -10, -5, -5, -10, -10, -20,
			-10, 0, 0, 0, 0, 0, 0, -10,
			-10, 0, 5, 5, 5, 5, 0, -10,
			-5, 0, 5, 5, 5, 5, 0, -5,
			0, 0, 5, 5, 5, 5, 0, -5,
			-10, 5, 5, 5, 5, 5, 0, -10,
			-10, 0, 5, 0, 0, 0, 0, -10,
			-20, -10, -10, -5, -5, -10, -10, -20,
		},
		// the king wants the centre once the queens are off, rather than the
		// corner it hides in during the middlegame
		board.King: {
			-50, -40, -30, -20, -20, -30, -40, -50,
			-30, -20, -10, 0, 0, -10, -20, -30,
			-30, -10, 20, 30, 30, 20, -10, -30,
			-30, -10, 30, 40, 40, 30, -10, -30,
			-30, -10, 30, 40, 40, 30, -10, -30,
			-30, -10, 20, 30, 30, 20, -10, -30,
			-30, -30, 0, 0, 0, 0, -30, -30,
			-50, -30, -30, -30, -30, -30, -30, -50,
		},
	},
}

func init() {
	for pt := board.Pawn; pt <= board.King; pt++ {
		for sq := board.A1; sq <= board.H8; sq++ {
			pieceSquareTablesMidGame[board.Black][pt][sq] = pieceSquareTablesMidGame[board.Black][pt][sq] + int16(scoreTable[pt])
		}
		for sq := board.A1; sq <= board.H8; sq++ {
			wsq := sq ^ 56
			pieceSquareTablesMidGame[board.White][pt][wsq] = pieceSquareTablesMidGame[board.Black][pt][sq]
		}
	}
	for pt := board.Pawn; pt <= board.King; pt++ {
		for sq := board.A1; sq <= board.H8; sq++ {
			pieceSquareTablesEndGame[board.Black][pt][sq] = pieceSquareTablesEndGame[board.Black][pt][sq] + int16(scoreTable[pt])
		}
		for sq := board.A1; sq <= board.H8; sq++ {
			wsq := sq ^ 56
			pieceSquareTablesEndGame[board.White][pt][wsq] = pieceSquareTablesEndGame[board.Black][pt][sq]
		}
	}
}
