package solver

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"
)

// DominoPlacement - the specific location and orientation of a domino in a Solution
type DominoPlacement struct {
	cell1Identifier string
	cell1Value      int
	cell2Identifier string
	cell2Value      int
	// string for pretty printing the domino
	printString string
}

func (p DominoPlacement) String() string {
	return fmt.Sprintf(
		"Domino %s placed with %d in Cell %s & %d in Cell %s\n",
		p.printString, p.cell1Value, p.cell1Identifier, p.cell2Value, p.cell2Identifier,
	)
}

// Solution - a complete layout of dominoes on the board
type Solution struct {
	dominoPlacements []DominoPlacement
}

func (s Solution) String() string {
	out := "Solution\n"
	for _, p := range s.dominoPlacements {
		out += "  " + p.String()
	}
	return out
}

func getCellValuesFromPlacements(placements *[]DominoPlacement) map[string] /* cell identifier */ int /*cell value */ {
	if placements == nil {
		panic("nil placements")
	}
	m := make(map[string]int)
	for _, p := range *placements {
		m[p.cell1Identifier] = p.cell1Value
		m[p.cell2Identifier] = p.cell2Value
	}
	return m
}

// GetPossibleSolutionsForArrangement - finds different potential solutions to check.
// Ideally, a lot of incorrect solution paths are eliminated here with early condition checks.
func GetPossibleSolutionsForArrangement(game *Game, dominoArrangement *DominoArrangement, outPossibleSolutions chan<- Solution) {
	if game == nil {
		panic("nil game")
	}
	if dominoArrangement == nil {
		panic("nil arrangement")
	}
	debugPrint(fmt.Println, strings.Repeat("*", 64))
	defer debugPrint(fmt.Println, strings.Repeat("*", 64)+"\n")
	debugPrint(fmt.Println, "Calculating possible solutions using arrangement...")
	debugPrint(fmt.Println, dominoArrangement.String())

	// track locations that have not been filled with a domino yet
	unfilledLocations := make([]DominoArrangementLocation, len(dominoArrangement.locations))
	copy(unfilledLocations, dominoArrangement.locations)

	// sort so the location with the fewest possibilities is evaluated first, reducing
	// the breadth of the decision tree
	//
	// for one puzzle, starting with the least # of dominoes vs. most reduced solve time from 90s to 6s
	slices.SortFunc(unfilledLocations, func(l, r DominoArrangementLocation) int {
		return len(*r.blacklistedDominoIDs) - len(*l.blacklistedDominoIDs)
	})

	// track unplaced and placed dominoes as time progresses
	unplacedDominoes := make(map[string]*domino)
	for _, d := range game.dominoes {
		unplacedDominoes[d.identifier] = d
	}
	placementsSoFar := make([]DominoPlacement, 0)

	// start placing dominoes
	placeDomino(game, unfilledLocations, unplacedDominoes, placementsSoFar, outPossibleSolutions)
}

// CheckSolution - returns if a possible solution successfully met all the conditions to solve the puzzle
func CheckSolution(game *Game, solution *Solution) bool {
	if game == nil {
		panic("nil game")
	}
	if solution == nil {
		panic("nil solution")
	}
	debugPrint(fmt.Println, strings.Repeat("*", 64))
	defer debugPrint(fmt.Println, strings.Repeat("*", 64)+"\n")
	debugPrint(fmt.Println, "Checking possible solution...")
	debugPrint(fmt.Println, solution.String())

	// check each condition, early returning if one fails
	cellValues := getCellValuesFromPlacements(&solution.dominoPlacements)
	for _, cond := range game.conditions {
		if ok, err := cond.check(cellValues); err != nil {
			panic("very unexpected error checking condition")
		} else if !ok {
			debugPrint(fmt.Printf, `Solution violates "%s"`+"\n", cond.String())
			return false
		}
	}

	debugPrint(fmt.Println, "Solution appears valid!")
	return true
}

// recursively places dominoes on the game board, testing along the way until a solution is reached
func placeDomino(
	game *Game,
	unfilledLocations []DominoArrangementLocation,
	unplacedDominoes map[string]*domino,
	placementsSoFar []DominoPlacement,
	outPossibleSolutions chan<- Solution,
) {
	if game == nil {
		panic("nil board")
	}
	if len(unfilledLocations) != len(unplacedDominoes) {
		panic("mismatch between number of dominoes and places to put them")
	}

	// base case - all locations have been filled with a
	if len(unfilledLocations) == 0 {
		// need to do this copy to prevent backtracking bugs (if using)
		placementsCopy := make([]DominoPlacement, len(placementsSoFar))
		copy(placementsCopy, placementsSoFar)
		newSolution := Solution{
			dominoPlacements: placementsCopy,
		}
		outPossibleSolutions <- newSolution
		debugPrint(fmt.Println, "All dominoes placed and possible solution added...")
		return
	}

	// check for violated conditions before moving on
	{
		cellValuesSoFar := getCellValuesFromPlacements(&placementsSoFar)
		for _, cond := range game.conditions {
			ok, err := cond.check(cellValuesSoFar)
			if err != nil {
				if errors.Is(err, errConditionNotReadyToCheck) {
					// condition just isn't ready to evaluate yet - move on with placing a domino
					continue
				}
				panic("very unexpected error checking condition" + err.Error())
			}
			if !ok {
				// condition failed! abort this path
				return
			}

		}
	}

	nextLocation := unfilledLocations[0]
	remainingLocations := make([]DominoArrangementLocation, len(unfilledLocations[1:]))
	copy(remainingLocations, unfilledLocations[1:])

	// try all dominoes
	for _, nextDomino := range slices.Collect(maps.Values(unplacedDominoes)) {
		// skip the domino if it is blacklisted for the location
		if _, blacklisted := (*nextLocation.blacklistedDominoIDs)[nextDomino.identifier]; blacklisted {
			continue
		}

		// try both orientations if they are different too
		type orientation struct {
			cell1Val, cell2Val int
		}
		orientations := []orientation{
			{
				cell1Val: nextDomino.val1,
				cell2Val: nextDomino.val2,
			},
		}
		if nextDomino.val1 != nextDomino.val2 {
			orientations = append(orientations, orientation{
				cell1Val: nextDomino.val2,
				cell2Val: nextDomino.val1,
			})
		} else {
			debugPrint(fmt.Printf, "not checking reverse orientation of domino %s in location %s...\n", nextDomino.String(), nextLocation.String())
		}

		for i, o := range orientations {
			if i == 0 {
				debugPrint(fmt.Printf, "placing domino %s in location %s...\n", nextDomino.String(), nextLocation.String())
			} else {
				debugPrint(fmt.Printf, "placing domino %s in location %s (reverse orientation)...\n", nextDomino.String(), nextLocation.String())
			}

			// generate the next placement
			placement := &DominoPlacement{
				cell1Identifier: nextLocation.cell1,
				cell1Value:      o.cell1Val,
				cell2Identifier: nextLocation.cell2,
				cell2Value:      o.cell2Val,
				printString:     nextDomino.String(),
			}

			// remove the domino since it will have been placed
			delete(unplacedDominoes, nextDomino.identifier)

			// track the placement
			placementsSoFar = append(placementsSoFar, *placement)

			// perform the next placement recursively (concurrently, if still allowed)
			placeDomino(game, remainingLocations, unplacedDominoes, placementsSoFar, outPossibleSolutions)

			// backtrack
			placementsSoFar = placementsSoFar[0 : len(placementsSoFar)-1]
			unplacedDominoes[nextDomino.identifier] = nextDomino
		}
	}
}
