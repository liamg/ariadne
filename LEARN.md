# Chess Engine — Learning Path

Build a chess engine from first principles in Go. Every component understood, nothing
copy-pasted. Maths picked up on the way, driven by need rather than study.

**Rule for the whole project:** write the slow, obviously-correct version first. Keep it.
Every fast version must be verified against it.

**Platform:** developed on Go 1.26 / linux/arm64, but intended to run competitively on
amd64 too. The core is portable scalar `uint64` — no PEXT, no GFNI, no SIMD. Any
architecture-specific path is an *optional accelerator* behind a build tag, must be a
drop-in replacement for a portable function that already exists and is already tested,
and must be justified by a benchmark. Realistically only NNUE inference will qualify.

---

## Phase 0 — Primitives

- [ ] Square indexing: 0–63, rank/file conversion, algebraic notation both ways
- [ ] `Bitboard uint64` type: set/clear/test a square
- [ ] Pretty-print a bitboard as an 8×8 grid (you will use this constantly)
- [ ] Learn `math/bits`: `OnesCount64`, `TrailingZeros64`, `LeadingZeros64`
- [ ] Iterate set bits: `for bb != 0 { sq := bits.TrailingZeros64(bb); bb &= bb - 1 }`
- [ ] Shift helpers with file-wrap masking (north, south, east, west, diagonals)

**Maths:** binary representation, sets as bit vectors, union/intersection/complement as
`|` `&` `^`. Why `bb & (bb-1)` clears the lowest set bit.

**Milestone:** print the starting position from bitboards.

---

## Phase 1 — Attacks

- [ ] Precomputed knight and king attack tables (`[64]Bitboard`)
- [ ] Pawn pushes, double pushes, captures, en passant targets
- [ ] `slowRayAttacks(sq, occupancy, isRook)` — a loop walking outward until blocked.
      Slow, obvious, correct. **Never delete this.**
- [ ] Relevant-occupancy masks: file+rank (or diagonals) minus the piece's own square,
      minus the last square of each ray
- [ ] `scatter(index, mask)` — spread an n-bit index across a mask's set bits
- [ ] Enumerate all 2^n occupancies for each square, compute true attacks with the slow version
- [ ] **Brute-force the magics**: random sparse candidates (`rand & rand & rand`),
      reject on destructive collision, keep what survives
- [ ] Magic lookup: `table[sq][((occ & mask[sq]) * magic[sq]) >> shift[sq]]`
- [ ] Assert magic output == slow output for all squares, all occupancies

**Maths:** combinatorics (2^n subsets), why edge squares carry no information, perfect
hashing, why multiplication is shift-and-add, constructive vs destructive collisions.

**Milestone:** generated your own magic numbers. Nothing copied.

---

## Phase 2 — Position & legality

- [ ] `Position` struct: piece bitboards, side to move, castling rights, ep square, clocks
- [ ] FEN parsing and generation
- [ ] `Make(move)` / `Unmake(move)` with an undo stack
- [ ] Pseudo-legal move generation for all pieces
- [ ] Attack detection: `isSquareAttackedBy(sq, side)`
- [ ] Legality: pinned pieces, evasions when in check, the en passant discovered-check edge case
- [ ] Castling with all its conditions

**Maths:** none really — this is where care beats cleverness. Encode moves compactly
(from, to, flags in 16 bits).

---

## Phase 3 — Perft

- [ ] `Perft(depth)` — count leaf nodes, nothing else
- [ ] Verify against published counts for the starting position to depth 6
- [ ] Verify against the standard test positions (Kiwipete, position 3, 4, 5, 6)
- [ ] `PerftDivide` — per-move breakdown, for bisecting bugs
- [ ] Make perft a CI-style test you can run in seconds

**Milestone:** all standard perft positions exact. Your move generation is now trustworthy
and you never doubt it again. Do not start Phase 4 before this passes.

---

## Phase 4 — Search

- [ ] Minimax to fixed depth (write it once, understand it, then throw it away)
- [ ] Negamax with alpha-beta
- [ ] Iterative deepening
- [ ] Zobrist hashing: random key per (piece, square), XOR in/out on make/unmake
- [ ] Transposition table: flat array, depth-preferred replacement, exact/lower/upper bounds
- [ ] Quiescence search — captures only, until the position is quiet
- [ ] Move ordering: TT move first, then MVV-LVA captures, then quiets
- [ ] Repetition and 50-move draw detection
- [ ] Mate scores that account for distance to mate

**Maths:** why perfect ordering gives b^(d/2) instead of b^d — the square root of the tree.
Zobrist as a group under XOR: every key is its own inverse, which is why unmake is free.

**Milestone:** finds mate in 3. Plays legal, non-embarrassing chess.

---

## Phase 5 — Evaluation

