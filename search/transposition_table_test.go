package search

import (
	"fmt"
	"math/rand/v2"
	"testing"
	"unsafe"

	"github.com/liamg/chess/board"
)

func TestTranspositionTableEntrySize(t *testing.T) {
	size := unsafe.Sizeof(transpositionTableEntry{})
	if size != 16 {
		t.Errorf("transpositionTableEntry size is %d, expected 16", size)
	}

	if unsafe.Sizeof(indexedEntry{}) > 64 {
		t.Errorf("transpositionTableClusterSize * transpositionTableEntry size is %d, expected <= 64", transpositionTableClusterSize*unsafe.Sizeof(transpositionTableEntry{}))
	}
}

func TestTranspositionTableMemoryAllocation(t *testing.T) {
	sizeMBs := []int{1, 2, 3, 4, 8, 16, 32, 64}
	for _, sizeMB := range sizeMBs {
		t.Run(fmt.Sprintf("Size %dMB", sizeMB), func(t *testing.T) {
			tt := newTranspositionTable(sizeMB)
			entries := len(tt.data)
			usedBytes := int(unsafe.Sizeof(indexedEntry{})) * entries
			if usedBytes > sizeMB*1048576 {
				t.Errorf("Allocated memory exceeds %dMB", sizeMB)
			}
		})
	}
}

func TestTranspositionTablePanicsOnZeroSize(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for zero size, but did not panic")
		}
	}()

	newTranspositionTable(0)
}

func TestTranspositionTableProbeInEntity(t *testing.T) {
	tt := newTranspositionTable(1)
	entry, ok := tt.probe(123, 1)
	if ok {
		t.Errorf("Expected no entry for key 123, but found one: %+v", entry)
	}
	if entry.Bound != unbounded {
		t.Errorf("Expected unbounded entry for key 123, but found: %+v", entry)
	}
}

func TestTranspositionTableStoreAndProbe(t *testing.T) {
	tt := newTranspositionTable(1)
	move := board.NewMove(board.A2, board.A3, board.Capture)
	tt.store(123, 100, move, 5, exact, 1, 1)
	retrieved, ok := tt.probe(123, 1)
	if !ok {
		t.Fatalf("Expected entry for key 123, but found none")
	}
	if retrieved.Key != 123 || retrieved.Score != 100 || retrieved.Move != move || retrieved.Depth != 5 || retrieved.Bound != exact || retrieved.Birth != 1 {
		t.Errorf("Retrieved entry does not match stored values: %+v", retrieved)
	}
}

func TestTranspositionTableProbeForDifferentKeyWithSameLowBits(t *testing.T) {
	tt := newTranspositionTable(1)
	move := board.NewMove(board.A2, board.A3, board.Capture)
	tt.store(123, 100, move, 5, exact, 1, 1)

	// Create a different key that has the same low bits as 123
	key := uint64(123) | (1 << 32) // This will have the same low bits as 123
	retrieved, ok := tt.probe(key, 1)
	if ok {
		t.Errorf("Expected no entry for key %d, but found one: %+v", key, retrieved)
	}
}

func TestTranspositionTableStoreAcrossCluster(t *testing.T) {
	tt := newTranspositionTable(1)

	for i := range transpositionTableClusterSize {
		tt.store(123+uint64(i*len(tt.data)), int16(i*100), board.NewMove(board.A1+board.Square(i), board.A2+board.Square(i), board.Capture), 5, exact, 1, 1)
	}
	for i := range transpositionTableClusterSize {
		retrieved, ok := tt.probe(123+uint64(i*len(tt.data)), 1)
		if !ok {
			t.Errorf("Expected entry for key %d, but found none", 123+(i*len(tt.data)))
		}
		if retrieved.Score != int16(i*100) {
			t.Errorf("Retrieved score for key %d does not match stored value: got %d, want %d", 123+(i*len(tt.data)), retrieved.Score, i*100)
		}
	}
	tt.store(123+uint64(len(tt.data)*transpositionTableClusterSize), int16(10000), board.NewMove(board.A1, board.A2, board.Capture), 5, exact, 1, 1)
	var foundNew bool
	for i := range transpositionTableClusterSize {
		retrieved, ok := tt.probe(123+uint64(i*len(tt.data)), 1)
		if !ok && i > 0 {
			t.Errorf("Expected entry for key %d, but found none", 123+(i*len(tt.data)))
		} else if i == 0 {
			if ok {
				t.Errorf("Expected entry for key %d to be evicted, but found one: %+v", 123+(i*len(tt.data)), retrieved)
			}
			continue
		}
		if retrieved.Score != int16(i*100) {
			if !foundNew && retrieved.Score == 10000 {
				foundNew = true
			} else if retrieved.Score != 10000 {
				t.Errorf("Retrieved score for key %d does not match stored value: got %d, want %d", 123+(i*len(tt.data)), retrieved.Score, i*100)
			} else {
				t.Errorf("New value stored more than once")
			}
		}
	}
}

