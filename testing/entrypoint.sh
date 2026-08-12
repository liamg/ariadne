#!/usr/bin/env bash
#
# Ariadne test harness entrypoint.
#
#   entrypoint.sh smoke   deterministic fixed-node canary (~30 games, minutes)
#   entrypoint.sh fast    short-TC regression SPRT
#   entrypoint.sh proper  the real SPRT
#
# Engines are mounted at /engines/{base,dev}. Output lands in /out.
#
# Exit codes:
#   0   pass (or SPRT accepted H1)
#   2   harness integrity failure - timeouts or crashes. Result is INVALID.
#   3   smoke score below floor
#   4   SPRT concluded H0 (rejected)
#   5   SPRT inconclusive (hit the round cap)
#   64  usage error

set -euo pipefail

MODE="${1:-}"
case "${MODE}" in
    smoke | fast | proper) ;;
    *)
        echo "usage: $(basename "$0") {smoke|fast|proper}" >&2
        exit 64
        ;;
esac

# Overridable so the harness can be exercised outside the container.
CONFIG="${CONFIG:-/etc/ariadne/config.env}"
if [[ ! -r "${CONFIG}" ]]; then
    echo "config not readable: ${CONFIG}" >&2
    exit 64
fi
# shellcheck source=config.env
. "${CONFIG}"

BASE_ENGINE="${BASE_ENGINE:-/engines/base}"
DEV_ENGINE="${DEV_ENGINE:-/engines/dev}"
CONCURRENCY="${CONCURRENCY:-1}"
OUT_DIR="${OUT_DIR:-/out}"
LABEL="${LABEL:-${MODE}}"

for engine in "${BASE_ENGINE}" "${DEV_ENGINE}"; do
    if [[ ! -x "${engine}" ]]; then
        echo "engine is missing or not executable: ${engine}" >&2
        exit 64
    fi
done

if [[ ! -r "${BOOK}" ]]; then
    echo "opening book is missing: ${BOOK}" >&2
    exit 64
fi

mkdir -p "${OUT_DIR}"

# ---- assemble the match -----------------------------------------------------

each=(
    "option.Hash=${HASH_MB}"
    "option.Move Overhead=${MOVE_OVERHEAD_MS}"
)

sprt=()
case "${MODE}" in
    smoke)
        each+=("nodes=${SMOKE_NODES}")
        rounds="${SMOKE_ROUNDS}"
        ;;
    fast)
        each+=("tc=${FAST_TC}")
        rounds="${FAST_MAX_ROUNDS}"
        elo0="${FAST_ELO0}"
        elo1="${FAST_ELO1}"
        sprt=(-sprt "elo0=${elo0}" "elo1=${elo1}"
            "alpha=${ALPHA}" "beta=${BETA}" "model=${SPRT_MODEL}")
        ;;
    proper)
        each+=("tc=${PROPER_TC}")
        rounds="${PROPER_MAX_ROUNDS}"
        elo0="${PROPER_ELO0}"
        elo1="${PROPER_ELO1}"
        if [[ "${SIMPLIFY:-0}" == "1" ]]; then
            elo0="${SIMPLIFY_ELO0}"
            elo1="${SIMPLIFY_ELO1}"
        fi
        sprt=(-sprt "elo0=${elo0}" "elo1=${elo1}"
            "alpha=${ALPHA}" "beta=${BETA}" "model=${SPRT_MODEL}")
        ;;
esac

# ---- run identity and resume ------------------------------------------------
#
# fastchess autosaves tournament state and can resume from it, which matters
# when a proper run is measured in days. But it does NOT check that the engines
# match the saved state - verified directly: swapping the dev binary and
# resuming produced no warning whatsoever, silently blending two different
# engines' games into a single result that looks perfectly valid.
#
# So resumption is keyed on a hash of everything that defines the experiment.
# Change any part of it and you get a clean fresh run instead of a corrupt one.
#
# Two things to know about resumed runs:
#   - The SAVED settings win over the command line, so the round cap comes from
#     the original invocation, not from this one.
#   - Games played since the last checkpoint are replayed, so the pgn can hold
#     up to AUTOSAVE_GAMES-1 duplicate games. The statistics are unaffected -
#     they come from the saved state - but the pgn is not a clean record.

