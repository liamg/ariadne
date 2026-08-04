package eval

import "github.com/liamg/chess/board"

var phaseWeights = [7]int{
	board.NoPieceType: 0,
	board.Pawn:        0,
	board.Knight:      1,
	board.Bishop:      1,
	board.Rook:        2,
	board.Queen:       4,
	board.King:        0,
}

const maxPhase = 24

var pieceMobilityWeightsMidGame = [7]int16{
	board.NoPieceType: 0,
	board.Pawn:        0,
	board.Knight:      4,
	board.Bishop:      5,
	board.Rook:        2,
	board.Queen:       1,
	board.King:        0,
}

var pieceMobilityWeightsEndGame = [7]int16{
	board.NoPieceType: 0,
	board.Pawn:        0,
	board.Knight:      4,
	board.Bishop:      5,
	board.Rook:        4,
	board.Queen:       2,
	board.King:        0,
}

const (
	bishopPairMG       = 40
	bishopPairEG       = 50
	rookOpenFileMG     = 20
	rookOpenFileEG     = 10
	rookHalfOpenFileMG = 10
	rookHalfOpenFileEG = 5
	doubledPawnMG      = -10
	doubledPawnEG      = -20
	isolatedPawnMG     = -15
	isolatedPawnEG     = -20
	backwardPawnMG     = -10
	backwardPawnEG     = -15
	unopposedPawnMG    = -8
	unopposedPawnEG    = -10
)

var passedWhitePawnsMG = []int{
	0, 0, 10, 15, 25, 45, 80, 0,
}

var passedWhitePawnsEG = []int{
	0, 0, 20, 30, 50, 90, 150, 0,
}

var (
	passedBlackPawnsMG []int
	passedBlackPawnsEG []int
)

func init() {
	passedBlackPawnsMG = make([]int, 8)
	passedBlackPawnsEG = make([]int, 8)
	for i := 0; i < 8; i++ {
		passedBlackPawnsMG[i] = passedWhitePawnsMG[7-i]
		passedBlackPawnsEG[i] = passedWhitePawnsEG[7-i]
	}
}

func northFill(pawns board.Bitboard) board.Bitboard {
	pawns |= pawns << 8
	pawns |= pawns << 16
	pawns |= pawns << 32
	return pawns
}

func southFill(pawns board.Bitboard) board.Bitboard {
	pawns |= pawns >> 8
	pawns |= pawns >> 16
	pawns |= pawns >> 32
	return pawns
}

// pawns with another friendly pawn on the same file. yields n-1 per file
func doubledPawns(colour board.Colour, pawns board.Bitboard) board.Bitboard {
	if colour == board.White {
		return pawns & northFill(pawns<<8)
	}
	return pawns & southFill(pawns>>8)
}

// pawns with no friendly pawn on either adjacent file. fill to whole files
// first, then spread sideways, so the ranks involved do not matter
func isolatedPawns(pawns board.Bitboard) board.Bitboard {
	files := northFill(pawns) | southFill(pawns)
	return pawns &^ (files.West() | files.East())
}

// pawns with no enemy pawn ahead of them on their own or an adjacent file.
// the shift before the fill makes it strictly ahead - an enemy pawn level with
// ours on the next file cannot stop it, as it captures away from us
func passedPawns(colour board.Colour, pawns, enemyPawns board.Bitboard) board.Bitboard {
	var stops board.Bitboard
	if colour == board.White {
		stops = southFill(enemyPawns >> 8)
	} else {
		stops = northFill(enemyPawns << 8)
	}
	stops |= stops.West() | stops.East()
	return pawns &^ stops
}

// pawns whose stop square is attacked by an enemy pawn and lies outside our own
// pawns' attack span - so they can never advance safely, and no friendly pawn
// can ever come up to defend them
func backwardPawns(colour board.Colour, pawns, enemyPawns board.Bitboard) board.Bitboard {
	if colour == board.White {
		span := northFill(pawns.NorthEast() | pawns.NorthWest())
		enemyAttacks := enemyPawns.SouthEast() | enemyPawns.SouthWest()
		return (pawns.North() & enemyAttacks &^ span).South()
	}
	span := southFill(pawns.SouthEast() | pawns.SouthWest())
	enemyAttacks := enemyPawns.NorthEast() | enemyPawns.NorthWest()
	return (pawns.South() & enemyAttacks &^ span).North()
}

// pawns with no enemy pawn ahead of them on their own file. a weak pawn on such
// a file is the one a rook can actually attack
func unopposedPawns(colour board.Colour, pawns, enemyPawns board.Bitboard) board.Bitboard {
	if colour == board.White {
		return pawns &^ southFill(enemyPawns)
	}
	return pawns &^ northFill(enemyPawns)
}

