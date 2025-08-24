package main

import (
	"errors"
	"fmt"
	"strconv"
)

type Condition struct {
	expression      ConditionExp
	operand         int      // goes with some conditions
	cellIdentifiers []string // cell identifiers
}

type ConditionExp int

const (
	ConditionExpSumEquals ConditionExp = iota
	ConditionExpEquivalent
	ConditionExpDistinct
	ConditionExpSumLessThan
	ConditionExpSumGreaterThan
)

// parses a Condition from input specification
func parseInputCondition(input []string) (*Condition, error) {
	debugPrint(fmt.Printf, "Parsing input condition - %v\n", input)

	if len(input) == 0 {
		return nil, errors.New("empty input condition")
	}

	inputPos := 0

	var outExpression ConditionExp
	var outOperand int
	outCellIdentifiers := make([]string, 0)

	// grab the expression first
	inExpression := input[0]
	operandRequired := false
	switch inExpression {
	// the ones that need an operand
	case "N":
		outExpression = ConditionExpSumEquals
		operandRequired = true
	case ">N":
		outExpression = ConditionExpSumGreaterThan
		operandRequired = true
	case "<N":
		outExpression = ConditionExpSumLessThan
		operandRequired = true
	// the ones that don't need an operand
	case "=":
		outExpression = ConditionExpEquivalent
	case "!=":
		outExpression = ConditionExpDistinct
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
		xPos, yPos, err := CellIdentifierToBoardPos(inCellPos)
		if err != nil {
			return nil, fmt.Errorf("failed to parse error condition cell position - %w", err)
		}
		outCellIdentifiers = append(outCellIdentifiers, BoardPosToCellIdentifier(xPos, yPos))
	}

	return &Condition{
		expression:      outExpression,
		operand:         outOperand,
		cellIdentifiers: outCellIdentifiers,
	}, nil
}
