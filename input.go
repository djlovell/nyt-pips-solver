package main

import (
	"errors"
	"fmt"
	"strconv"
)

/*
Input Specification

When describing game boards, a grid system is used.
The locations of cells, dominoes, conditions, etc. are consistently described and will be referred to
by their (x/column, y/row) position, zero-indexed and increasing down and to the right when visualizing the game board.

When inputting the board, a uniform grid (consistent row widths and column heights) is expected, and shall
be a [][]string.

Cells used in the game (aka that will have a domino placed on them) will be described as:
- "O" for in-play
- "X" if they are not part of the game

"Conditions" in the game apply to groups of cells and can be any of the following expressions:
- "N", where N is a number that domino faces covering the group of cells must add up to
- "=", where the domino faces must all be equivalent
- "!=", where the domino faces must all be distinct
- ">N", where the domino faces must be cumulatively greater than N
- "<N", where the domino faces must be cumulatively less than N

When inputting, conditions shall be a [][]string, each []string condition following the format:
{"{{expression from above}}", "{{integer operand for expressions that have N in them}}", "X1:Y1", "X2:Y2", etc. for any other cells the condition covers}
Example: {"N", "12", "0:0", "0:1"} implies that the first two cells in the top left corner of the board must add up to 12

Not every cell needs to be covered by a condition, and conditions may overlap across cells

Dominoes shall be input as a []string, each string of the format "N|M", with no particular polarity, but where 0 <= N,M <= 6.
*/

// in-game tutorial (at least as of 8/22)
var inputCells = [][]string{
	{"X", "O", "O"},
	{"X", "O", "O"},
	{"O", "O", "X"},
}

var inputConditions = [][]string{
	{"N", "0", "0:2"},
	{"=", "1:1", "1:2"},
	{"N", "10", "2:0", "2:1"},
}

var inputDominoes = []string{
	"5|5", "0|2", "2|3",
}

// parses a Cell from an input specification
func parseInputCell(s string) (*Cell, error) {
	switch s {
	case "X":
		return &Cell{inPlay: false}, nil
	case "O":
		return &Cell{inPlay: true}, nil
	default:
		return nil, fmt.Errorf("%s is an unknown input cell type", s)
	}
}

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