func TestTranspositionTableReset(t *testing.T) {
	tt := newTranspositionTable(1)
	tt.store(123, 100, board.NewMove(board.A2, board.A3, board.Capture), 5, exact, 1, 1)
	tt.reset()
	entry, ok := tt.probe(123, 1)
	if ok {
		t.Errorf("Expected no entry for key 123 after reset, but found one: %+v", entry)
	}
	if entry.Bound != unbounded {
		t.Errorf("Expected unbounded entry for key 123 after reset, but found: %+v", entry)
	}
}

func TestTranspositionTableReplacement(t *testing.T) {
	tests := []struct {
		name           string
		initialDepth   int8
		initialBirth   int8
		newDepth       int8
		newBirth       int8
		expectReplaced bool
	}{
		{
			name:           "Replace with deeper entry",
			initialDepth:   5,
			initialBirth:   1,
			newDepth:       10,
			newBirth:       2,
			expectReplaced: true,
		},
		{
			name:           "Do not replace with shallower entry",
			initialDepth:   10,
			initialBirth:   1,
			newDepth:       5,
			newBirth:       2,
			expectReplaced: false,
		},
		{
			name:           "Replace with same depth but newer entry",
			initialDepth:   5,
			initialBirth:   1,
			newDepth:       5,
			newBirth:       2,
			expectReplaced: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tt := newTranspositionTable(1)
			tt.store(123, 100, board.NewMove(board.A2, board.A3, board.Capture), test.initialDepth, exact, uint8(test.initialBirth), 1)
			tt.store(123, 200, board.NewMove(board.A3, board.A4, board.Capture), test.newDepth, exact, uint8(test.newBirth), 1)

			retrieved, ok := tt.probe(123, 1)
			if !ok {
				t.Fatalf("Expected entry for key 123, but found none")
			}
			replaced := retrieved.Score == 200
			if replaced != test.expectReplaced {
				t.Errorf("Replacement behavior mismatch: got %v, want %v", replaced, test.expectReplaced)
			}
		})
	}
}

func TestTranspositionTableEvictsShallowestFromCluster(t *testing.T) {
	tt := newTranspositionTable(1)
	for i := range transpositionTableClusterSize {
		tt.store(123+uint64(i*len(tt.data)), int16(i), board.NewMove(board.A1+board.Square(i), board.A2+board.Square(i), board.Capture), 20-int8(i), exact, 1, 1)
	}
	tt.store(123+uint64(transpositionTableClusterSize*len(tt.data)), int16(100), board.NewMove(board.A1, board.A2, board.Capture), 5, exact, 1, 1)
	retrieved, ok := tt.probe(123+uint64(transpositionTableClusterSize*len(tt.data)), 1)
	if !ok {
		t.Fatalf("Expected entry for key %d, but found none", 123+(transpositionTableClusterSize*len(tt.data)))
	}
	if retrieved.Score != 100 {
		t.Errorf("Retrieved score for key %d does not match stored value: got %d, want %d", 123+(transpositionTableClusterSize*len(tt.data)), retrieved.Score, 100)
	}
	// Check that the shallowest entry was evicted
	evictedKey := 123 + uint64((transpositionTableClusterSize-1)*len(tt.data))
	_, ok = tt.probe(evictedKey, 1)
	if ok {
		t.Errorf("Expected entry for key %d to be evicted, but found one", evictedKey)
	}

	for i := range transpositionTableClusterSize + 1 {
		if i == 3 {
			continue // This was the evicted entry
		}
		_, ok := tt.probe(123+uint64(i*len(tt.data)), 1)
		if !ok {
			t.Errorf("Expected entry for key %d, but found none", 123+(i*len(tt.data)))
		}
	}
}

