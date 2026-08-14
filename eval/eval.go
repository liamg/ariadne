package eval

import "github.com/liamg/ariadne/board"

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
	kingDangerDivisor  = 32
	kingDangerMax      = 500
)

var passedWhitePawnsMG = []int{
	0, 0, 10, 15, 25, 45, 80, 0,
}

var passedWhitePawnsEG = []int{
	0, 0, 20, 30, 50, 90, 150, 0,
}

var kingAttackWeights = [7]int16{
	board.NoPieceType: 0,
	board.Pawn:        0,
	board.Knight:      2,
	board.Bishop:      2,
	board.Rook:        3,
	board.Queen:       5,
	board.King:        0,
}

var (
	passedBlackPawnsMG []int
	passedBlackPawnsEG []int
)

var kingZones [64]board.Bitboard

func init() {
	passedBlackPawnsMG = make([]int, 8)
	passedBlackPawnsEG = make([]int, 8)
	for i := range 8 {
		passedBlackPawnsMG[i] = passedWhitePawnsMG[7-i]
		passedBlackPawnsEG[i] = passedWhitePawnsEG[7-i]
	}

	for sq := board.A1; sq <= board.H8; sq++ {
		file := min(max(sq.File(), board.FileB), board.FileG)
		rank := min(max(sq.Rank(), board.Rank2), board.Rank7)
		centreSq, _ := board.SquareFromFileAndRank(file, rank)
		centre := centreSq.Bitboard()
		kingZones[sq] = centre | centre.NorthWest() | centre.North() | centre.NorthEast() | centre.West() | centre.East() | centre.SouthWest() | centre.South() | centre.SouthEast()
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

type Evaluator struct {
	pawnsTable pawnsTable
}

func New() *Evaluator {
	return &Evaluator{
		pawnsTable: newPawnsTable(),
	}
}

func (e *Evaluator) Evaluate(p *board.Position) Score {
	var sq board.Square

	occ := p.Occupancy()

	var phase int

	whitePawns := p.Pieces(board.White, board.Pawn)
	blackPawns := p.Pieces(board.Black, board.Pawn)

	whitePieces := p.PiecesByColour(board.White)
	blackPieces := p.PiecesByColour(board.Black)

	whiteKing := p.KingSquare(board.White)
	blackKing := p.KingSquare(board.Black)

	pawnEntry, ok := e.pawnsTable.probe(whitePawns, blackPawns)
	if !ok {
		var pawnScoreMidGame Score
		var pawnScoreEndGame Score
		doubledWhitePawns := doubledPawns(board.White, whitePawns)
		doubledBlackPawns := doubledPawns(board.Black, blackPawns)

		pawnScoreMidGame += Score((doubledWhitePawns.Count() - doubledBlackPawns.Count()) * doubledPawnMG)
		pawnScoreEndGame += Score((doubledWhitePawns.Count() - doubledBlackPawns.Count()) * doubledPawnEG)

		isolatedWhitePawns := isolatedPawns(whitePawns)
		isolatedBlackPawns := isolatedPawns(blackPawns)

		passedWhitePawns := passedPawns(board.White, whitePawns, blackPawns)
		passedBlackPawns := passedPawns(board.Black, blackPawns, whitePawns)

		backwardWhitePawns := backwardPawns(board.White, whitePawns, blackPawns)
		backwardBlackPawns := backwardPawns(board.Black, blackPawns, whitePawns)

		pawnScoreMidGame += Score((backwardWhitePawns.Count() - backwardBlackPawns.Count()) * backwardPawnMG)
		pawnScoreEndGame += Score((backwardWhitePawns.Count() - backwardBlackPawns.Count()) * backwardPawnEG)

		// a backward pawn on a file with no enemy pawn ahead of it is the one a rook
		// can sit opposite, so it takes an extra penalty on top
		unopposedWhitePawns := backwardWhitePawns & unopposedPawns(board.White, whitePawns, blackPawns)
		unopposedBlackPawns := backwardBlackPawns & unopposedPawns(board.Black, blackPawns, whitePawns)

		pawnScoreMidGame += Score((unopposedWhitePawns.Count() - unopposedBlackPawns.Count()) * unopposedPawnMG)
		pawnScoreEndGame += Score((unopposedWhitePawns.Count() - unopposedBlackPawns.Count()) * unopposedPawnEG)

		for passedWhitePawns != 0 {
			sq, passedWhitePawns = passedWhitePawns.PopSquare()
			pawnScoreMidGame += Score(passedWhitePawnsMG[sq.Rank()])
			pawnScoreEndGame += Score(passedWhitePawnsEG[sq.Rank()])
		}
		for passedBlackPawns != 0 {
			sq, passedBlackPawns = passedBlackPawns.PopSquare()
			pawnScoreMidGame -= Score(passedBlackPawnsMG[sq.Rank()])
			pawnScoreEndGame -= Score(passedBlackPawnsEG[sq.Rank()])
		}

		pawnScoreMidGame += Score((isolatedWhitePawns.Count() - isolatedBlackPawns.Count()) * isolatedPawnMG)
		pawnScoreEndGame += Score((isolatedWhitePawns.Count() - isolatedBlackPawns.Count()) * isolatedPawnEG)

		pawnEntry.whitePawns = whitePawns
		pawnEntry.blackPawns = blackPawns
		pawnEntry.kingSquare[0] = 0xFF
		pawnEntry.kingSquare[1] = 0xFF
		pawnEntry.midGameScore = pawnScoreMidGame
		pawnEntry.endGameScore = pawnScoreEndGame
	}

	if pawnEntry.kingSquare[board.White] != byte(whiteKing) {
		whitePenalty := Score(calculateKingShelterAndStormPenalty(whiteKing, whitePawns, blackPawns))
		pawnEntry.shelter[board.White] = int16(whitePenalty)
		pawnEntry.kingSquare[board.White] = byte(whiteKing)
	}

	if pawnEntry.kingSquare[board.Black] != byte(blackKing) {
		blackPenalty := Score(calculateKingShelterAndStormPenalty(blackKing.FlipVertical(), blackPawns.FlipVertical(), whitePawns.FlipVertical()))
		pawnEntry.shelter[board.Black] = int16(blackPenalty)
		pawnEntry.kingSquare[board.Black] = byte(blackKing)
	}

	shelterPenaltyMidGame := Score(pawnEntry.shelter[board.White] - pawnEntry.shelter[board.Black])
	scoreMidGame := pawnEntry.midGameScore - shelterPenaltyMidGame
	scoreEndGame := pawnEntry.endGameScore

	// white pieces attacking the black king
	whiteKingAttackers := int16(0)
	whiteKingAttackWeight := int16(0)
	// black pieces attacking the white king
	blackKingAttackers := int16(0)
	blackKingAttackWeight := int16(0)

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

				att := p.AttacksWithCustomOccupancy(pt, sq, occ) &^ whitePieces
				if hits := int16((att & kingZones[blackKing]).Count()); hits > 0 {
					whiteKingAttackers++
					whiteKingAttackWeight += (kingAttackWeights[pt]) * hits
				}
				mobility = int16(att.Count())
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
				att := p.AttacksWithCustomOccupancy(pt, sq, occ) &^ blackPieces
				if hits := int16((att & kingZones[whiteKing]).Count()); hits > 0 {
					blackKingAttackers++
					blackKingAttackWeight += (kingAttackWeights[pt]) * hits
				}
				mobility = int16(att.Count())
			}
			scoreMidGame -= Score((mobility * pieceMobilityWeightsMidGame[pt]) + pieceSquareTablesMidGame[board.Black][pt][sq])
			scoreEndGame -= Score((mobility * pieceMobilityWeightsEndGame[pt]) + pieceSquareTablesEndGame[board.Black][pt][sq])
		}
	}

	if whiteKingAttackers >= 2 {
		danger := Score(whiteKingAttackWeight) * Score(whiteKingAttackers)
		scoreMidGame += min((danger*danger)/kingDangerDivisor, kingDangerMax)
	}
	if blackKingAttackers >= 2 {
		danger := Score(blackKingAttackWeight) * Score(blackKingAttackers)
		scoreMidGame -= min((danger*danger)/kingDangerDivisor, kingDangerMax)
	}
	phase = min(maxPhase, phase)
	score := (scoreMidGame*(Score(phase)) + scoreEndGame*(Score(maxPhase-phase))) / maxPhase

	if p.SideToMove() == board.Black {
		score = -score
	}
	return score
}

