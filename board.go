package main

import (
	"fmt"
	"strings"
)

type Board [][]*Cell

// Boards specs are expected to be consistently sized (e.g. 2x2, 16x4), and specify which cells are used for the game
func InitializeBoard(spec [][]CellType) Board {
	board := make([][]*Cell, 0)
	for _, specRow := range spec {
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

func (b Board) Print() {
	fmt.Println("This is kinda what the board looks like...")
	fmt.Println(b.String() + "\n")
}

func (b Board) String() string {
	rowStrings := []string{}
	for _, r := range b {
		rowCells := []string{}
		for _, c := range r {
			rowCells = append(rowCells, c.String())
		}
		rowStrings = append(rowStrings, strings.Join(rowCells, " "))
	}
	return strings.Join(rowStrings, "\n")
}
