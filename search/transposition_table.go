package search

import (
	"fmt"
	"math/bits"
	"unsafe"

	"github.com/liamg/ariadne/board"
	"github.com/liamg/ariadne/eval"
)

type bound uint8

const (
	unbounded bound = iota
	exact
	lowerBound
	upperBound
)

type indexedEntry [transpositionTableClusterSize]transpositionTableEntry

type transpositionTableEntry struct {
	Key   uint64 // full zobrist hash for now
	Score int16  // NOTE: smaller than the actual underlying score type, but this is fine, as _stored_ scores will always be -32000>=x<=32000
	Move  board.Move
	Depth int8
	Bound bound
	Birth uint8
}

const transpositionTableClusterSize = 4

type transpositionTable struct {
	data []indexedEntry
}

func newTranspositionTable(sizeMB int) *transpositionTable {
	if sizeMB <= 0 {
		panic(fmt.Errorf("transposition table size must be greater than 0MB, got %dMB", sizeMB))
	}
	bytes := 1048576 * sizeMB
	entries := bytes / int(unsafe.Sizeof(indexedEntry{}))
	exp := bits.Len(uint(entries)) - 1
	roundedEntries := 1 << exp

	return &transpositionTable{
		data: make([]indexedEntry, roundedEntries),
	}
}

func (t *transpositionTable) probe(key uint64, ply int) (transpositionTableEntry, bool) {
	for _, entry := range t.data[key&uint64(len(t.data)-1)] {
		if entry.Key == key && entry.Bound != unbounded {
			entry.Score = scoreFromTT(entry.Score, ply)
			return entry, true
		}
	}
	return transpositionTableEntry{}, false
}

func (t *transpositionTable) reset() {
	for i := range t.data {
		for j := range t.data[i] {
			t.data[i][j].Bound = unbounded
		}
	}
}

func (t *transpositionTable) store(key uint64, score int16, move board.Move, depth int, bound bound, age uint8, ply int) {
	cluster := &t.data[key&uint64(len(t.data)-1)]

	targetIndex := -1
	var lowestValue int
	for i, entry := range cluster {

		if entry.Bound == unbounded {
			targetIndex = i
			break
		}

		// overwrite existing value
		if entry.Key == key {
			// ...but only if the new result is at least as deep
			// we don't want to overwrite valuable info with a crappy shallow search
			if depth < int(entry.Depth) {
				return
			}
			cluster[i] = transpositionTableEntry{
				Key:   key,
				Score: scoreToTT(score, ply),
				Move:  move,
				Depth: int8(depth),
				Bound: bound,
				Birth: age,
			}
			return
		}

		if v := int(entry.Depth) - 2*int(age-entry.Birth); v < lowestValue || targetIndex == -1 {
			lowestValue = v
			targetIndex = i
		}
	}

	// we need to replace the OLDEST (or first empty) candidate
	cluster[targetIndex] = transpositionTableEntry{
		Key:   key,
		Score: scoreToTT(score, ply),
		Move:  move,
		Depth: int8(depth),
		Bound: bound,
		Birth: age,
	}
}

const (
	lowMateThreshold  = int16(-eval.Mate + MaxPly)
	highMateThreshold = int16(eval.Mate - MaxPly)
)

func scoreToTT(score int16, ply int) int16 {
	if score < lowMateThreshold {
		return score - int16(ply)
	}
	if score > highMateThreshold {
		return score + int16(ply)
	}
	return score
}

func scoreFromTT(score int16, ply int) int16 {
	if score < lowMateThreshold {
		return score + int16(ply)
	}
	if score > highMateThreshold {
		return score - int16(ply)
	}
	return score
}