// for each square, a bitboard of the squares ahead of it (north)
var ranksFrom [9]board.Bitboard

func init() {
	for rank := board.Rank1; rank <= board.Rank8; rank++ {
		ranksFrom[rank] = ^board.Bitboard(0) << (rank * 8)
	}
}

// done from the point of view of white - we flip on the way in
func calculateKingShelterAndStormPenalty(king board.Square, friendlyPawns, enemyPawns board.Bitboard) int16 {
	clampedFile := min(max(king.File(), board.FileB), board.FileG)

	rank := king.Rank()

	penalty := int16(0)
	for file := clampedFile - 1; file <= clampedFile+1; file++ {
		shelterPawns := friendlyPawns & ranksFrom[rank] & file.Mask()
		// no pawns on this file on or ahead of the kings rank
		var blocking board.Rank
		if shelterPawns == 0 {
			penalty += 30
		} else {
			nearestPawn, _ := shelterPawns.PopSquare()
			blocking = nearestPawn.Rank()
			switch blocking {
			case board.Rank2:
				// nothing to do, no penalty if king is protected nicely
			case board.Rank3:
				penalty += 10
			case board.Rank4:
				penalty += 20
			default:
				penalty += 30
			}
		}
		stormingPawns := enemyPawns & ranksFrom[rank+1] & file.Mask()
		if stormingPawns != 0 {
			nearestStormer, _ := stormingPawns.PopSquare()
			stormingRank := nearestStormer.Rank()
			if nearestStormer.Bitboard().South()&shelterPawns != 0 {
				switch stormingRank {
				case board.Rank3:
					penalty += 5
				case board.Rank4:
					penalty += 3
				case board.Rank5:
					penalty += 1
				}
			} else {
				switch stormingRank {
				case board.Rank2:
					penalty += 30
				case board.Rank3:
					penalty += 20
				case board.Rank4:
					penalty += 10
				case board.Rank5:
					penalty += 5
				}
			}
		}

	}

	return penalty
}
