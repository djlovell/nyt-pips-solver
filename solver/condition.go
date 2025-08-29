package solver

import (
	"errors"
	"fmt"
	"strings"

	"djlovell/nyt_pips_solver/input"
)

type condition struct {
	expression      conditionExp
	operand         int      // goes with some conditions
	cellIdentifiers []string // cell identifiers
}

type conditionExp int

const (
	conditionExpSumEquals conditionExp = iota
	conditionExpSumLessThan
	conditionExpSumGreaterThan
	conditionExpEquivalent
	conditionExpDistinct
)

func (c condition) String() string {
	s := "Cell"
	switch len(c.cellIdentifiers) {
	case 0:
		panic("no cells for condition")
	case 1:
		s += " " + c.cellIdentifiers[0]
	case 2:
		s += fmt.Sprintf("s %s & %s", c.cellIdentifiers[0], c.cellIdentifiers[1])
	default: // > 2 cells
		s += fmt.Sprintf("s %s", strings.Join(c.cellIdentifiers, ", "))
		if idx := strings.LastIndex(s, ","); idx != -1 {
			sBeforeAnd := s[0 : idx+1]
			sAfterAnd := s[idx+1:]
			s = sBeforeAnd + " &" + sAfterAnd
		}
	}
	s += " must "
	switch c.expression {
	case conditionExpSumEquals:
		s += fmt.Sprintf("add up to %d", c.operand)
	case conditionExpSumLessThan:
		s += fmt.Sprintf("add up to less than %d", c.operand)
	case conditionExpSumGreaterThan:
		s += fmt.Sprintf("add up to greater than %d", c.operand)
	case conditionExpEquivalent:
		s += "all be the same"
	case conditionExpDistinct:
		s += "all be different"
	default:
		panic("unhandled expression type")
	}
	return s
}

var errConditionNotReadyToCheck = errors.New("condition's cells are not ready to check yet")

// check - returns if cell values satisfy the condition or not
// returns errConditionNotReadyToCheck if not all cells have been filled
func (c condition) check(cellValues map[string] /* cell identifier */ int /*cell value */) (bool, error) {
	for _, cell := range c.cellIdentifiers {
		if _, ok := cellValues[cell]; !ok {
			return false, errConditionNotReadyToCheck
		}
	}

	switch c.expression {
	case conditionExpSumEquals:
		sum := 0
		for _, cell := range c.cellIdentifiers {
			sum += cellValues[cell]
		}
		if sum != c.operand {
			return false, nil
		}
	case conditionExpSumLessThan:
		sum := 0
		for _, cell := range c.cellIdentifiers {
			sum += cellValues[cell]
		}
		if sum >= c.operand {
			return false, nil
		}
	case conditionExpSumGreaterThan:
		sum := 0
		for _, cell := range c.cellIdentifiers {
			sum += cellValues[cell]
		}
		if sum <= c.operand {
			return false, nil
		}
	case conditionExpEquivalent:
		// just use first value as the "norm" and fail if anything else doesn't match
		expectedVal := cellValues[c.cellIdentifiers[0]]
		for _, cell := range c.cellIdentifiers {
			if cellValues[cell] != expectedVal {
				return false, nil
			}
		}
	case conditionExpDistinct:
		foundVals := make(map[int]bool)
		for _, cell := range c.cellIdentifiers {
			cellVal := cellValues[cell]
			if _, alreadyExists := foundVals[cellVal]; alreadyExists {
				return false, nil
			}
			foundVals[cellVal] = true
		}
	default:
		panic("unexpected condition expression type")
	}

	return true, nil
}

// parses a Condition from input specification
func parseInputCondition(input *input.Condition) (*condition, error) {
	if input == nil {
		panic("nil input condition")
	}

	var outExpression conditionExp
	var outOperand int
	outCellIdentifiers := make([]string, 0)

	// grab the expression first
	if input.Expression == nil {
		return nil, errors.New(`input file condition missing "expression"`)
	}
	operandRequired := false
	switch e := *input.Expression; e {
	// the ones that need an operand
	case "N":
		outExpression = conditionExpSumEquals
		operandRequired = true
	case ">N":
		outExpression = conditionExpSumGreaterThan
		operandRequired = true
	case "<N":
		outExpression = conditionExpSumLessThan
		operandRequired = true
	// the ones that don't need an operand
	case "=":
		outExpression = conditionExpEquivalent
	case "!=":
		outExpression = conditionExpDistinct
	default:
		return nil, fmt.Errorf("%s is not a recognized input condition expression", e)
	}

	// figure out if an operand is expected, and if needed, try to parse that next
	if operandRequired {
		if input.Operand == nil {
			return nil, fmt.Errorf(`input file condition expression "%s" missing "operand"`, *input.Expression)
		}
		if o := *input.Operand; o < 0 {
			return nil, fmt.Errorf("input file condition operand cannot be a negative number")
		} else {
			outOperand = o
		}
	}

	// parse cell positions next
	for _, c := range input.Cells {
		if c.X == nil {
			return nil, errors.New(`input file condition cell missing "x" position`)
		}
		if c.Y == nil {
			return nil, errors.New(`input file condition cell missing "y" position`)
		}
		outCellIdentifiers = append(outCellIdentifiers, boardPosToCellIdentifier(*c.X, *c.Y))
	}

	return &condition{
		expression:      outExpression,
		operand:         outOperand,
		cellIdentifiers: outCellIdentifiers,
	}, nil
}
