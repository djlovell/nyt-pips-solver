package solver

import (
	"errors"
	"fmt"
	"maps"
	"slices"
)

// DominoArrangementLocation - defines a grouping of cells where a domino could be placed based on which cells are in play
type DominoArrangementLocation struct {
	cell1 string // identifier
	cell2 string // identifier
}

func (a DominoArrangementLocation) String() string {
	identifiers := []string{a.cell1, a.cell2}
	slices.Sort(identifiers)
	return fmt.Sprintf("Cells %s-%s\n", identifiers[0], identifiers[1])
}

// DominoArrangement - defines a set of locations on a board where dominoes could fit
type DominoArrangement struct {
	locations []DominoArrangementLocation
}

func (a DominoArrangement) String() string {
	out := "Possible Arrangement\n"
	for _, l := range a.locations {
		out += "  " + l.String()
	}
	return out
}

// GetDominoArrangements - determines possible arrangements for laying dominoes on a board.
// Pre-computing valid domino positions will simplify solving later.
func GetDominoArrangements(game *Game) ([]DominoArrangement, error) {
	if game == nil {
		panic("nil board")
	}
	fmt.Println("Calculating possible arrangements for dominoes on the board...")

	// create a map of cells in play to keep track of
	cellsRemaining := branchUnarrangedCells(game.inPlayCellsByIdentifier)

	locations := make([]DominoArrangementLocation, 0) // tracks locations of fitted dominoes for a possible solution
	arrangements := make([]DominoArrangement, 0)      // tracks discovered fit solutions

	findDominoArrangements(game, cellsRemaining, locations, &arrangements)
	if len(arrangements) == 0 {
		return nil, errors.New("no solutions found")
	}

	fmt.Printf("%d possible domino arrangements found...\n", len(arrangements))
	for _, solution := range arrangements {
		fmt.Println(solution.String())
	}

	return arrangements, nil
}

// attempts to recurse through different ways of fitting dominoes to a board without using loops
// each recursive call will fit a domino into a cell and one of its neighbors, then remove the two from the remaining cells
func findDominoArrangements(game *Game, unarrangedCells map[string]*cell, locations []DominoArrangementLocation, outArrangements *[]DominoArrangement) {
	if game == nil {
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
		debugPrint(fmt.Println, "All cells accounted for and arrangement added...")
		return
	}

	// grab the next cell to fit a domino in - is this efficient?
	nextCellIdentifier := slices.Collect(maps.Keys(unarrangedCells))[0]
	nextCell := unarrangedCells[nextCellIdentifier]

	// if at any point we encounter a cell that has no remaining neighbors that aren't accounted for...we have ran into an invalid fitment
	neighborFound := false

	for _, neighbor := range []*cell{nextCell.neighborRight, nextCell.neighborBelow, nextCell.neighborLeft, nextCell.neighborAbove} {
		// is there a neighbor at all?
		if neighbor == nil {
			continue
		}
		// if the neighbor has been used by another domino, skip it
		if _, ok := unarrangedCells[neighbor.identifier()]; !ok {
			continue
		}

		neighborFound = true

		// add the domino location to the list
		locationsNew := branchArrangementLocations(locations)
		locationsNew = append(locationsNew, DominoArrangementLocation{
			cell1: nextCell.identifier(),
			cell2: neighbor.identifier(),
		})

		// remove the cell and neighbor from the list and continue
		cellsRemainingNew := branchUnarrangedCells(unarrangedCells)
		delete(cellsRemainingNew, nextCell.identifier())
		delete(cellsRemainingNew, neighbor.identifier())
		findDominoArrangements(game, cellsRemainingNew, locationsNew, outArrangements)
	}

	if !neighborFound {
		debugPrint(fmt.Printf, "Attempted arrangement resulted in an orphaned cell - %d cells unarranged...\n", len(unarrangedCells))
		return
	}
}

func branchUnarrangedCells(in map[string]*cell) map[string]*cell {
	out := make(map[string]*cell)
	maps.Copy(out, in)
	return out
}

func branchArrangementLocations(in []DominoArrangementLocation) []DominoArrangementLocation {
	out := make([]DominoArrangementLocation, len(in))
	copy(out, in)
	return out
}
