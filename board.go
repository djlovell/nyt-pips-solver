package main

import (
	"errors"
	"fmt"
	"strings"
)

type Board struct {
	cells [][]*Cell
}

// Loads a game board (utilized cells and conditions)
func InitializeBoard(inputCells [][]string) (*Board, error) {
	board := new(Board)

	// cell initialization
	{
		cells := make([][]*Cell, 0)
		for _, specRow := range inputCells {
			cellRow := make([]*Cell, 0)
			for _, c := range specRow {
				cellRow = append(cellRow, parseInputCell(c))
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

	// condition initialization TODO

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
