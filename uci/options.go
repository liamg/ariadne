package uci

import "strconv"

type engineOptions struct {
	debug                      bool
	transpositionTableMemoryMB int
	moveOverheadMS             int
}

var defaultEngineOptions = engineOptions{
	debug:                      false,
	transpositionTableMemoryMB: 64,
	moveOverheadMS:             50,
}

type optionType string

const (
	optionTypeSpin   optionType = "spin"
	optionTypeButton optionType = "button"
	optionTypeCheck  optionType = "check"
	optionTypeCombo  optionType = "combo"
	optionTypeString optionType = "string"
)

type option struct {
	name         string
	optionType   optionType
	defaultValue any
	min          int64
	max          int64
	options      []string
	apply        func(*engine, responder, string)
}

// var {name, type, default, min, max, apply}
var options = []option{
	{
		name:         "Hash",
		optionType:   optionTypeSpin,
		defaultValue: defaultEngineOptions.transpositionTableMemoryMB,
		min:          1,
		max:          33554432, // insane max of 32TB to match stockfish
		apply: func(e *engine, r responder, value string) {
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				r.debug("error parsing Hash option value: " + err.Error())
				return
			}
			e.options.transpositionTableMemoryMB = int(v)
			e.searcher = e.createSearcher(r)
		},
	},
	{
		name:       "Clear Hash",
		optionType: optionTypeButton,
		apply: func(e *engine, r responder, _ string) {
			// clear the transposition table
			r.debug("clearing transposition table")
			e.clearTranspositionTable()
		},
	},
	{
		name:         "Move Overhead",
		optionType:   optionTypeSpin,
		defaultValue: defaultEngineOptions.moveOverheadMS,
		min:          0,
		max:          5000,
		apply: func(e *engine, r responder, value string) {
			v, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				r.debug("error parsing Move Overhead option value: " + err.Error())
				return
			}
			e.options.moveOverheadMS = int(v)
		},
	},
}