func TestTranspositionTableEvictsOldestFromCluster(t *testing.T) {
	tt := newTranspositionTable(1)
	for i := range transpositionTableClusterSize {
		tt.store(123+uint64(i*len(tt.data)), int16(i), board.NewMove(board.A1+board.Square(i), board.A2+board.Square(i), board.Capture), 20, exact, uint8((2+i)%4), 1)
	}
	tt.store(123+uint64(transpositionTableClusterSize*len(tt.data)), int16(100), board.NewMove(board.A1, board.A2, board.Capture), 5, exact, transpositionTableClusterSize+1, 1)
	retrieved, ok := tt.probe(123+uint64(transpositionTableClusterSize*len(tt.data)), 1)
	if !ok {
		t.Fatalf("Expected entry for key %d, but found none", 123+(transpositionTableClusterSize*len(tt.data)))
	}
	if retrieved.Score != 100 {
		t.Errorf("Retrieved score for key %d does not match stored value: got %d, want %d", 123+(transpositionTableClusterSize*len(tt.data)), retrieved.Score, 100)
	}
	// Check that the oldest entry was evicted
	evictedKey := uint64(123 + (2 * len(tt.data)))
	_, ok = tt.probe(evictedKey, 1)
	if ok {
		t.Errorf("Expected entry for key %d to be evicted, but found one", evictedKey)
	}

	for i := range transpositionTableClusterSize + 1 {
		if i == 2 {
			continue // This was the evicted entry
		}
		_, ok := tt.probe(123+uint64(i*len(tt.data)), 1)
		if !ok {
			t.Errorf("Expected entry for key %d, but found none", 123+(i*len(tt.data)))
		}
	}
}

func TestTranspositionTableEvictsOldestFromClusterWithWraparound(t *testing.T) {
	tt := newTranspositionTable(1)
	for i := range transpositionTableClusterSize {
		age := uint8(3 - i)
		if i == 1 {
			age = 250
		}
		tt.store(123+uint64(i*len(tt.data)), int16(i), board.NewMove(board.A1+board.Square(i), board.A2+board.Square(i), board.Capture), 20, exact, age, 1)
	}
	tt.store(123+uint64(transpositionTableClusterSize*len(tt.data)), int16(100), board.NewMove(board.A1, board.A2, board.Capture), 5, exact, transpositionTableClusterSize+1, 1)
	retrieved, ok := tt.probe(123+uint64(transpositionTableClusterSize*len(tt.data)), 1)
	if !ok {
		t.Fatalf("Expected entry for key %d, but found none", 123+(transpositionTableClusterSize*len(tt.data)))
	}
	if retrieved.Score != 100 {
		t.Errorf("Retrieved score for key %d does not match stored value: got %d, want %d", 123+(transpositionTableClusterSize*len(tt.data)), retrieved.Score, 100)
	}
	// Check that the oldest entry was evicted
	evictedKey := 123 + uint64(len(tt.data))
	_, ok = tt.probe(evictedKey, 1)
	if ok {
		t.Errorf("Expected entry for key %d to be evicted, but found one", evictedKey)
	}
	for i := range transpositionTableClusterSize + 1 {
		if i == 1 {
			continue // This was the evicted entry
		}
		_, ok := tt.probe(123+uint64(i*len(tt.data)), 1)
		if !ok {
			t.Errorf("Expected entry for key %d, but found none", 123+(i*len(tt.data)))
		}
	}
}

func TestTranspositionTableFuzz(t *testing.T) {
	tt := newTranspositionTable(1)
	rnd := rand.New(rand.NewPCG(1, 2))
	count := 1000000
	store := make(map[uint64]transpositionTableEntry, count)
	age := uint8(0)
	for i := range count {
		key := rnd.Uint64()
		if _, ok := store[key]; ok {
			continue // Skip duplicate keys
		}
		entry := transpositionTableEntry{
			Key:   key,
			Score: int16(rnd.IntN(200) - 100),
			Move:  board.NewMove(board.Square(rnd.IntN(64)), board.Square(rnd.IntN(64)), board.Capture),
			Depth: int8(rnd.IntN(20)),
			Bound: bound(1 + rnd.IntN(3)),
			Birth: age,
		}
		store[key] = entry
		tt.store(key, entry.Score, entry.Move, entry.Depth, entry.Bound, entry.Birth, 1)

		if i&0xffff == 0 {
			age++
		}
	}
	for key, val := range store {
		retrieved, ok := tt.probe(key, 1)
		if ok && retrieved != val {
			t.Errorf("Retrieved entry for key %d does not match stored value: got %+v, want %+v", key, retrieved, val)
		}
	}
}
