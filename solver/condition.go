package solver

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
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

// check - returns if cell values satisfy the condition or not
func (c condition) check(cellValues map[string] /* cell identifier */ int /*cell value */) bool {
	switch c.expression {
	case conditionExpSumEquals:
		sum := 0
		for _, cell := range c.cellIdentifiers {
			sum += cellValues[cell]
		}
		if sum != c.operand {
			return false
		}
	case conditionExpSumLessThan:
		sum := 0
		for _, cell := range c.cellIdentifiers {
			sum += cellValues[cell]
		}
		if sum >= c.operand {
			return false
		}
	case conditionExpSumGreaterThan:
		sum := 0
		for _, cell := range c.cellIdentifiers {
			sum += cellValues[cell]
		}
		if sum <= c.operand {
			return false
		}
	case conditionExpEquivalent:
		// just use first value as the "norm" and fail if anything else doesn't match
		expectedVal := cellValues[c.cellIdentifiers[0]]
		for _, cell := range c.cellIdentifiers {
			if cellValues[cell] != expectedVal {
				return false
			}
		}
	case conditionExpDistinct:
		foundVals := make(map[int]bool)
		for _, cell := range c.cellIdentifiers {
			cellVal := cellValues[cell]
			if _, alreadyExists := foundVals[cellVal]; alreadyExists {
				return false
			}
			foundVals[cellVal] = true
		}
	default:
		panic("unexpected condition expression type")
	}

	return true
}

// parses a Condition from input specification
func parseInputCondition(input []string) (*condition, error) {
	debugPrint(fmt.Printf, "Parsing input condition - %v\n", input)

	if len(input) == 0 {
		return nil, errors.New("empty input condition")
	}

	inputPos := 0

	var outExpression conditionExp
	var outOperand int
	outCellIdentifiers := make([]string, 0)

	// grab the expression first
	inExpression := input[0]
	operandRequired := false
	switch inExpression {
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
		return nil, fmt.Errorf("%s is not a recognized input condition expression", inExpression)
	}
	inputPos++

	// figure out if an operand is expected, and if needed, try to parse that next
	if operandRequired {
		inOperand := input[1]
		if len(input) < 3 {
			return nil, errors.New("invalid input condition - missing either an operand or cell locations")
		}
		if i, err := strconv.Atoi(inOperand); err != nil {
			return nil, fmt.Errorf("failed to parse error condition operand - %w", err)
		} else if i < 0 {
			return nil, errors.New("error condition operand should not be negative")
		} else {
			outOperand = i
		}
		inputPos++
	} else {
		if len(input) < 2 {
			return nil, errors.New("invalid input condition - missing cell locations")
		}
	}

	// parse cell positions next
	for i := inputPos; i < len(input); i++ {
		inCellPos := input[i]
		xPos, yPos, err := cellIdentifierToBoardPos(inCellPos)
		if err != nil {
			return nil, fmt.Errorf("failed to parse error condition cell position - %w", err)
		}
		outCellIdentifiers = append(outCellIdentifiers, boardPosToCellIdentifier(xPos, yPos))
	}

	return &condition{
		expression:      outExpression,
		operand:         outOperand,
		cellIdentifiers: outCellIdentifiers,
	}, nil
}
