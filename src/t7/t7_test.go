package t7

import "testing"

func Test7_0(t *testing.T) {
	m := map[string]func(op int) int{}
	m["1"] = func(op int) int {
		return op
	}

	m["2"] = func(op int) int {
		return op * op
	}

	m["3"] = func(op int) int {
		return op * op * op
	}

	t.Log(m["1"](2), m["2"](2), m["3"](2))
}