resume_args=()
resuming=0
append_outputs=false
state=""

if [[ "${MODE}" != "smoke" ]]; then
    run_id="$(
        printf '%s\n' \
            "$(cat /etc/ariadne/harness-tag 2>/dev/null || echo unknown)" \
            "${MODE}" \
            "${elo0}" "${elo1}" \
            "${CONCURRENCY}" \
            "$(sha256sum "${DEV_ENGINE}" | cut -d' ' -f1)" \
            "$(sha256sum "${BASE_ENGINE}" | cut -d' ' -f1)" |
            sha256sum | cut -c1-12
    )"

    # Stable across resumes, so a resumed run appends to the same pgn and logs
    # instead of scattering them over timestamped files.
    LABEL="${MODE}-${run_id}"

    state_dir="${OUT_DIR}/state"
    mkdir -p "${state_dir}"
    state="${state_dir}/${run_id}.json"

    if [[ -f "${state}" ]]; then
        resuming=1
        append_outputs=true
        resume_args=(-config "file=${state}" "outname=${state}")
    else
        # file= pointing at a missing path is a hard error, so a fresh run gets
        # outname only.
        resume_args=(-config "outname=${state}")
    fi

    # Checkpoint granularity: how much work a kill can cost. fastchess does NOT
    # checkpoint on SIGTERM, so this is the real worst case - 20 games is a few
    # minutes at proper time controls.
    resume_args+=(-autosaveinterval "${AUTOSAVE_GAMES:-20}")
fi

raw_log="${OUT_DIR}/${LABEL}.log"
log="${OUT_DIR}/${LABEL}.clean.log"
warn_log="${OUT_DIR}/${LABEL}.warn.log"
pgn="${OUT_DIR}/${LABEL}.pgn"

# A completed or invalidated tournament must never be resumed. Resuming a
# finished one returns instantly with a bogus "result"; resuming a tainted one
# carries the tainted games forward.
clear_state() {
    [[ -n "${state}" && -f "${state}" ]] && rm -f "${state}"
    return 0
}

# -strict makes fastchess exit non-zero on any warning, which catches UCI
# protocol regressions (bad mate signs, illegal moves, malformed info lines)
# that the integrity checks below would otherwise only notice after the fact.
# Enabled for the cheap modes only: aborting a 40-hour proper run on a single
# warning costs more than finishing it and rejecting the result afterwards,
# which the warning count below does anyway.
#
# ALLOW_WARNINGS=1 opts out. The case that needs it: base is an older commit
# that predates a UCI protocol fix, so the warnings come from the reference
# binary, not from the patch under test. Timeouts, crashes and stalls stay
# fatal regardless - those invalidate results, warnings merely indicate a bug
# somewhere.
ALLOW_WARNINGS="${ALLOW_WARNINGS:-0}"

strict=()
if [[ "${MODE}" != "proper" && "${ALLOW_WARNINGS}" != "1" ]]; then
    strict=(-strict)
fi

# Deliberately absent:
#   -recover      a crash must fail the run, not be silently papered over
#   -resign/-draw adjudication acts on engine-reported scores, so enabling it
#                 makes results depend on score reporting being trustworthy
args=(
    "${strict[@]}"
    -engine "cmd=${DEV_ENGINE}" name=dev
    -engine "cmd=${BASE_ENGINE}" name=base
    -each "${each[@]}"
    -openings "file=${BOOK}" format=epd order=random
    -srand "${SRAND}"
    -repeat
    -rounds "${rounds}"
    -games 2
    -concurrency "${CONCURRENCY}"
    -report penta=true
    -pgnout "file=${pgn}" "append=${append_outputs}"
    # Disconnects, stalls and time losses are logged here, NOT to stdout.
    # Without this the integrity checks below have nothing to read.
    -log "file=${warn_log}" level=warn "append=${append_outputs}"
    -event "ariadne-${MODE}"
    "${sprt[@]}"
    "${resume_args[@]}"
)

