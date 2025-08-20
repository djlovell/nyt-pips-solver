package main

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

// PossibleDominoLocation - defines a grouping of cells where a domino could be placed based on which cells are in play
type PossibleDominoLocation struct {
	cell1 *Cell
	cell2 *Cell
}

// unique identifier for location (for de-duplication)
func (l *PossibleDominoLocation) Identifier() string {
	identifiers := []string{l.cell1.Identifier(), l.cell2.Identifier()}
	slices.Sort(identifiers)
	return strings.Join(identifiers, "-")
}

// PossibleDominoFitSolution - defines a set of locations on a board where dominoes could fit
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
func (b Board) DominoFitCheck() {
	fmt.Println("Attempting to check if dominoes will fit on the board")

	// create a map of cells in play to keep track of
	cellsRemaining := make(map[string]*Cell)
	for yIdx := 0; yIdx < len(b); yIdx++ {
		for xIdx := 0; xIdx < len(b[yIdx]); xIdx++ {
			if cell := b[yIdx][xIdx]; cell.inPlay {
				cellsRemaining[cell.Identifier()] = cell
			}
		}
	}

	locations := make([]PossibleDominoLocation, 0)    // tracks locations of fitted dominoes for a possible solution
	solutions := make([]PossibleDominoFitSolution, 0) // tracks discovered fit solutions

	b.determinePossibleDominoLocations(cellsRemaining, locations, &solutions)
	if len(solutions) == 0 {
		panic("dominoes cannot fit on the defined board")
	}

	fmt.Printf("%d possible domino fit solution(s) found...\n", len(solutions))
	for _, solution := range solutions {
		fmt.Println(solution.String())
	}
}

// attempts to recurse through different ways of fitting dominoes to a board without using loops
// each recursive call will fit a domino into a cell and one of its neighbors, then remove the two from the remaining cells
func (b Board) determinePossibleDominoLocations(cellsRemaining map[string]*Cell, locations []PossibleDominoLocation, solutionsOut *[]PossibleDominoFitSolution) {
	// base case - all cells have been accounted for
	if len(cellsRemaining) == 0 {
		// don't add if the solution is functionally the same - TODO
		newSolution := PossibleDominoFitSolution{
			locations: locations,
		}

		*solutionsOut = append(*solutionsOut, newSolution)

		// fmt.Println("all cells accounted for! solution added...")
		return
	}

	// grab the next cell to work - is this efficient?
	nextCellIdentifier := slices.Collect(maps.Keys(cellsRemaining))[0]
	nextCell := cellsRemaining[nextCellIdentifier]

	// if at any point we encounter a cell that has no remaining neighbors that aren't accounted for...we have ran into an invalid fitment
	neighborFound := false

	for _, neighbor := range []*Cell{nextCell.neighborRight, nextCell.neighborBelow, nextCell.neighborLeft, nextCell.neighborAbove} {
		// is there a neighbor at all?
		if neighbor == nil {
			continue
		}
		// has the neighbor already been accounting for?
		if _, ok := cellsRemaining[neighbor.Identifier()]; !ok {
			continue
		}

		neighborFound = true

		// add the domino location to the list
		locationsNew := copyLocations(locations)
		locationsNew = append(locationsNew, PossibleDominoLocation{
			cell1: nextCell,
			cell2: neighbor,
		})
		// remove the cell and neighbor from the list and continue
		cellsRemainingNew := copyCellsRemaining(cellsRemaining)
		delete(cellsRemainingNew, nextCell.Identifier())
		delete(cellsRemainingNew, neighbor.Identifier())
		b.determinePossibleDominoLocations(cellsRemainingNew, locationsNew, solutionsOut)
	}

	if !neighborFound {
		// fmt.Printf("attempted fit resulted in an unfillable cell, %d cells unaccounted for...\n", len(cellsRemaining))
		return
	}
}

func copyCellsRemaining(in map[string]*Cell) map[string]*Cell {
	out := make(map[string]*Cell)
	maps.Copy(out, in)
	return out
}

func copyLocations(in []PossibleDominoLocation) []PossibleDominoLocation {
	out := make([]PossibleDominoLocation, len(in))
	copy(out, in)
	return out
}
