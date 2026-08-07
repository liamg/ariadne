package uci

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/liamg/ariadne/version"
)

func TestHandshake(t *testing.T) {
	harness := newTestHarness(t)
	defer harness.Close()

	u := New(harness, harness)
	go func() { _ = u.Run(t.Context()) }()

	harness.Send("uci")
	harness.WaitForExact("uciok", time.Second*10)

	messages := harness.Messages()

	var ids []string
	var receivedOptions []string
	for i, msg := range messages {
		if len(msg) >= 3 && msg[:3] == "id " {
			if len(receivedOptions) > 0 {
				t.Errorf("id message received after option messages: %s", msg)
			}
			ids = append(ids, msg)
		}
		if len(msg) >= 7 && msg[:7] == "option " {
			receivedOptions = append(receivedOptions, msg)
		}
		if msg == "uciok" {
			if i < len(messages)-1 {
				t.Errorf("uciok message received before end of messages: %s", msg)
			}
		}
	}

	var name string
	var v string
	var author string
	for _, id := range ids {
		fields := strings.Fields(id)
		if len(fields) < 3 {
			t.Fatalf("invalid id message: %s", id)
		}
		switch fields[1] {
		case "name":
			name = fields[2]
			if len(fields) > 3 {
				v = fields[3]
			}
		case "author":
			author = strings.Join(fields[2:], " ")
		default:
			t.Errorf("unknown id field: %s", fields[1])
		}
	}

	if name != version.Name {
		t.Errorf("unexpected engine name: %s", name)
	}

	if v != version.Version {
		t.Errorf("unexpected engine version: %s", v)
	}

	if author != version.Author {
		t.Errorf("unexpected engine author: %s", author)
	}

	expectedOptions := make(map[string]option, len(options))
	for _, opt := range options {
		expectedOptions[opt.name] = opt
	}
	actualOptions := make(map[string]option, len(receivedOptions))

	for _, optionrec := range receivedOptions {
		if !strings.HasPrefix(optionrec, "option name ") {
			t.Errorf("invalid option message: %s", optionrec)
			continue
		}
		name, remainder, ok := strings.Cut(optionrec[len("option name "):], " type ")
		if !ok {
			t.Errorf("invalid option message: %s", optionrec)
			continue
		}

		var opt option

		ot, remainder, _ := strings.Cut(remainder, " ")
		opt.optionType = optionType(ot)

		remainder = strings.TrimSpace(remainder)
		switch opt.optionType {
		case optionTypeSpin:
			if !strings.HasPrefix(remainder, "default ") {
				t.Errorf("invalid spin option message: %s", optionrec)
				continue
			}
			remainder = remainder[len("default "):]
			defaultValue, remainder, _ := strings.Cut(remainder, " ")
			opt.defaultValue = defaultValue
			remainder = strings.TrimSpace(remainder)
			if !strings.HasPrefix(remainder, "min ") {
				t.Errorf("invalid spin option message: %s", optionrec)
				continue
			}
			remainder = remainder[len("min "):]
			minStr, remainder, _ := strings.Cut(remainder, " ")
			minValue, err := strconv.ParseInt(minStr, 10, 64)
			if err != nil {
				t.Errorf("invalid spin option message: %s", optionrec)
				continue
			}
			opt.min = minValue
			remainder = strings.TrimSpace(remainder)
			if !strings.HasPrefix(remainder, "max ") {
				t.Errorf("invalid spin option message: %s", optionrec)
				continue
			}
			remainder = remainder[len("max "):]
			maxValue, err := strconv.ParseInt(remainder, 10, 64)
			if err != nil {
				t.Errorf("invalid spin option message: %s", optionrec)
				continue
			}
			opt.max = maxValue
		case optionTypeCheck:
			if !strings.HasPrefix(remainder, "default ") {
				t.Errorf("invalid check option message: %s", optionrec)
				continue
			}
			remainder = remainder[len("default "):]
			defaultValue := remainder
			opt.defaultValue = defaultValue == "true"
		case optionTypeCombo:
			if !strings.HasPrefix(remainder, "default ") {
				t.Errorf("invalid combo option message: %s", optionrec)
				continue
			}
			remainder = remainder[len("default "):]
			defaultValue, remainder, _ := strings.Cut(remainder, " ")
			opt.defaultValue = defaultValue
			remainder = strings.TrimSpace(remainder)
			for strings.HasPrefix(remainder, "var ") {
				remainder = remainder[len("var "):]
				var varValue string
				varValue, remainder, _ = strings.Cut(remainder, " ")
				opt.options = append(opt.options, varValue)
				remainder = strings.TrimSpace(remainder)
			}
		case optionTypeString:
			if !strings.HasPrefix(remainder, "default ") {
				t.Errorf("invalid string option message: %s", optionrec)
				continue
			}
			remainder = remainder[len("default "):]
			defaultValue := remainder
			opt.defaultValue = defaultValue
		case optionTypeButton:
			if strings.TrimSpace(remainder) != "" {
				t.Errorf("invalid button option message: %s", optionrec)
				continue
			}
		}
		actualOptions[name] = opt
	}

	for expected, expectedOption := range expectedOptions {
		actualOption, ok := actualOptions[expected]
		if !ok {
			t.Errorf("expected option not found: %s", expected)
			continue
		}
		if actualOption.optionType != expectedOption.optionType {
			t.Errorf("option type mismatch for %s: expected %s, got %s", expected, expectedOption.optionType, actualOption.optionType)
		}
		if fmt.Sprintf("%v", actualOption.defaultValue) != fmt.Sprintf("%v", expectedOption.defaultValue) {
			t.Errorf("option default value mismatch for %s: expected %v, got %v", expected, expectedOption.defaultValue, actualOption.defaultValue)
		}
		if actualOption.min != expectedOption.min {
			t.Errorf("option min value mismatch for %s: expected %d, got %d", expected, expectedOption.min, actualOption.min)
		}
		if actualOption.max != expectedOption.max {
			t.Errorf("option max value mismatch for %s: expected %d, got %d", expected, expectedOption.max, actualOption.max)
		}
		if len(actualOption.options) != len(expectedOption.options) {
			t.Errorf("option options length mismatch for %s: expected %d, got %d", expected, len(expectedOption.options), len(actualOption.options))
		} else {
			for i := range actualOption.options {
				if fmt.Sprintf("%v", actualOption.options[i]) != fmt.Sprintf("%v", expectedOption.options[i]) {
					t.Errorf("option options value mismatch for %s at index %d: expected %v, got %v", expected, i, expectedOption.options[i], actualOption.options[i])
				}
			}
		}
	}
}
