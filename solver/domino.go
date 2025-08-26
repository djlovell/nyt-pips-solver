package solver

import (
	"djlovell/nyt_pips_solver/input"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type domino struct {
	identifier string // in case we need to uniquely identify these
	val1       int
	val2       int
}

func (d domino) String() string {
	return fmt.Sprintf("[%d|%d]", d.val1, d.val2)
}

// parses a domino from an input specification
func parseInputDomino(d *input.Domino) (*domino, error) {
	if d == nil {
		panic("nil input domino")
	}
	if d.Val1 == nil {
		return nil, errors.New(`input file domino missing "val1`)
	}
	if d.Val2 == nil {
		return nil, errors.New(`input file domino missing "val2`)
	}
	val1, val2 := *d.Val1, *d.Val2
	if val1 < 0 || val1 > 6 || val2 < 0 || val2 > 6 {
		return nil, errors.New("domino values must be between 0 and 6")
	}
	return &domino{
		identifier: uuid.NewString(),
		val1:       val1,
		val2:       val2,
	}, nil
}
