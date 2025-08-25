package main

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
)

// DominoArrangementLocation - defines a grouping of cells where a domino could be placed based on which cells are in play
type DominoArrangementLocation struct {
	cell1 string // identifier
	cell2 string // identifier
}

// unique identifier for location (for de-duplication)
func (l *DominoArrangementLocation) Identifier() string {
	identifiers := []string{l.cell1, l.cell2}
	slices.Sort(identifiers)
	return strings.Join(identifiers, "-")
}

// DominoArrangement - defines a set of locations on a board where dominoes could fit
type DominoArrangement struct {
	locations []DominoArrangementLocation
}

func (s DominoArrangement) String() string {
	out := "Possible Fit Solution\n"
	for _, l := range s.locations {
		out += "\t" + l.Identifier() + "\n"
	}
	return out
}

// GetDominoArrangements - determines possible arrangements for laying dominoes on a board.
// Pre-computing valid domino positions will simplify solving later.
func GetDominoArrangements(board *Board) error {
	if board == nil {
		panic("nil board")
	}
	fmt.Println("Calculating possible arrangements for dominoes on the board...")

	// create a map of cells in play to keep track of
	cellsRemaining := make(map[string]*Cell)
	for yIdx := 0; yIdx < len(board.cells); yIdx++ {
		for xIdx := 0; xIdx < len(board.cells[yIdx]); xIdx++ {
			if cell := board.cells[yIdx][xIdx]; cell.inPlay {
				cellsRemaining[cell.Identifier()] = cell
			}
		}
	}

	locations := make([]DominoArrangementLocation, 0) // tracks locations of fitted dominoes for a possible solution
	solutions := make([]DominoArrangement, 0)         // tracks discovered fit solutions

	findDominoArrangements(board, cellsRemaining, locations, &solutions)
	if len(solutions) == 0 {
		return errors.New("no solutions found")
	}

	fmt.Printf("%d possible domino arrangements found...\n", len(solutions))
	for _, solution := range solutions {
		fmt.Println(solution.String())
	}

	return nil
}

// attempts to recurse through different ways of fitting dominoes to a board without using loops
// each recursive call will fit a domino into a cell and one of its neighbors, then remove the two from the remaining cells
func findDominoArrangements(board *Board, unarrangedCells map[string]*Cell, locations []DominoArrangementLocation, outArrangements *[]DominoArrangement) {
	if board == nil {
		panic("nil board")
	}
	if outArrangements == nil {
		panic("nil output arrangements")
	}

	// base case - all cells have been accounted for in the arrangement, so save it as a solution
	if len(unarrangedCells) == 0 {
		newSolution := DominoArrangement{
			locations: locations,
		}
		*outArrangements = append(*outArrangements, newSolution)
		debugPrint(fmt.Println, "all cells accounted for! solution added...")
		return
	}

	// grab the next cell to fit a domino in - is this efficient?
	nextCellIdentifier := slices.Collect(maps.Keys(unarrangedCells))[0]
	nextCell := unarrangedCells[nextCellIdentifier]

	// if at any point we encounter a cell that has no remaining neighbors that aren't accounted for...we have ran into an invalid fitment
	neighborFound := false

	for _, neighbor := range []*Cell{nextCell.neighborRight, nextCell.neighborBelow, nextCell.neighborLeft, nextCell.neighborAbove} {
		// is there a neighbor at all?
		if neighbor == nil {
			continue
		}
		// if the neighbor has been used by another domino, skip it
		if _, ok := unarrangedCells[neighbor.Identifier()]; !ok {
			continue
		}

		neighborFound = true

		// add the domino location to the list
		locationsNew := branchArrangementLocations(locations)
		locationsNew = append(locationsNew, DominoArrangementLocation{
			cell1: nextCell.Identifier(),
			cell2: neighbor.Identifier(),
		})

		// remove the cell and neighbor from the list and continue
		cellsRemainingNew := branchUnarrangedCells(unarrangedCells)
		delete(cellsRemainingNew, nextCell.Identifier())
		delete(cellsRemainingNew, neighbor.Identifier())
		findDominoArrangements(board, cellsRemainingNew, locationsNew, outArrangements)
	}

	if !neighborFound {
		debugPrint(fmt.Printf, "attempted arrangement resulted in an orphaned cell - %d cells unarranged...\n", len(unarrangedCells))
		return
	}
}

func branchUnarrangedCells(in map[string]*Cell) map[string]*Cell {
	out := make(map[string]*Cell)
	maps.Copy(out, in)
	return out
}

func branchArrangementLocations(in []DominoArrangementLocation) []DominoArrangementLocation {
	out := make([]DominoArrangementLocation, len(in))
	copy(out, in)
	return out
}
