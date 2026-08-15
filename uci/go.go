package uci

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/liamg/ariadne/search"
)

func init() {
	register("go", func(ctx context.Context, e *engine, r responder, input string) error {
		limits := search.Limits{
			Depth:          0, // default depth to maximum
			MoveOverheadMS: int64(e.options.moveOverheadMS),
		}

		fields := strings.Fields(input)

		for i := 0; i < len(fields); i++ {
			switch fields[i] {
			case "movestogo":
				if i+1 < len(fields) {
					i++
					if t, err := strconv.ParseInt(fields[i], 10, 64); err == nil && t > 0 {
						limits.MovesBeforeControlSwitch = t
					}
				}
			case "winc":
				if i+1 < len(fields) {
					i++
					if t, err := strconv.ParseInt(fields[i], 10, 64); err == nil && t > 0 {
						limits.WhiteIncrementMS = t
					}
				}
			case "binc":
				if i+1 < len(fields) {
					i++
					if t, err := strconv.ParseInt(fields[i], 10, 64); err == nil && t > 0 {
						limits.BlackIncrementMS = t
					}
				}
			case "movetime":
				if i+1 < len(fields) {
					i++
					if t, err := strconv.ParseInt(fields[i], 10, 64); err == nil && t > 0 {
						limits.MoveTimeMS = t
					}
				}
			case "wtime":
				if i+1 < len(fields) {
					i++
					if t, err := strconv.ParseInt(fields[i], 10, 64); err == nil && t >= 0 {
						limits.WhiteTimeMS = t
						limits.HasTimeControl = true
					}
				}
			case "btime":
				if i+1 < len(fields) {
					i++
					if t, err := strconv.ParseInt(fields[i], 10, 64); err == nil && t >= 0 {
						limits.BlackTimeMS = t
						limits.HasTimeControl = true
					}
				}
			case "depth":
				if i+1 < len(fields) {
					i++
					if d, err := strconv.Atoi(fields[i]); err == nil && d > 0 {
						limits.Depth = d
					}
				}
			case "nodes":
				if i+1 < len(fields) {
					i++
					if n, err := strconv.ParseInt(fields[i], 10, 64); err == nil && n > 0 {
						limits.Nodes = n
					}
				}
			case "ponder":
				limits.Ponder = true
			case "mate":
				if i+1 < len(fields) {
					i++
					if m, err := strconv.Atoi(fields[i]); err == nil && m > 0 {
						limits.Depth = (m * 2) - 1
					}
				}
			case "infinite":
				limits.Depth = 0 // go deep to find mate
				limits.Infinite = true
			default:
				r.debugf("unknown go option: %s", strings.Join(fields[i:], " "))
			}
		}

		go func() {
			result, err := e.goSearch(ctx, limits)
			if err != nil {
				r.debugf("error during search: %v", err)
				if errors.Is(err, errAlreadySearching) {
					return
				}
				// otherwise fall through and emit a null move
			}

			if len(result.PV) > 1 {
				r.sendf("bestmove %s ponder %s", result.BestMove, result.PV[1])
			} else {
				r.sendf("bestmove %s", result.BestMove)
			}
		}()

		return nil
	})
}
