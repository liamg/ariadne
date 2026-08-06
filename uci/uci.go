package uci

import (
	"context"
	"fmt"
	"strings"

	"github.com/liamg/ariadne/version"
)

func init() {
	register("uci", func(ctx context.Context, e *engine, r responder, input string) error {
		// engine details
		r.sendf("id name %s %s", version.Name, version.Version)
		r.sendf("id author %s", version.Author)

		// send available options
		for _, opt := range options {
			switch opt.optionType {
			case optionTypeSpin:
				r.sendf("option name %s type spin default %d min %d max %d", opt.name, opt.defaultValue, opt.min, opt.max)
			case optionTypeCheck:
				r.sendf("option name %s type check default %t", opt.name, opt.defaultValue)
			case optionTypeCombo:
				vars := []string{}
				for _, v := range opt.options {
					vars = append(vars, encodeEmpty(v))
				}
				r.sendf("option name %s type combo default %s %s", opt.name, opt.defaultValue, strings.Join(vars, " var "))
			case optionTypeString:
				r.sendf("option name %s type string default %s", opt.name, encodeEmpty(opt.defaultValue))
			case optionTypeButton:
				r.sendf("option name %s type button", opt.name)
			default:
				r.debugf("unknown option type: %s", opt.optionType)
			}
		}

		// done
		r.send("uciok")
		return nil
	})
}

const empty = "<empty>"

func encodeEmpty(input any) string {
	s := fmt.Sprintf("%s", input)
	if s == "" {
		return empty
	}
	return s
}

func decodeEmpty(input string) string {
	if input == empty {
		return ""
	}
	return input
}
