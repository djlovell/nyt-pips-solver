package solver

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type domino struct {
	val1 int
	val2 int
}

func (d domino) String() string {
	return fmt.Sprintf("[%d|%d]", d.val1, d.val2)
}

// parses a domino from an input specification
func parseInputDomino(s string) (*domino, error) {
	vals := strings.Split(s, "|")
	if len(vals) != 2 {
		return nil, fmt.Errorf("%s is not a valid domino input format", s)
	}
	val1, err := strconv.Atoi(vals[0])
	if err != nil {
		return nil, fmt.Errorf("failed to parse domino value - %w", err)
	}
	val2, err := strconv.Atoi(vals[1])
	if err != nil {
		return nil, fmt.Errorf("failed to parse domino value - %w", err)
	}
	if val1 < 0 || val1 > 6 || val2 < 0 || val2 > 6 {
		return nil, errors.New("domino values must be between 0 and 6")
	}
	return &domino{
		val1: val1,
		val2: val2,
	}, nil
}
