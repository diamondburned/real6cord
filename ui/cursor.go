package ui

import (
	"strconv"
)

// RETs pos vertical, pos horizontal (line, col)
func getCursorPos() (int, int) {
	rs := queryTerm("\033[6n")
	var X, Y string

	var grab int

ParseLoop:
	for _, r := range rs {
		switch r {
		case '[':
			grab = 1
		case ';':
			grab = 2
		case 'R':
			break ParseLoop
		default:
			switch grab {
			case 1:
				X += string(r)
			case 2:
				Y += string(r)
			}
		}
	}

	l, err := strconv.Atoi(X)
	if err != nil {
		panic(err)
	}

	c, err := strconv.Atoi(Y)
	if err != nil {
		panic(err)
	}

	return l, c
}
