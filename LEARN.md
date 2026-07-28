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

**On the maths:** none of this needs a degree. Each phase explains what you actually need
in plain language. Where a term sounds intimidating, it's usually one idea in a bad suit.

---

## Standing resources

Bookmark these now; they're referenced throughout.

- **[Chess Programming Wiki](https://www.chessprogramming.org/)** — the canonical reference
  for everything in this document. If a term confuses you, it has a page.
- **[TalkChess forum](https://talkchess.com/)** — where engine authors actually talk. Search
  before asking; decades of answers are already there.
- **[Stockfish source](https://github.com/official-stockfish/Stockfish)** — reference
  implementation of essentially every idea here. Read it *after* trying yourself.
- **[Blunder](https://github.com/algerbrex/blunder)** — a strong engine written in Go, with
  an author's blog series explaining the decisions. Closest thing to a peer project.
- **Sebastian Lague — "Coding Adventure: Chess"** (YouTube) — the friendliest visual
  introduction to the whole problem. Watch before starting.

---

## Phase 0 — Primitives

- [x] Square indexing: 0–63, rank/file conversion, algebraic notation both ways
- [x] `Bitboard uint64` type: set/clear/test a square
- [x] Pretty-print a bitboard as an 8×8 grid (you will use this constantly)
- [x] Learn `math/bits`: `OnesCount64`, `TrailingZeros64`, `LeadingZeros64`
- [x] Iterate set bits: `for bb != 0 { sq := bits.TrailingZeros64(bb); bb &= bb - 1 }`
- [x] Shift helpers with file-wrap masking (north, south, east, west, diagonals)

**Maths — sets as bits.** A 64-bit integer is 64 yes/no answers. Treat it as "which squares
are in this set" and the operators become set theory: `|` is union, `&` is intersection,
`^` is symmetric difference, `&^` is set subtraction. "White pawns that can capture
something" is one `&`. That's the entire appeal.

The one trick worth understanding rather than memorising: `bb & (bb - 1)` clears the
lowest set bit. Subtracting 1 flips the lowest 1 to 0 and turns every 0 below it into 1;
`&`-ing with the original keeps only the bits above. Write out `0b1011000 - 1` on paper
once and it's obvious forever.

**Resources**

- [Bitboards](https://www.chessprogramming.org/Bitboards) — CPW overview
- [General Setwise Operations](https://www.chessprogramming.org/General_Setwise_Operations) —
  the complete catalogue of bit tricks. Skim now, return often.
- [Go `math/bits` docs](https://pkg.go.dev/math/bits)

**Milestone:** print the starting position from bitboards.

---

## Phase 1 — Attacks

- [x] Precomputed knight and king attack tables (`[64]Bitboard`)
- [x] Pawn pushes, double pushes, captures, en passant targets
- [x] `slowRayAttacks(sq, occupancy, isRook)` — a loop walking outward until blocked.
      Slow, obvious, correct. **Never delete this.**
- [x] Relevant-occupancy masks: file+rank (or diagonals) minus the piece's own square,
      minus the last square of each ray
- [x] `scatter(index, mask)` — spread an n-bit index across a mask's set bits
- [x] Enumerate all 2^n occupancies for each square, compute true attacks with the slow version
- [x] **Brute-force the magics**: random sparse candidates (`rand & rand & rand`),
      reject on destructive collision, keep what survives
- [x] Magic lookup: `table[sq][((occ & mask[sq]) * magic[sq]) >> shift[sq]]`
- [x] Assert magic output == slow output for all squares, all occupancies

**Maths — counting possibilities.** If 10 squares can each be occupied or empty
independently, there are 2×2×2… ten times = 2¹⁰ = 1024 possible arrangements. That's all
"combinatorics" means here: multiply the choices. Since there are only 1024 inputs, you
precompute all 1024 answers. Not clever — just exhaustive.

**Maths — perfect hashing.** You need to turn 10 bits scattered across a 64-bit word into a
dense number 0–1023 so it can index an array. A "hash function" is anything that maps big
inputs to small numbers; a *perfect* one produces no harmful collisions. Multiplication
works because multiplying by a constant is just adding shifted copies of the input — pick
the right constant and your scattered bits land where you want them. Nobody can derive
that constant, so you guess randomly until one works.

**Resources**

- [Magic Bitboards](https://www.chessprogramming.org/Magic_Bitboards) — the reference page
- [Looking for Magics](https://www.chessprogramming.org/Looking_for_Magics) — Tord Romstad's
  original search method, including the `rand & rand & rand` trick
- [Sliding Piece Attacks](https://www.chessprogramming.org/Sliding_Piece_Attacks) — the
  alternatives (kindergarten, hyperbola quintessence) if you get curious

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

**Maths:** none. This phase is pure care. Encode moves compactly — from (6 bits), to
(6 bits), flags (4 bits) fits in a `uint16`.

**Resources**

- [Forsyth-Edwards Notation](https://www.chessprogramming.org/Forsyth-Edwards_Notation) — FEN spec
- [Legal Move](https://www.chessprogramming.org/Legal_Move) and
  [Pin](https://www.chessprogramming.org/Pin) — the cases that bite
- [Encoding Moves](https://www.chessprogramming.org/Encoding_Moves)

---

## Phase 3 — Perft

- [ ] `Perft(depth)` — count leaf nodes, nothing else
- [ ] Verify against published counts for the starting position to depth 6
- [ ] Verify against the standard test positions (Kiwipete, position 3, 4, 5, 6)
- [ ] `PerftDivide` — per-move breakdown, for bisecting bugs
- [ ] Make perft a CI-style test you can run in seconds

**How to debug a perft mismatch:** run `PerftDivide` on your engine and on a known-good one
(Stockfish has `go perft N`). Compare the per-move counts, find the move whose subtree count
differs, play that move, repeat one level down. You'll land on the exact broken position in
a few steps. This is a binary search and it never fails.

**Resources**

- [Perft](https://www.chessprogramming.org/Perft) and
  [Perft Results](https://www.chessprogramming.org/Perft_Results) — the standard positions
  and their exact node counts. These numbers are your test suite.

**Milestone:** all standard perft positions exact. Your move generation is now trustworthy
and you never doubt it again. **Do not start Phase 4 before this passes.**

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

**Maths — why move ordering dominates.** With ~35 legal moves per position, searching 8
moves deep means 35⁸ ≈ 2.2 trillion positions. Alpha-beta with *perfect* move ordering cuts
that to 35⁴ ≈ 1.5 million — the square root. With random ordering you get roughly 35⁶.

The practical reading: ordering isn't a 10% optimisation, it's the difference between depth
8 and depth 4. Almost everything in Phase 7 exists to improve ordering.

**Maths — XOR is its own undo.** Zobrist hashing gives every (piece, square) pair a random
64-bit number and XORs them together to identify a position. It works because `a ^ b ^ b = a`
— XORing twice cancels out. So making a move is "XOR out the piece's old square, XOR in the
new one", and unmaking is *the identical operation*. No undo log needed. (Formally this is a
group where every element is its own inverse; you don't need the vocabulary, just the fact.)

**Resources**

- [Alpha-Beta](https://www.chessprogramming.org/Alpha-Beta) and
  [Negamax](https://www.chessprogramming.org/Negamax)
- [Quiescence Search](https://www.chessprogramming.org/Quiescence_Search) — skip this and
  your engine hallucinates constantly
- [Zobrist Hashing](https://www.chessprogramming.org/Zobrist_Hashing)
- [Transposition Table](https://www.chessprogramming.org/Transposition_Table) — read the
  section on bound types (exact/lower/upper) twice; it's the usual source of subtle bugs
- [MVV-LVA](https://www.chessprogramming.org/MVV-LVA)

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

**Maths — linear interpolation.** A knight on the rim is bad in the middlegame and fine in
an endgame, so you keep two values and blend them: `score = (mg*phase + eg*(24-phase)) / 24`,
where `phase` counts remaining material. That's all "tapered eval" means — a weighted
average that slides as pieces come off.

**Worth noticing:** a piece-square table is 768 numbers (12 piece types × 64 squares), and
evaluating it is "add up the values for the pieces that are present". That is *exactly* a
one-layer neural network with 768 binary inputs. Keep this in mind — Phase 9 is the same
idea with a hidden layer bolted on.

**Resources**

- [Simplified Evaluation Function](https://www.chessprogramming.org/Simplified_Evaluation_Function)
  — a complete, sane starting PSQT set with values you can type in today
- [Tapered Eval](https://www.chessprogramming.org/Tapered_Eval)
- [Evaluation](https://www.chessprogramming.org/Evaluation) — the full menu of terms

**Milestone:** ~1900–2200 Elo. A real engine.

---

## Phase 6 — Interface & testing

**This is the phase that separates an engine from a toy.** Everything after it depends on
being able to answer "did that change help?" honestly.

- [ ] UCI protocol: `uci`, `isready`, `position`, `go`, `stop`, `bestmove`
- [ ] Time management: allocate per move, respect increment, stop cleanly
- [ ] Run under a GUI (Cute Chess, Arena, Banksia)
- [ ] Get an opening book / EPD position set so games aren't all identical
- [ ] Automated engine-vs-engine matches with `fastchess`
- [ ] **Set up SPRT and use it for every single change from here on**

### The maths you need here, properly

**Standard deviation — how much things wobble.** Flip a fair coin 100 times and you expect
50 heads, but you'll actually get something like 44–56. Standard deviation is the number
that quantifies that spread. The one fact that matters:

> The wobble shrinks with the **square root** of the sample size.
> Four times as many games = half the noise. Sixteen times = a quarter.

That square root is brutal, and it's why measuring small improvements is so expensive. To
halve your uncertainty you must quadruple your work.

**Elo — a scale, not a score.** Elo converts a rating gap into an expected result:
`score = 1 / (1 + 10^(-gap/400))`. A 400-point gap means the stronger side scores ~91%. A
**5-point** gap — a genuinely worthwhile engine improvement — means scoring **50.72%**
instead of 50.00%.

Now combine those two facts. You're trying to detect a 0.72 percentage-point difference,
and noise only falls as √N. That's why it takes **20,000–40,000 games**. A change that wins
55 out of 100 games has told you *nothing* — that's comfortably inside noise.

**SPRT — Sequential Probability Ratio Test.** A statistical test (Abraham Wald, 1945) that
checks the evidence after *every game* and stops the moment the answer is clear, rather than
committing to a fixed number of games up front.

You state two hypotheses and two error rates:

```
H0: the change is worth   0 Elo   (useless)
H1: the change is worth  +5 Elo   (worth keeping)
alpha = beta = 0.05               (5% chance of each kind of mistake)
```

After each game it asks "how much more likely are these results under H1 than under H0?",
accumulates that as a running total, and stops when the total crosses an upper bound
(accept) or lower bound (reject). It's a random walk with two exits.

The payoff is early stopping: a patch that loses 30 Elo gets rejected after ~300 games
instead of 40,000. Only genuinely marginal changes cost you the full budget. In practice
this more than halves your total compute.

```bash
fastchess \
  -engine cmd=./engine-new name=new \
  -engine cmd=./engine-old name=old \
  -each tc=10+0.1 -concurrency 8 \
  -openings file=book.epd format=epd order=random \
  -repeat -games 80000 \
  -sprt elo0=0 elo1=5 alpha=0.05 beta=0.05
```

Two flags that aren't optional. **`-openings`**: engines are deterministic, so without varied
starting positions you'd play the same game 40,000 times. **`-repeat`**: each opening is
played twice with colours swapped, so a lucky opening doesn't bias the result.

**The discipline this forces:** one change at a time. Bundle three changes and a passing
SPRT tells you the *bundle* is net-positive — not which parts helped and which were
dragging it down. This is also the real cost of engine development: tens of thousands of
games per change, forever. Budget for it.

**Resources**

- [fastchess](https://github.com/Disservin/fastchess) — the match runner (replaced
  cutechess-cli; also see its [man page](https://github.com/Disservin/fastchess/blob/master/man.md))
- [Sequential Probability Ratio Test](https://www.chessprogramming.org/Sequential_Probability_Ratio_Test)
  — CPW's chess-specific treatment
- [SPRT on Wikipedia](https://en.wikipedia.org/wiki/Sequential_probability_ratio_test) —
  the general statistics, if you want the derivation
- [UCI protocol spec](https://www.chessprogramming.org/UCI)
- [Fishtest](https://tests.stockfishchess.org/) — Stockfish's distributed testing framework.
  Watch real patches pass and fail; excellent intuition for how noisy this all is.
- [CCRL rating lists](https://computerchess.org.uk/ccrl/) — where engines get ranked, once
  yours is worth submitting

**Milestone:** you can no longer fool yourself about whether a change helped.

---

## Phase 7 — Search refinement

Each of these is one idea, added and SPRT-tested individually. **Never two at once.**

- [ ] Killer moves and history heuristic
- [ ] Aspiration windows
- [ ] Principal variation search (null-window re-search)
- [ ] Null move pruning
- [ ] Late move reductions
- [ ] Futility pruning, reverse futility
- [ ] Static exchange evaluation (SEE) for capture ordering and pruning
- [ ] Check extensions
- [ ] Multithreading (Lazy SMP)

**The theme:** almost all of these are bets. "This move is probably bad, so search it
shallowly (or not at all)." They're unsound in theory — you *will* miss things — but they
win far more from extra depth than they lose from occasional blindness. Which is exactly
why every one of them needs SPRT rather than reasoning.

**Resources**

- [Null Move Pruning](https://www.chessprogramming.org/Null_Move_Pruning)
- [Late Move Reductions](https://www.chessprogramming.org/Late_Move_Reductions)
- [Principal Variation Search](https://www.chessprogramming.org/Principal_Variation_Search)
- [Static Exchange Evaluation](https://www.chessprogramming.org/Static_Exchange_Evaluation)
- [Lazy SMP](https://www.chessprogramming.org/Lazy_SMP)

**Milestone:** ~2600–2800 Elo. Roughly the ceiling without a neural network.

---

## Phase 8 — Neural networks from scratch

Prerequisite for Phase 9. **Do not skip.** Build the intuition on a toy before touching
10 million weights.

- [ ] Write a 2→4→1 network in Go by hand. Forward pass only.
- [ ] Derive and implement backpropagation for it yourself, on paper first
- [ ] Train it on XOR — the simplest function a single layer *cannot* learn
- [ ] Watch the 3Blue1Brown series
- [ ] Work through Karpathy's micrograd

### The maths, in plain language

**A neuron** is: multiply each input by a weight, add them up, add a bias, then apply one
bending function. That's it. Material counting (`1×pawns + 3×knights + 5×rooks`) is a
neuron. Your piece-square table is a layer of them.

**Why a "nonlinearity" is required.** If layer 1 multiplies by A and layer 2 multiplies by B,
chaining them gives B×A — which is just *another single multiplication*. Two stacked linear
layers collapse into one, so depth would buy you nothing. Inserting a bend between them —
`clamp(x, 0, 127)` — stops the collapse. That kink is the entire reason deep networks work.

**A derivative** is the answer to: *"if I nudge this number up slightly, does the output go
up or down, and by how much?"* Nothing more. Concretely:

```
w = 0.500  →  error 10.00
w = 0.501  →  error  9.98     raising w lowered the error
                              slope ≈ (9.98 - 10.00) / 0.001 = -20
```

**Gradient descent** is then obvious: compute that slope for every weight, and nudge each
one in the direction that reduces error. Repeat a few hundred million times. The "gradient"
is just the collection of all those slopes at once.

**Backpropagation** is the trick for computing all of them in one backward sweep instead of
testing each weight individually. It's the chain rule from calculus, applied efficiently.

**You will never write this in production.** PyTorch computes gradients automatically. You do
it once by hand on the toy network so it stops being magic, then let the library do it forever.

**Resources**

- [3Blue1Brown — Neural Networks](https://www.3blue1brown.com/topics/neural-networks) —
  four videos, visual, the best intuition-builder that exists. Start here.
- [Karpathy — Neural Networks: Zero to Hero](https://karpathy.ai/zero-to-hero.html) —
  builds an autograd engine from scratch in ~100 lines of Python. Written for programmers.
  After this, PyTorch is no longer mysterious.

---

## Phase 9 — NNUE

- [ ] Read the NNUE document below. Twice. It's long and it's worth it.
- [ ] Data generation: self-play at shallow depth, label positions with search scores
- [ ] Start simple: 768 inputs (12 piece types × 64 squares) → 256 → 1, from both perspectives
- [ ] Train in PyTorch. Export weights as a flat binary file.
- [ ] Inference in Go, integer arithmetic, clipped ReLU
- [ ] **Verify Go inference matches PyTorch** on thousands of positions before trusting it
- [ ] Incremental accumulator: `acc -= w[oldFeature]; acc += w[newFeature]` on make/unmake
- [ ] SPRT against the hand-crafted eval
- [ ] Iterate: stronger engine → better data → stronger net → repeat
- [ ] Later: larger feature sets (king-bucketed, HalfKP-style), bigger layers

**The one idea.** Inputs are binary and sparse (~32 of 768 are 1), so the first layer needs
no multiplication at all — just add the weight-columns of the active features. And when a
piece moves, exactly one feature turns off and one turns on:

```go
acc.sub(weights[oldFeature])   // two vector operations
acc.add(weights[newFeature])   // instead of thirty
```

You never recompute the big layer; you *maintain* it across make/unmake. That's the
"efficiently updatable" in the name, and it's the whole invention.

**Maths — quantisation.** Weights are stored as int8/int16 rather than floats, and
activations clamp to 0–127. Integers are faster and the precision loss costs almost nothing.
This is also why SIMD matters here and nowhere else in the engine.

**Watch out for:** validation loss is not playing strength. A net that predicts labels more
accurately can easily produce a *weaker* engine. Only SPRT decides.

**Resources**

- [NNUE — the official document](https://github.com/official-stockfish/nnue-pytorch/blob/master/docs/nnue.md)
  ([readable HTML version](https://official-stockfish.github.io/docs/nnue-pytorch-wiki/docs/nnue.html))
  — comprehensive, from architecture through quantisation to SIMD implementation. The single
  best resource that exists on this topic.
- [nnue-pytorch](https://github.com/official-stockfish/nnue-pytorch) — the trainer itself
- [NNUE on CPW](https://www.chessprogramming.org/NNUE) — shorter overview first, if the
  document above is too much on a first pass

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

**The idea in one paragraph.** Bitboards answer "is this square attacked?" with one bit, and
you build the answer with `OR`. But `OR` throws away information: if a knight and a bishop
both attack e5 and the knight leaves, you can't remove just the knight's contribution —
you have to recompute. Engines do this millions of times a second. Store a *count* instead
and removal becomes subtraction, which is perfectly reversible.

The formal way to say that: `OR` gives you a lattice (no inverses), addition gives you a
module over the integers (inverses exist). You don't need the vocabulary — the fact is the
whole point, and it's the same insight that makes the NNUE accumulator work.

**Bit-slicing** is the trick that keeps it fast: rather than 64 separate 4-bit counters,
store all the bit-0s in one `uint64`, all the bit-1s in another, and so on. Then a
ripple-carry adder in ~8 instructions increments all 64 counters at once.

**Resources**

- Nothing to link. That's the point.

---

## Notes to self

- Slow correct version first, always. Fast version verified against it, always.
- One change at a time, SPRT-tested. Intuition about Elo is worthless.
- Perft before search. Search before eval. Eval before NNUE.
- When stuck, the Chess Programming Wiki almost certainly has the page.
- Beating Stockfish is not the goal and never will be. Understanding every line is.
