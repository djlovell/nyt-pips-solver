package main

import (
	"errors"
	"fmt"
	"strings"
)

type Board struct {
	cells      [][]*Cell
	conditions []*Condition
}

// Loads a game board (utilized cells and conditions)
func InitializeBoard(inputCells [][]string, inputConditions [][]string) (*Board, error) {
	board := new(Board)

	// cell initialization
	{
		cells := make([][]*Cell, 0)
		for _, specRow := range inputCells {
			cellRow := make([]*Cell, 0)
			for _, c := range specRow {
				if p, err := parseInputCell(c); err != nil {
					return nil, err
				} else {
					cellRow = append(cellRow, p)
				}
			}
			cells = append(cells, cellRow)
		}

		// fill in cell positions and establish neighbors starting from the bottom/right of the board
		for yIdx := len(cells) - 1; yIdx >= 0; yIdx-- {
			for xIdx := len(cells[yIdx]) - 1; xIdx >= 0; xIdx-- {
				cell := cells[yIdx][xIdx]
				if cell == nil {
					return nil, errors.New("nil cell found during initialization")
				}

				// set the position of the current cell
				cell.posX = xIdx
				cell.posY = yIdx

				// unused cells do not get neighbors nor get to be neighbors
				if !cell.inPlay {
					continue
				}

				// establish neighbor associations with the neighbor to the left
				if xIdx > 0 {
					neighborLeft := cells[yIdx][xIdx-1]
					if neighborLeft == nil {
						return nil, errors.New("nil cell found during initialization")
					}
					if neighborLeft.inPlay {
						cell.neighborLeft = neighborLeft
						neighborLeft.neighborRight = cell
					}
				}

				// establish neighbor associations with the neighbor above
				if yIdx > 0 {
					neighborAbove := cells[yIdx-1][xIdx]
					if neighborAbove == nil {
						return nil, errors.New("nil cell found during initialization")
					}
					if neighborAbove.inPlay {
						cell.neighborAbove = neighborAbove
						neighborAbove.neighborBelow = cell
					}
				}
			}
		}
		board.cells = cells
	}

	// condition initialization
	{
		conditions := make([]*Condition, 0)
		for _, inputCond := range inputConditions {
			condition, err := parseInputCondition(inputCond)
			if err != nil {
				return nil, err
			}
			// fact check that all the locations are actually valid cells
			for _, conditionCellLoc := range condition.cellIdentifiers {
				xPos, yPos, err := CellIdentifierToBoardPos(conditionCellLoc)
				if err != nil {
					return nil, err
				}
				if yPos >= len(board.cells) {
					return nil, errors.New("input condition cell location out of y range")
				}
				if xPos >= len(board.cells[yPos]) {
					return nil, errors.New("input condition cell location out of x range")
				}
				if cell := board.cells[yPos][xPos]; !cell.inPlay {
					return nil, errors.New("input condition contains cell not in play")
				}
			}
			conditions = append(conditions, condition)
		}
		board.conditions = conditions
	}

	return board, nil
}

// I hate it but the board looks prettier
func (b Board) Print() {
	xIdxMax := 0
	fmt.Println("This is kinda what the board looks like...")
	rowStrings := []string{}
	for yIdx, r := range b.cells {
		rowCells := []string{fmt.Sprintf("%-3d", yIdx)} // TODO: make Y index pad up to 3 digits
		for xIdx, c := range r {
			cellStr := "["
			if c.inPlay {
				cellStr += " "
			} else {
				cellStr += "X"
			}
			cellStr += "]"
			rowCells = append(rowCells, cellStr)
			xIdxMax = max(xIdxMax, xIdx)
		}
		rowStrings = append(rowStrings, strings.Join(rowCells, " "))
	}
	headerRowValues := []string{"   "}
	for i := 0; i <= xIdxMax; i++ {
		headerRowValues = append(headerRowValues, fmt.Sprintf("%-3d", i))
	}
	rowStrings = append([]string{strings.Join(headerRowValues, " ")}, rowStrings...)
	fmt.Println(strings.Join(rowStrings, "\n") + "\n")
}