# Everything needed to identify what produced this result, so pasting the output
# into a PR carries its own provenance. Printed before the match AND again at
# the end: a proper run emits millions of lines, and a header only at the top
# would have scrolled away long before the result you actually paste.
print_run_header() {
    echo "harness     : $(cat /etc/ariadne/harness-tag 2>/dev/null || echo unknown)"
    echo "mode        : ${MODE}"
    echo "dev         : ${DEV_ENGINE}"
    echo "base        : ${BASE_ENGINE}"
    echo "book        : ${BOOK}"
    echo "hash        : ${HASH_MB}MB"
    echo "overhead    : ${MOVE_OVERHEAD_MS}ms"
    if ((${#sprt[@]})); then
        echo "bounds      : [${elo0}, ${elo1}] ${SPRT_MODEL}"
    fi
    echo "concurrency : ${CONCURRENCY}"
    echo "max games   : $((rounds * 2))"
    if [[ "${MODE}" != "smoke" ]]; then
        if ((resuming)); then
            echo "resuming    : yes - run ${run_id}, continuing saved tournament"
            echo "              (round cap comes from the ORIGINAL invocation)"
        else
            echo "resuming    : no - fresh run ${run_id}"
        fi
    fi
}

print_run_header
echo

# fastchess autosaves tournament state to config.json in its working directory.
# Run from the output directory so it lands with the pgn and logs instead of
# wherever the harness happened to be invoked from. All paths above are
# absolute, so this is safe.
cd "${OUT_DIR}"

set +e
fastchess "${args[@]}" 2>&1 | tee "${raw_log}"
fastchess_rc="${PIPESTATUS[0]}"
set -e

sed 's/\x1b\[[0-9;]*m//g' "${raw_log}" > "${log}"

# ---- integrity: this gates everything else ----------------------------------
#
# A single flagged game skews the result in a direction that has nothing to do
# with the patch. Silence about it is worse than a failure, so this is checked
# before the verdict and overrides it.

[[ -f "${warn_log}" ]] || : > "${warn_log}"

# Two sources, because neither is complete on its own. The end-of-run summary
# is only printed when the tournament finishes normally - an interrupted run
# skips it entirely - so the warn log is the reliable one.
summary_timeouts=$(awk '/^[[:space:]]*Timeouts:/ { s += $2 } END { print s + 0 }' "${log}")
summary_crashes=$(awk '/^[[:space:]]*Crashed:/  { s += $2 } END { print s + 0 }' "${log}")

logged_timeouts=$(grep -c "loses on time" "${warn_log}" || true)
logged_crashes=$(grep -c "disconnects" "${warn_log}" || true)
logged_stalls=$(grep -c "stalls" "${warn_log}" || true)

timeouts=$((summary_timeouts > logged_timeouts ? summary_timeouts : logged_timeouts))
crashes=$((summary_crashes > logged_crashes ? summary_crashes : logged_crashes))

warnings=$(grep -c "^Warning;" "${log}" || true)
interrupted=$(grep -c "Tournament was interrupted" "${log}" || true)

echo
echo "---- run ----"
print_run_header
echo
echo "---- integrity ----"
echo "timeouts    : ${timeouts}"
echo "crashes     : ${crashes}"
echo "stalls      : ${logged_stalls}"
echo "warnings    : ${warnings}$( ((warnings > 0)) && [[ "${ALLOW_WARNINGS}" == "1" ]] && echo "  (tolerated)")"
echo "interrupted : ${interrupted}"

fatal_warnings="${warnings}"
if [[ "${ALLOW_WARNINGS}" == "1" ]]; then
    fatal_warnings=0
    if ((warnings > 0)); then
        echo
        echo "NOTE: ${warnings} warning(s) tolerated via ALLOW_WARNINGS=1."
        echo "      Check ${warn_log} and satisfy yourself they come from base."
    fi
fi

if ((timeouts > 0 || crashes > 0 || logged_stalls > 0 || fatal_warnings > 0 || interrupted > 0)); then
    echo
    echo "INVALID: this run cannot be trusted. Discard the result." >&2
    ((timeouts > 0)) && echo "  - engines lost ${timeouts} game(s) on time" >&2
    ((crashes > 0)) && echo "  - engines disconnected ${crashes} time(s)" >&2
    ((logged_stalls > 0)) && echo "  - engines stalled ${logged_stalls} time(s)" >&2
    ((fatal_warnings > 0)) && echo "  - fastchess emitted ${warnings} warning(s)" >&2
    ((interrupted > 0)) && echo "  - the tournament was interrupted before finishing" >&2
    echo >&2
    echo "  logs: ${raw_log}" >&2
    echo "        ${warn_log}" >&2
    clear_state
    echo "  for time losses, raise Move Overhead or lower concurrency" >&2
    ((fatal_warnings > 0)) && echo "  if the warnings come from an older base that predates a" >&2
    ((fatal_warnings > 0)) && echo "  protocol fix, rerun with ALLOW_WARNINGS=1" >&2
    exit 2
fi

if ((fastchess_rc != 0)); then
    echo "fastchess exited ${fastchess_rc}" >&2
    exit 2
fi

# ---- verdict ----------------------------------------------------------------

score_pct=$(sed -n 's/.*Points: [0-9.]* (\([0-9.]*\) %).*/\1/p' "${log}" | tail -1)
games=$(sed -n 's/^Games: \([0-9]*\),.*/\1/p' "${log}" | tail -1)
verdict=$(grep -oE "SPRT \(.*\) completed - H[01] was accepted" "${log}" | tail -1 || true)

echo
echo "---- result ----"
echo "games : ${games:-0}"
echo "score : ${score_pct:-?} %"
# Only the final results block - fastchess prints an interim one every
# -ratinginterval games, and mixing them together is misleading.
awk '/^Results of /{ block = "" } { block = block $0 ORS } END { printf "%s", block }' "${log}" \
    | grep -E "^(Elo:|LLR:|Ptnml)" || true
echo "pgn   : ${pgn}"

if [[ "${MODE}" == "smoke" ]]; then
    # No statistics here - just "did it fall over or play like a disaster".
    if [[ -z "${score_pct}" ]]; then
        echo "could not parse a score from the match output" >&2
        exit 2
    fi
    if awk -v s="${score_pct}" -v f="${SMOKE_MIN_SCORE_PCT}" 'BEGIN { exit !(s < f) }'; then
        echo
        echo "FAIL: scored ${score_pct}%, floor is ${SMOKE_MIN_SCORE_PCT}%." >&2
        exit 3
    fi
    echo
    echo "PASS: no crashes, no time losses, score ${score_pct}% is above the ${SMOKE_MIN_SCORE_PCT}% floor."
    echo "This is a smoke test. It is NOT evidence the change is an improvement."
    exit 0
fi

case "${verdict}" in
    *"H1 was accepted")
        echo
        clear_state
        echo "SPRT ACCEPTED (H1)."
        exit 0
        ;;
    *"H0 was accepted")
        echo
        clear_state
        echo "SPRT REJECTED (H0)." >&2
        exit 4
        ;;
    *)
        echo
        echo "SPRT INCONCLUSIVE - hit the ${rounds}-round cap without crossing a bound."
        if [[ "${MODE}" == "fast" ]]; then
            # Failing to prove a regression is a pass for a regression check.
            clear_state
            echo "No regression proven, which is all this mode can tell you."
            exit 0
        fi
        clear_state
        exit 5
        ;;
esac
