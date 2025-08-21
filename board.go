package main

import (
	"fmt"
	"strings"
)

type Board [][]*Cell

// Boards specs are expected to be consistently sized (e.g. 2x2, 16x4), and specify which cells are used for the game
func InitializeBoard(input [][]CellType) Board {
	board := make([][]*Cell, 0)
	for _, specRow := range input {
		cellRow := make([]*Cell, 0)
		for _, t := range specRow {
			cellRow = append(cellRow, NewCell(t))
		}
		board = append(board, cellRow)
	}

	// fill in cell positions and establish neighbors starting from the bottom/right of the board
	for yIdx := len(board) - 1; yIdx >= 0; yIdx-- {
		for xIdx := len(board[yIdx]) - 1; xIdx >= 0; xIdx-- {
			cell := board[yIdx][xIdx]
			if cell == nil {
				panic("cell is nil")
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
				neighborLeft := board[yIdx][xIdx-1]
				if neighborLeft == nil {
					panic("neighbor is nil")
				}
				if neighborLeft.inPlay {
					cell.neighborLeft = neighborLeft
					neighborLeft.neighborRight = cell
				}
			}

			// establish neighbor associations with the neighbor above
			if yIdx > 0 {
				neighborAbove := board[yIdx-1][xIdx]
				if neighborAbove == nil {
					panic("neighbor is nil")
				}
				if neighborAbove.inPlay {
					cell.neighborAbove = neighborAbove
					neighborAbove.neighborBelow = cell
				}
			}
		}
	}
	return board
}

// I hate it but the board looks prettier
func (b Board) Print() {
	xIdxMax := 0
	fmt.Println("This is kinda what the board looks like...")
	rowStrings := []string{}
	for yIdx, r := range b {
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
