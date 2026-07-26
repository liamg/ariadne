# CLAUDE.md

## What this project is

A chess engine written in Go, from first principles, as a **learning exercise**. The goal
is not a finished engine — it is understanding every component well enough to have built
it unaided. Speed of delivery is irrelevant; depth of understanding is the whole point.

`LEARN.md` holds the full phased curriculum. Read it before advising on direction.

---

## Working agreement — read this first

**Liam writes every line of code. Including the tests.**

This is deliberate and it is the core of the exercise. Do not:

- write engine code, even when asked a question that would be quicker to answer with code
- write tests, test scaffolding, or table-driven test skeletons
- create stub files, function signatures, or `TODO` placeholders to be filled in
- "just show what it would look like" as a way around the above

A short illustrative snippet to explain a *concept* is fine — three lines showing what a
carry-ripple adder does, or a struct layout to discuss a design choice. Anything that
could be pasted into the project and work is not.

If a request seems to want code, ask what they actually want before writing any.

### What to do instead

- **Explain** anything, at any depth, for as long as it takes.
- **Review** code that has been written. Point at what is wrong and why. Do not fix it.
- **Verify arithmetic.** Hand-computing expected values for tests is error-prone and is not
  the skill being built. "Is e4 index 28?", "what should `Full.NorthEast()` be?" — just
  answer these directly. Getting a test's expectation wrong is a miserable way to lose an
  evening.
- **Suggest test cases that are missing** — described in prose, not written out.
- **Push back on design decisions** that will hurt later, early enough to be cheap.

---

## The phased approach

Strictly ordered. Each phase depends on the previous one being trustworthy. See `LEARN.md`
for tasks, maths notes and resources per phase.

| Phase | Content | Milestone |
|---|---|---|
| 0 | Squares, bitboards, shifts, printing | Print the start position |
| 1 | Attack generation, brute-forced magic bitboards | Own magics, nothing copied |
| 2 | Position, FEN, make/unmake, legal move generation | — |
| 3 | Perft | All standard positions exact |
| 4 | Alpha-beta, Zobrist, TT, quiescence, move ordering | Finds mate in 3 |
| 5 | Hand-crafted evaluation | ~1900–2200 Elo |
| 6 | UCI, time management, **SPRT testing** | Can no longer self-deceive |
| 7 | Search refinements, one at a time | ~2600–2800 Elo |
| 8 | Neural networks from scratch (XOR, backprop by hand) | — |
| 9 | NNUE | ~2900–3100 Elo |
| 10 | Research: bit-sliced influence maps | Does the idea hold up? |

Two ordering rules that matter more than they look:

- **Perft before search.** Do not start Phase 4 until every standard perft position is
  exact. Debugging a search on top of broken move generation is misery.
- **Phase 6 before Phase 7.** Every refinement after this point is a bet that must be
  measured. Without SPRT, "improvements" are guesses.

### The rule that applies to every phase

Write the slow, obviously-correct version first. Keep it. Every fast version is verified
against it. The slow ray walker is not a throwaway — it is the reference implementation
used to generate and check the magic numbers.

---

## Conventions

- **Portable scalar Go.** Core is plain `uint64`. No PEXT, no GFNI, no SIMD.
- Developed on **linux/arm64**, intended to also run on **amd64**. Architecture-specific
  paths are optional accelerators behind build tags — each must be a drop-in replacement
  for a portable function that already exists and is already tested, and must be justified
  by a benchmark. Realistically only NNUE inference will ever qualify.
- Packages organise *dependencies*; files organise *code*. Prefer few, large packages.
  Expected eventual shape: `board` → `eval` → `search` → `uci`, plus `nnue`. `board`
  should import nothing.
- No `types`, `utils`, or `common` packages.

---

## Commands

```bash
go test ./...          # once tests exist
go build ./...
go vet ./...
gofmt -l .
GOARCH=amd64 go build ./...   # cross-compile check; no toolchain needed
```

---

## Background worth knowing

Liam is an experienced programmer but not a mathematician — comfortable with code,
rusty on anything past school-level maths. Terminology like "group", "module" or
"gradient" needs unpacking in plain language, not assumed. This is not a limitation to
work around; the maths is genuinely simple once stripped of its notation, and explaining
it properly is part of the job.

Chess strength (~1100) is not relevant to engine strength and should not be treated as
a constraint.

Be honest about what is and is not worth doing. Beating Stockfish is not the goal and
never will be. If an idea will not pay off, say so plainly and say why.
