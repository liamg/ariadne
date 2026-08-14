default: test

.PHONY: test test-full generate bench benchmark release-dryrun \
	harness-image harness-push sprt-engines sprt-smoke sprt-fast sprt

# ---- SPRT testing ----------------------------------------------------------

HARNESS_IMAGE ?= ghcr.io/liamg/ariadne-harness

# Content-addressed: the image is part of the frozen test config, so its tag
# must change when its inputs do, and must NOT be :latest.
HARNESS_TAG ?= $(shell cat testing/Dockerfile testing/entrypoint.sh testing/config.env \
	| sha256sum | cut -c1-12)

SPRT_BIN := $(CURDIR)/testing/bin
SPRT_OUT := $(CURDIR)/testing/pgn
SPRT_WT  := $(CURDIR)/testing/.base-worktree

# Compare against the commit this work branched from, not whatever main has
# moved on to - otherwise main shifting mid-run corrupts the comparison.
BASE ?= $(shell git merge-base HEAD main)

# One less than the CPU count, leaving something for the rest of the box. Each
# game holds two engine processes, but only one searches at a time - the other
# is blocked on stdin - so this is N busy threads, not 2N.
#
# Counts LOGICAL CPUs, so on an SMT box two searches share a physical core.
# That costs nps and adds timing jitter; it hits both engines equally, so it
# does not bias the result, but time losses INVALIDATE a run and lowering this
# is the fix. Also part of the run identity: change it and a saved run will
# start fresh rather than resume.
CONCURRENCY ?= $(shell n=$$(getconf _NPROCESSORS_ONLN 2>/dev/null || echo 2); \
	if [ "$$n" -gt 1 ]; then echo $$((n - 1)); else echo 1; fi)

# How many games a kill can cost. fastchess does not checkpoint on SIGTERM,
# so this is the real worst case for an interrupted run.
AUTOSAVE_GAMES ?= 20

STAMP = $(shell date +%Y%m%d-%H%M%S)

test:
	@echo "Running short test suite..."
	@go test -short -tags magicgen ./...

test-full:
	@echo "Running full test suite..."
	@go test -tags magicgen ./...

release-dryrun:
	@echo "Running release dry-run..."
	goreleaser release --snapshot --clean

generate:
	@echo "Generating code..."
	@go generate ./...

bench:
	@echo "Running benc..."
	go run . bench

benchmark:
	@echo "Running benchmarks..."
	go test -run=^$$ -bench=. -benchmem ./...

harness-image:
	@echo "Building harness image $(HARNESS_IMAGE):$(HARNESS_TAG)..."
	docker build --build-arg HARNESS_TAG=$(HARNESS_TAG) \
		-t $(HARNESS_IMAGE):$(HARNESS_TAG) testing

harness-push: harness-image
	docker push $(HARNESS_IMAGE):$(HARNESS_TAG)

# Builds dev from the working tree and base from the merge-base. Both bench
# numbers are printed: if they match, the change is non-functional and there is
# nothing for SPRT to measure.
sprt-engines:
	@mkdir -p $(SPRT_BIN) $(SPRT_OUT)
	@echo "Building dev from the working tree..."
	@go build -o $(SPRT_BIN)/dev .
	@echo "Building base from $(BASE)..."
	@rm -rf $(SPRT_WT)
	@git worktree add --detach --quiet $(SPRT_WT) $(BASE)
	@cd $(SPRT_WT) && go build -o $(SPRT_BIN)/base . \
		|| (echo "base does not build - no baseline available" >&2; exit 1)
	@git worktree remove --force $(SPRT_WT)
	@echo
	@dev=$$($(SPRT_BIN)/dev  bench | awk '/^Total nodes searched/ { print $$NF }'); \
	 base=$$($(SPRT_BIN)/base bench | awk '/^Total nodes searched/ { print $$NF }'); \
	 echo "dev  bench: $$dev"; \
	 echo "base bench: $$base"; \
	 if [ "$$dev" = "$$base" ]; then \
	   touch $(SPRT_BIN)/.bench-identical; \
	   echo; \
	   echo "NOTE: node counts are identical - the search explores the same tree."; \
	 else \
	   rm -f $(SPRT_BIN)/.bench-identical; \
	 fi
	@echo

# Identical node counts mean the search is unchanged. That is usually a no-op
# and testing it wastes days. The exception is a pure speedup: same tree, higher
# nps, which IS a strength change at a time control. So this blocks the
# expensive target rather than every target, and FORCE=1 overrides it.
define bench_guard
	@if [ -f $(SPRT_BIN)/.bench-identical ] && [ -z "$(FORCE)" ]; then \
		echo; \
		echo "REFUSING TO RUN: dev and base have identical bench node counts."; \
		echo; \
		echo "  The search is unchanged, so there is almost certainly nothing"; \
		echo "  to measure and this would burn days for a guaranteed null result."; \
		echo; \
		echo "  The one real exception is a pure speedup - identical tree, higher"; \
		echo "  nps. Compare the nps figures from the bench output above; if dev"; \
		echo "  is genuinely faster, rerun with FORCE=1."; \
		echo; \
		exit 1; \
	fi
endef

define run_harness
	docker run --rm \
		--user $$(id -u):$$(id -g) \
		-v $(SPRT_BIN):/engines:ro \
		-v $(SPRT_OUT):/out \
		-e CONCURRENCY=$(CONCURRENCY) \
		-e LABEL=$(1)-$(STAMP) \
		-e SIMPLIFY=$(if $(SIMPLIFY),$(SIMPLIFY),0) \
		-e ALLOW_WARNINGS=$(if $(ALLOW_WARNINGS),$(ALLOW_WARNINGS),0) \
		-e AUTOSAVE_GAMES=$(AUTOSAVE_GAMES) \
		$(if $(CPUSET),--cpuset-cpus=$(CPUSET),) \
		$(HARNESS_IMAGE):$(HARNESS_TAG) $(1)
endef

# Minutes. Fixed nodes, deterministic. Catches disasters, proves nothing good.
# Not guarded - it is cheap, and with identical binaries it is a valid harness
# self-check that should come back at exactly 50%.
sprt-smoke: sprt-engines
	$(call run_harness,smoke)

# Hours. "Did this make things seriously worse?"
sprt-fast: sprt-engines
	$(bench_guard)
	$(call run_harness,fast)

# Days on this box, hours on a rented one. The only target whose verdict counts.
# Pass SIMPLIFY=1 for non-regression bounds when removing something.
sprt: sprt-engines
	$(bench_guard)
	$(call run_harness,proper)
