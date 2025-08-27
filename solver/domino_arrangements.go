package solver

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

// DominoArrangementLocation - defines a grouping of cells where a domino could be placed based on which cells are in play
type DominoArrangementLocation struct {
	cell1 string // identifier
	cell2 string // identifier
	// potential optimization - pre-determine dominoes that can't go in this location
	blacklistedDominoIDs *map[string]any
}

func (a DominoArrangementLocation) String() string {
	identifiers := []string{a.cell1, a.cell2}
	slices.Sort(identifiers)
	return fmt.Sprintf("Cells %s-%s\n", identifiers[0], identifiers[1])
}

// experiment - filter down dominoes that can go in this location for later checking
func (a *DominoArrangementLocation) addBlacklistedDominoIDs(g *Game) *DominoArrangementLocation {
	conditionsForLocation := append(
		g.inPlayCellsByIdentifier[a.cell1].applicableConditions,
		g.inPlayCellsByIdentifier[a.cell2].applicableConditions...,
	)
	invalidDominoes := make(map[string]any)
	for _, d := range g.dominoes {
		for _, c := range conditionsForLocation {
			// if both values in a domino fail a condition...blacklist the domino
			switch c.expression {
			case conditionExpSumEquals:
				// both domino values exceed
				if d.val1 > c.operand && d.val2 > c.operand {
					invalidDominoes[d.identifier] = true
					debugPrint(fmt.Printf, "Domino %s blacklisted for location %s\n", d.String(), a.String())
				}
			case conditionExpSumLessThan:
				// both domino values meet or exceed
				if d.val1 >= c.operand && d.val2 >= c.operand {
					invalidDominoes[d.identifier] = true
					debugPrint(fmt.Printf, "Domino %s blacklisted for location %s\n", d.String(), a.String())
				}
			case conditionExpSumGreaterThan:
				// the condition only uses one cell and neither domino value is sufficient
				if len(c.cellIdentifiers) == 1 {
					if d.val1 <= c.operand && d.val2 <= c.operand {
						invalidDominoes[d.identifier] = true
						debugPrint(fmt.Printf, "Domino %s blacklisted for location %s\n", d.String(), a.String())
					}
				}
			case conditionExpEquivalent:
				// harder to check for without visiting other cells
			case conditionExpDistinct:
				// harder to check for without visiting other cells
			default:
				panic("unhandled condition expression type")
			}
		}
	}

	a.blacklistedDominoIDs = &invalidDominoes
	return a
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
func GetDominoArrangements(game *Game, outArrangements chan<- DominoArrangement) {
	if game == nil {
		panic("nil board")
	}
	debugPrint(fmt.Println, strings.Repeat("*", 64))
	defer debugPrint(fmt.Println, strings.Repeat("*", 64)+"\n")
	debugPrint(fmt.Println, "Calculating possible arrangements for dominoes on the board...")

	// create a map of cells in play to keep track of
	cellsRemaining := branchUnarrangedCells(game.inPlayCellsByIdentifier)

	locations := make([]DominoArrangementLocation, 0) // tracks locations of fitted dominoes for a possible arrangement

	findDominoArrangements(game, cellsRemaining, locations, outArrangements)
}

// attempts to recurse through different ways of fitting dominoes to a board without using loops
// each recursive call will fit a domino into a cell and one of its neighbors, then remove the two from the remaining cells
func findDominoArrangements(game *Game, unarrangedCells map[string]*cell, locations []DominoArrangementLocation, outArrangements chan<- DominoArrangement) {
	if game == nil {
		panic("nil board")
	}
	if outArrangements == nil {
		panic("nil output arrangements")
	}

	// base case - all cells have been accounted for in the arrangement, so save it
	if len(unarrangedCells) == 0 {
		newSolution := DominoArrangement{
			locations: locations,
		}
		outArrangements <- newSolution
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
		addedLocation := &DominoArrangementLocation{
			cell1: nextCell.identifier(),
			cell2: neighbor.identifier(),
		}
		locationsNew = append(locationsNew, *addedLocation.addBlacklistedDominoIDs(game))

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
