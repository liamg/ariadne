package uci

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func init() {
	register("setoption", func(ctx context.Context, e *engine, r responder, input string) error {
		if !strings.HasPrefix(input, "name ") {
			r.debugf("setoption: missing name prefix: %s", input)
			return nil
		}

		input = strings.TrimSpace(input[4:])

		var value string
		name, after, ok := strings.Cut(input, " value ")
		if ok {
			value = decodeEmpty(strings.TrimSpace(after))
		}

		for _, option := range options {
			if option.name == name {
				switch option.optionType {
				case optionTypeSpin:
					// clamp between min and max
					val, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						r.debugf("setoption: invalid spin value: %s", value)
						return nil
					}
					val = max(option.min, min(option.max, val))
					value = fmt.Sprintf("%d", val)
				case optionTypeCombo:
					// check if value is in options
					if !slices.Contains(option.options, value) {
						r.debugf("setoption: invalid combo value: %s", value)
						return nil
					}
				}
				option.apply(e, r, value)
				return nil
			}
		}

		r.debugf("setoption: unknown option name: %s", name)
		return nil
	})
}
