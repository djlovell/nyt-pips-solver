package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

type Board [][]*cell

// make this a property determined at initialization time later...
func (b Board) NumCells() int {
	if len(b) == 0 {
		return 0
	}
	if len(b[0]) == 0 {
		return 0
	}
	return len(b) * len(b[0])
}

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

// which cells a domino could fit in on a board
type PossibleDominoLocation struct {
	cell1 *cell
	cell2 *cell
}

// unique identifier for location (for de-duplication)
func (l *PossibleDominoLocation) Identifier() string {
	identifiers := []string{l.cell1.Identifier(), l.cell2.Identifier()}
	slices.Sort(identifiers)
	return strings.Join(identifiers, "-")
}

func (l PossibleDominoLocation) String() string {
	return fmt.Sprintf("Cells %s - %s", l.cell1.Identifier(), l.cell2.Identifier())
}

type PossibleDominoFitSolution struct {
	locations []PossibleDominoLocation
}

func (s PossibleDominoFitSolution) String() string {
	out := "Possible Fit Solution\n"
	for _, l := range s.locations {
		out += "\t" + l.Identifier() + "\n"
	}
	return out
}

// figures out if and where dominoes can be laid
func (b Board) WillDominoesFit() bool {
	fmt.Println("Attempting to check if dominoes will fit on the board")
	filledCells := map[string]any{}                // track which cells have been accounted for
	locations := make([]PossibleDominoLocation, 0) // tracks locations of fitted dominoes for a possible solution
	solutions := make([]PossibleDominoFitSolution, 0)

	b.exploreDominoFitPath(filledCells, locations, &solutions)

	fmt.Printf("%d possible fitments found...\n", len(solutions))
	for _, solution := range solutions {
		fmt.Println(solution.String())
	}

	return false
}

// attempts to recurse through different ways of fitting dominoes to a board. Should return
// a status or error if the board is invalid, and eventually a list of possible fitments to
// try dominoes against later
//
// current algo does not work. It does a single double-loop pass through the board, which
// does not ensure that every cell is accounted for
func (b Board) exploreDominoFitPath(filledCells map[string]any, locations []PossibleDominoLocation, solutionsOut *[]PossibleDominoFitSolution) {
	for yIdx := 0; yIdx < len(b); yIdx++ {
		for xIdx := 0; xIdx < len(b[yIdx]); xIdx++ {
			cell := b[yIdx][xIdx]

			// base case - cell has already been used up
			if _, filled := filledCells[cell.Identifier()]; filled {
				continue
			}

			// if cell if unused, it doesn't need a domino, but is accounted for
			if !cell.inPlay {
				filledCells[cell.Identifier()] = true
				continue
			}

			cellOrphaned := true

			// try filling with a right neighbor
			if neighbor := cell.neighborRight; neighbor != nil {
				if _, neighborFilled := filledCells[neighbor.Identifier()]; !neighborFilled {
					locationsNew := copyLocations(locations)
					locationsNew = append(locationsNew, PossibleDominoLocation{
						cell1: cell,
						cell2: neighbor,
					})
					filledCellsNew := copyFilledCellMap(filledCells)
					fmt.Printf("looking right - cells %s and %s can be used for domino\n", cell.Identifier(), neighbor.Identifier())
					filledCellsNew[cell.Identifier()] = true
					filledCellsNew[neighbor.Identifier()] = true
					b.exploreDominoFitPath(filledCellsNew, locationsNew, solutionsOut)
				}

				cellOrphaned = false
			}

			// try filling with a below neighbor
			if neighbor := cell.neighborBelow; neighbor != nil {
				if _, neighborFilled := filledCells[neighbor.Identifier()]; !neighborFilled {
					locationsNew := copyLocations(locations)
					locationsNew = append(locationsNew, PossibleDominoLocation{
						cell1: cell,
						cell2: neighbor,
					})
					filledCellsNew := copyFilledCellMap(filledCells)
					fmt.Printf("looking below - cells %s and %s can be used for domino\n", cell.Identifier(), neighbor.Identifier())
					filledCellsNew[cell.Identifier()] = true
					filledCellsNew[neighbor.Identifier()] = true
					b.exploreDominoFitPath(filledCellsNew, locationsNew, solutionsOut)
				}

				cellOrphaned = false
			}

			// try filling with a left neighbor
			if neighbor := cell.neighborLeft; neighbor != nil {
				if _, neighborFilled := filledCells[neighbor.Identifier()]; !neighborFilled {
					locationsNew := copyLocations(locations)
					locationsNew = append(locationsNew, PossibleDominoLocation{
						cell1: cell,
						cell2: neighbor,
					})
					filledCellsNew := copyFilledCellMap(filledCells)
					fmt.Printf("looking left - cells %s and %s can be used for domino\n", cell.Identifier(), neighbor.Identifier())
					filledCellsNew[cell.Identifier()] = true
					filledCellsNew[neighbor.Identifier()] = true
					b.exploreDominoFitPath(filledCellsNew, locationsNew, solutionsOut)
				}

				cellOrphaned = false
			}

			// try filling with an above neighbor
			if neighbor := cell.neighborAbove; neighbor != nil {
				if _, neighborFilled := filledCells[neighbor.Identifier()]; !neighborFilled {
					locationsNew := copyLocations(locations)
					locationsNew = append(locationsNew, PossibleDominoLocation{
						cell1: cell,
						cell2: neighbor,
					})
					filledCellsNew := copyFilledCellMap(filledCells)
					fmt.Printf("looking above - cells %s and %s can be used for domino\n", cell.Identifier(), neighbor.Identifier())
					filledCellsNew[cell.Identifier()] = true
					filledCellsNew[neighbor.Identifier()] = true
					b.exploreDominoFitPath(filledCellsNew, locationsNew, solutionsOut)
				}

				cellOrphaned = false
			}

			// FAIL - cell has no neighbors...invalid board
			if cellOrphaned {
				panic("invalid board...cell stranded with no neighbors...no dominoes will fit")
			}

		}
	}

	// the end! found a solution...maybe. We at least got done iterating through the cells. fix later
	fmt.Println("we made it...through the board!")

	// not a solution if not every cell has been accounted for
	if len(filledCells) != b.NumCells() {
		return
	}

	// don't add if the solution is functionally the same - TODO
	newSolution := PossibleDominoFitSolution{
		locations: locations,
	}

	*solutionsOut = append(*solutionsOut, newSolution)
}

func copyFilledCellMap(in map[string]any) map[string]any {
	out := make(map[string]any)
	maps.Copy(out, in)
	return out
}

func copyLocations(in []PossibleDominoLocation) []PossibleDominoLocation {
	out := make([]PossibleDominoLocation, len(in))
	copy(out, in)
	return out
}