func Evaluate(p *board.Position) Score {
	var scoreMidGame Score
	var scoreEndGame Score
	var sq board.Square

	occ := p.Occupancy()

	var phase int

	whitePawns := p.Pieces(board.White, board.Pawn)
	blackPawns := p.Pieces(board.Black, board.Pawn)

	doubledWhitePawns := doubledPawns(board.White, whitePawns)
	doubledBlackPawns := doubledPawns(board.Black, blackPawns)

	scoreMidGame += Score((doubledWhitePawns.Count() - doubledBlackPawns.Count()) * doubledPawnMG)
	scoreEndGame += Score((doubledWhitePawns.Count() - doubledBlackPawns.Count()) * doubledPawnEG)

	isolatedWhitePawns := isolatedPawns(whitePawns)
	isolatedBlackPawns := isolatedPawns(blackPawns)

	passedWhitePawns := passedPawns(board.White, whitePawns, blackPawns)
	passedBlackPawns := passedPawns(board.Black, blackPawns, whitePawns)

	backwardWhitePawns := backwardPawns(board.White, whitePawns, blackPawns)
	backwardBlackPawns := backwardPawns(board.Black, blackPawns, whitePawns)

	scoreMidGame += Score((backwardWhitePawns.Count() - backwardBlackPawns.Count()) * backwardPawnMG)
	scoreEndGame += Score((backwardWhitePawns.Count() - backwardBlackPawns.Count()) * backwardPawnEG)

	// a backward pawn on a file with no enemy pawn ahead of it is the one a rook
	// can sit opposite, so it takes an extra penalty on top
	unopposedWhitePawns := backwardWhitePawns & unopposedPawns(board.White, whitePawns, blackPawns)
	unopposedBlackPawns := backwardBlackPawns & unopposedPawns(board.Black, blackPawns, whitePawns)

	scoreMidGame += Score((unopposedWhitePawns.Count() - unopposedBlackPawns.Count()) * unopposedPawnMG)
	scoreEndGame += Score((unopposedWhitePawns.Count() - unopposedBlackPawns.Count()) * unopposedPawnEG)

	for passedWhitePawns != 0 {
		sq, passedWhitePawns = passedWhitePawns.PopSquare()
		scoreMidGame += Score(passedWhitePawnsMG[sq.Rank()])
		scoreEndGame += Score(passedWhitePawnsEG[sq.Rank()])
	}
	for passedBlackPawns != 0 {
		sq, passedBlackPawns = passedBlackPawns.PopSquare()
		scoreMidGame -= Score(passedBlackPawnsMG[sq.Rank()])
		scoreEndGame -= Score(passedBlackPawnsEG[sq.Rank()])
	}

	scoreMidGame += Score((isolatedWhitePawns.Count() - isolatedBlackPawns.Count()) * isolatedPawnMG)
	scoreEndGame += Score((isolatedWhitePawns.Count() - isolatedBlackPawns.Count()) * isolatedPawnEG)

	whitePieces := p.PiecesByColour(board.White)
	blackPieces := p.PiecesByColour(board.Black)

	for pt := board.Pawn; pt <= board.King; pt++ {

		var mobility int16
		whitePiecesByType := p.Pieces(board.White, pt)
		blackPiecesByType := p.Pieces(board.Black, pt)
		whiteCount := whitePiecesByType.Count()
		blackCount := blackPiecesByType.Count()
		switch pt {
		case board.Bishop:
			if whiteCount > 1 {
				scoreMidGame += bishopPairMG
				scoreEndGame += bishopPairEG
			}
			if blackCount > 1 {
				scoreMidGame -= bishopPairMG
				scoreEndGame -= bishopPairEG
			}
		}
		phase += (whiteCount + blackCount) * phaseWeights[pt]
		for whitePiecesByType != 0 {
			sq, whitePiecesByType = whitePiecesByType.PopSquare()
			if pt > board.Pawn && pt < board.King {
				if pt == board.Rook {
					mask := sq.File().Mask()
					if whitePawns&mask == 0 {
						if blackPawns&mask == 0 {
							scoreMidGame += rookOpenFileMG
							scoreEndGame += rookOpenFileEG
						} else {
							scoreMidGame += rookHalfOpenFileMG
							scoreEndGame += rookHalfOpenFileEG
						}
					}
				}

				mobility = int16((p.AttacksWithCustomOccupancy(pt, sq, occ) &^ whitePieces).Count())
			}
			scoreMidGame += Score((mobility * pieceMobilityWeightsMidGame[pt]) + pieceSquareTablesMidGame[board.White][pt][sq])
			scoreEndGame += Score((mobility * pieceMobilityWeightsEndGame[pt]) + pieceSquareTablesEndGame[board.White][pt][sq])
		}
		for blackPiecesByType != 0 {
			sq, blackPiecesByType = blackPiecesByType.PopSquare()
			if pt > board.Pawn && pt < board.King {
				if pt == board.Rook {
					mask := sq.File().Mask()
					if blackPawns&mask == 0 {
						if whitePawns&mask == 0 {
							scoreMidGame -= rookOpenFileMG
							scoreEndGame -= rookOpenFileEG
						} else {
							scoreMidGame -= rookHalfOpenFileMG
							scoreEndGame -= rookHalfOpenFileEG
						}
					}
				}
				mobility = int16((p.AttacksWithCustomOccupancy(pt, sq, occ) &^ blackPieces).Count())
			}
			scoreMidGame -= Score((mobility * pieceMobilityWeightsMidGame[pt]) + pieceSquareTablesMidGame[board.Black][pt][sq])
			scoreEndGame -= Score((mobility * pieceMobilityWeightsEndGame[pt]) + pieceSquareTablesEndGame[board.Black][pt][sq])
		}
	}

	phase = min(maxPhase, phase)
	score := (scoreMidGame*(Score(phase)) + scoreEndGame*(Score(maxPhase-phase))) / maxPhase

	if p.SideToMove() == board.Black {
		score = -score
	}
	return score
}
