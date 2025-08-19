package main

import (
	"fmt"
	"maps"
	"strings"
)

type Board [][]*cell

// Boards specs are expected to be consistently sized (e.g. 2x2, 16x4), and specify which cells are used for the game
func InitializeBoard(spec [][]CellType) Board {
	board := make([][]*cell, 0)
	for _, specRow := range spec {
		cellRow := make([]*cell, 0)
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
				if !neighborLeft.inPlay {
					continue
				}
				cell.neighborLeft = neighborLeft
				neighborLeft.neighborRight = cell
			}

			// establish neighbor associations with the neighbor above
			if yIdx > 0 {
				neighborAbove := board[yIdx-1][xIdx]
				if neighborAbove == nil {
					panic("neighbor is nil")
				}
				if !neighborAbove.inPlay {
					continue
				}
				cell.neighborAbove = neighborAbove
				neighborAbove.neighborBelow = cell
			}
		}
	}
	return board
}

func (b Board) Print() {
	fmt.Println(b.String())
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

// type DominoPlace struct {
// 	cell1 *cell
// 	cell2 *cell
// }

// figures out if and where dominoes can be laid
func (b Board) WillDominoesFit() bool {
	fmt.Println("Attempting to check if dominoes will fit on the board")
	filledCells := map[string]any{} // track which cells have been accounted for

	return b.exploreDominoFitPath(filledCells)
}

// attempts to recurse through different ways of fitting dominoes to a board. Should return
// a status or error if the board is invalid, and eventually a list of possible fitments to
// try dominoes against later
//
// current algo does not work. It does a single double-loop pass through the board, which
// does not ensure that every cell is accounted for
func (b Board) exploreDominoFitPath(filledCells map[string]any) bool {
	for yIdx := 0; yIdx < len(b); yIdx++ {
		for xIdx := 0; xIdx < len(b[yIdx]); xIdx++ {
			cell := b[yIdx][xIdx]

			// base case - cell has already been used up
			if _, filled := filledCells[cell.Identifier()]; filled {
				continue
			}

			// if cell if unused, it doesn't need a domino
			if !cell.inPlay {
				filledCells[cell.Identifier()] = true
				continue
			}

			foundPath := false

			// try filling with a right neighbor
			if cell.neighborRight != nil {
				filledCellsNew := copyFilledCellMap(filledCells)
				fmt.Printf("looking right, cells %s and %s can be used for domino\n", cell.Identifier(), cell.neighborRight.Identifier())
				filledCellsNew[cell.Identifier()] = true
				filledCellsNew[cell.neighborRight.Identifier()] = true
				if foundPath = b.exploreDominoFitPath(filledCellsNew); foundPath {
					return true
				}
			}

			// try filling with a below neighbor
			if cell.neighborBelow != nil {
				filledCellsNew := copyFilledCellMap(filledCells)
				fmt.Printf("looking down, cells %s and %s can be used for domino\n", cell.Identifier(), cell.neighborBelow.Identifier())
				filledCellsNew[cell.Identifier()] = true
				filledCellsNew[cell.neighborBelow.Identifier()] = true
				if foundPath = b.exploreDominoFitPath(filledCellsNew); foundPath {
					return true
				}
			}

			// try filling with a left neighbor
			if cell.neighborLeft != nil {
				filledCellsNew := copyFilledCellMap(filledCells)
				fmt.Printf("looking left, cells %s and %s can be used for domino\n", cell.Identifier(), cell.neighborLeft.Identifier())
				filledCellsNew[cell.Identifier()] = true
				filledCellsNew[cell.neighborLeft.Identifier()] = true
				if foundPath = b.exploreDominoFitPath(filledCellsNew); foundPath {
					return true
				}
			}

			// try filling with an above neighbor
			if cell.neighborAbove != nil {
				filledCellsNew := copyFilledCellMap(filledCells)
				fmt.Printf("looking up, cells %s and %s can be used for domino\n", cell.Identifier(), cell.neighborAbove.Identifier())
				filledCellsNew[cell.Identifier()] = true
				filledCellsNew[cell.neighborAbove.Identifier()] = true
				if foundPath = b.exploreDominoFitPath(filledCellsNew); foundPath {
					return true
				}
			}

			// FAIL - cell has no neighbors...invalid board
			panic("invalid board...cell stranded with no neighbors...no dominoes will fit")
		}
	}

	fmt.Println("we made it!")
	return true
}

func copyFilledCellMap(in map[string]any) map[string]any {
	out := make(map[string]any)
	maps.Copy(out, in)
	return out
}
