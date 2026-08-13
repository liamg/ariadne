# Ariadne test harness.
#
# Contains fastchess, a reference Stockfish, and the opening book. It does NOT
# contain Ariadne - base and dev binaries are built per run and mounted at
# /engines, because they change every run and baking them in would mean
# rebuilding the image constantly.
#
# The image is part of the frozen test config. Rebuilding it with a different
# fastchess version or a different book changes what results mean, so it is
# tagged by content rather than :latest.

# ---- build fastchess from a pinned tag --------------------------------------

FROM ubuntu:24.04 AS fastchess

ARG FASTCHESS_VERSION=v1.8.2-alpha

RUN apt-get update && apt-get install -y --no-install-recommends \
        build-essential \
        ca-certificates \
        git \
    && rm -rf /var/lib/apt/lists/*

# NATIVE= overrides fastchess's default -march=native. Without it the binary is
# tuned to whichever machine built the image, and GitHub's runner fleet is
# heterogeneous - a build on an AVX-512 Intel runner dies with SIGILL (exit 132)
# the moment it runs on an AMD one. The failure is a CPU lottery, so it looks
# like flakiness rather than a build bug.
#
# Costs nothing: fastchess is a referee shuttling UCI text between two
# processes, not a hot vector loop. Note that fastchess's own `build=release`
# is NOT the fix - it hardcodes -march=x86-64 and -static, breaking arm64.
RUN git clone --depth 1 --branch "${FASTCHESS_VERSION}" \
        https://github.com/Disservin/fastchess.git /src \
    && make -C /src NATIVE= -j"$(nproc)" \
    && /src/fastchess --version

# ---- opening book -----------------------------------------------------------
#
# Provisional. This is Stockfish's UHO_Lichess_4852_v1 (2,632,036 positions,
# 175MB unpacked), pinned to a commit and checksummed so the build is
# reproducible. We keep a deterministic 20,000-position subset: that is enough
# for a 40,000-game run at two games per opening with no reuse, and it keeps
# the image small.
#
# Replace this stage with a COPY of the generated CC0 book when it exists.
# Swapping the book is a frozen-config change.

FROM ubuntu:24.04 AS book

ARG BOOK_COMMIT=426eca422c202ef381a1380ada1873a7f185c4b9
ARG BOOK_SHA256=4e298f11e8acfa106babe02968f2e61582145e7874c59284690b20b9650e0e07
ARG BOOK_POSITIONS=20000

RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
        curl \
        unzip \
    && rm -rf /var/lib/apt/lists/*

RUN curl -fsSL -o /tmp/book.zip \
        "https://raw.githubusercontent.com/official-stockfish/books/${BOOK_COMMIT}/UHO_Lichess_4852_v1.epd.zip" \
    && echo "${BOOK_SHA256}  /tmp/book.zip" | sha256sum --check --strict

# Fixed stride rather than a shuffle, so the subset is identical on every build.
RUN mkdir -p /books \
    && total="$(unzip -p /tmp/book.zip | wc -l)" \
    && stride="$((total / BOOK_POSITIONS))" \
    && unzip -p /tmp/book.zip \
        | awk -v s="${stride}" 'NR % s == 1' \
        | head -n "${BOOK_POSITIONS}" > /books/book.epd \
    && rm /tmp/book.zip \
    && lines="$(wc -l < /books/book.epd)" \
    && [ "${lines}" -eq "${BOOK_POSITIONS}" ] \
        || (echo "expected ${BOOK_POSITIONS} positions, got ${lines}" >&2; exit 1)

# ---- runtime ----------------------------------------------------------------

FROM ubuntu:24.04

# Stockfish is a stable reference opponent and the evaluator used to generate
# the book. Pinning matters: an opponent that silently changes strength
# invalidates every cross-time comparison. Record the version in the image.
#
# Debian installs the binary to /usr/games, which is NOT on the default
# container PATH, so it gets symlinked somewhere findable.
RUN apt-get update && apt-get install -y --no-install-recommends \
        stockfish \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && dpkg-query --showformat='${Version}' --show stockfish > /etc/stockfish-version \
    && echo "stockfish $(cat /etc/stockfish-version)" \
    && ln -sf /usr/games/stockfish /usr/local/bin/stockfish

# Fail the build, not a test run three weeks from now, if the NNUE net is
# missing or broken.
#
# 2>&1 is required: stockfish writes its bench summary to stderr, so piping
# stdout alone finds nothing even though the line is plainly visible in a
# terminal (which shows both streams).
RUN printf 'bench\nquit\n' | stockfish 2>&1 | grep -q "Nodes searched" \
    || (echo "stockfish is installed but cannot run bench" >&2; exit 1)

COPY --from=fastchess /src/fastchess /usr/local/bin/fastchess

COPY --from=book /books/book.epd /books/book.epd

# Stamped into the image so every run's output states which harness produced it.
# Provenance travels with the result instead of living in a file someone has to
# remember to update.
ARG HARNESS_TAG=unknown
RUN mkdir -p /etc/ariadne && echo "${HARNESS_TAG}" > /etc/ariadne/harness-tag

COPY config.env /etc/ariadne/config.env
COPY entrypoint.sh /usr/local/bin/entrypoint.sh
RUN chmod +x /usr/local/bin/entrypoint.sh

# Mounted per run.
VOLUME ["/engines", "/out"]

WORKDIR /out

LABEL org.opencontainers.image.source="https://github.com/liamg/ariadne"
LABEL org.opencontainers.image.description="fastchess + stockfish + book, for Ariadne SPRT testing"

ENTRYPOINT ["/usr/local/bin/entrypoint.sh"]