- [ ] Material counting
- [ ] Piece-square tables
- [ ] Tapered eval: interpolate midgame/endgame by remaining material
- [ ] Mobility (count of legal moves per piece)
- [ ] Pawn structure: passed, isolated, doubled, backward
- [ ] King safety: attackers in the king zone, pawn shield
- [ ] Bishop pair, rook on open file

**Maths:** linear interpolation. Recognise that a piece-square table *is* a one-layer
neural network — 768 binary inputs, one weighted sum. This matters later.

**Milestone:** ~1900–2200 Elo. A real engine.

---

## Phase 6 — Interface & testing

- [ ] UCI protocol: `uci`, `isready`, `position`, `go`, `stop`, `bestmove`
- [ ] Time management: allocate per move, respect increment, stop cleanly
- [ ] Run under a GUI (Cute Chess, Arena, Banksia)
- [ ] Automated matches against other engines via `cutechess-cli`
- [ ] **SPRT testing** — decide whether a change actually gained Elo
- [ ] An opening book or position set so games aren't all identical

**Maths:** statistics. Elo as a logistic model. Sequential probability ratio tests, why a
change that "looks better" over 100 games is usually noise.

**Milestone:** you can no longer fool yourself about whether a change helped. This phase
is what separates an engine from a toy.

---

## Phase 7 — Search refinement

Each of these is one idea, added and SPRT-tested individually. Never two at once.

- [ ] Killer moves and history heuristic
- [ ] Aspiration windows
- [ ] Principal variation search (null-window re-search)
- [ ] Null move pruning
- [ ] Late move reductions
- [ ] Futility pruning, reverse futility
- [ ] Static exchange evaluation (SEE) for capture ordering and pruning
- [ ] Check extensions
- [ ] Multithreading (Lazy SMP)

**Milestone:** ~2600–2800 Elo with a hand-crafted eval. This is roughly the ceiling
without a neural network.

---

## Phase 8 — Neural networks from scratch

Prerequisite for Phase 9. Do not skip.

- [ ] Write a 2→4→1 network in Go by hand. Forward pass only.
- [ ] Derive and implement backpropagation for it yourself, on paper first
- [ ] Train it on XOR — the simplest function a single layer *cannot* learn
- [ ] Watch 3Blue1Brown's neural network series
- [ ] Work through Karpathy's "Neural Networks: Zero to Hero" (micrograd)

**Maths:** derivatives as "if I nudge this, how much does the output move". The chain rule.
Gradient descent. Why a nonlinearity is required for depth to mean anything — two stacked
linear layers collapse into one matrix.

---

## Phase 9 — NNUE

- [ ] Data generation: self-play at shallow depth, label positions with search scores
- [ ] Start simple: 768 inputs (12 piece types × 64 squares) → 256 → 1, from both perspectives
- [ ] Train in PyTorch. Export weights as a flat binary file.
- [ ] Inference in Go, integer arithmetic, clipped ReLU
- [ ] **Verify Go inference matches PyTorch** on thousands of positions before trusting it
- [ ] Incremental accumulator: `acc -= w[oldFeature]; acc += w[newFeature]` on make/unmake
- [ ] SPRT against the hand-crafted eval
- [ ] Iterate: stronger engine → better data → stronger net → repeat
- [ ] Later: larger feature sets (king-bucketed, HalfKP-style), bigger layers

**Maths:** linear algebra as bulk arithmetic. Quantisation — why int8/int16 rather than
floats. Why the accumulator works: addition has inverses, `OR` does not.

**Milestone:** ~2900–3100 Elo. Genuinely strong.

---

## Phase 10 — The research idea

The reason for doing any of this differently from everyone else.

- [ ] Bit-sliced influence maps: store an attacker *count* per square, not a bit.
      Four `uint64` bit-planes, ripple-carry add/subtract, 64 counters updated at once.
- [ ] Incremental maintenance across make/unmake, including slider ray fixups when
      blockers appear or vanish
- [ ] Benchmark: `make + incremental update` vs `make + recompute`. **This one number
      decides whether the idea is real.**
- [ ] If it holds: derive mobility, king safety, and SEE directly from the counts
- [ ] Explore fusing the influence map with the NNUE accumulator — both are additive
      structures maintained across make/unmake. Can they be one object?

**Maths:** why moving from a Boolean lattice (`OR`, no inverses) to a `Z`-module
(`+`/`−`, invertible) is the whole point. Bit-slicing as a transpose of the data layout.

---

## Notes to self

- Slow correct version first, always. Fast version verified against it, always.
- One change at a time, SPRT-tested. Intuition about Elo is worthless.
- Perft before search. Search before eval. Eval before NNUE.
- Beating Stockfish is not the goal and never will be. Understanding every line is.
