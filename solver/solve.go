package solver

import (
	"fmt"
	"maps"
	"strings"
	"sync"
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

func (s *Solution) getCellValues() map[string] /* cell identifier */ int /*cell value */ {
	m := make(map[string]int)
	for _, p := range s.dominoPlacements {
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
	unfilledLocations := dominoArrangement.locations

	// track unplaced and placed dominoes as time progresses
	unplacedDominoes := make(map[string]*domino)
	for _, d := range game.dominoes {
		unplacedDominoes[d.identifier] = d
	}
	placementsSoFar := make([]DominoPlacement, 0)

	// we can allow domino placement paths to be explored in parallel up to a certain depth without thrashing
	// set to 10, I was able to observe 100% CPU/RAM usage and 100% sustained SSD activity...
	//
	// I started seeing negative return after 4 layers
	maxConcurrentLayers := 4
	placeDomino(game, unfilledLocations, unplacedDominoes, placementsSoFar, outPossibleSolutions, maxConcurrentLayers)
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
	cellValues := solution.getCellValues()
	for _, cond := range game.conditions {
		if ok := cond.check(cellValues); !ok {
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
	concurrentLayersRemaining int,
) {
	if game == nil {
		panic("nil board")
	}
	if len(unfilledLocations) != len(unplacedDominoes) {
		panic("mismatch between number of dominoes and places to put them")
	}

	// base case - all locations have been filled with a
	if len(unfilledLocations) == 0 {
		newSolution := Solution{
			dominoPlacements: placementsSoFar,
		}
		// TODO: test conditions at the end to see if we succeeded
		outPossibleSolutions <- newSolution
		debugPrint(fmt.Println, "All dominoes placed and possible solution added...")
		return
	}

	// get the next location to fill
	nextLocation := unfilledLocations[0]
	remainingLocations := unfilledLocations[1:]

	// try all dominoes
	wg := new(sync.WaitGroup)
	for _, nextDomino := range unplacedDominoes {
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
			// pre-check this placement to see if it violates any conditions
			if fail := placement.doesPlacementFailConditionsEarly(game); fail {
				continue
			}

			// remove the domino since it will have been placed
			unplacedDominoesNew := branchUnplacedDominoes(unplacedDominoes)
			delete(unplacedDominoesNew, nextDomino.identifier)

			// track the placement
			placementsSoFarNew := branchPlacementsSoFar(placementsSoFar)
			placementsSoFarNew = append(placementsSoFarNew, *placement)

			// perform the next placement recursively (concurrently, if still allowed)
			nextRecursiveCall := func() {
				placeDomino(game, remainingLocations, unplacedDominoesNew, placementsSoFarNew, outPossibleSolutions, concurrentLayersRemaining-1)
			}
			if concurrentLayersRemaining > 0 {
				wg.Go(nextRecursiveCall)
			} else {
				nextRecursiveCall()
			}
		}
	}
	wg.Wait()
}

// TODO: this really belongs in condition.go, and could use more thought/refinement,
// but serves the purpose of early-eliminating solution paths where a placement
// single-handedly violates a condition. e.g. placing a 6 in a cell with "<5" condition.
func (p DominoPlacement) doesPlacementFailConditionsEarly(g *Game) bool {
	if g == nil {
		panic("nil board")
	}
	cell1, ok := g.inPlayCellsByIdentifier[p.cell1Identifier]
	if !ok {
		panic("cell identifier not found in in play cells by identifier - check game loading")
	}
	cell2, ok := g.inPlayCellsByIdentifier[p.cell2Identifier]
	if !ok {
		panic("cell identifier not found in in play cells by identifier - check game loading")
	}

	for _, c := range []struct {
		cell *cell
		val  int
	}{
		{
			cell: cell1,
			val:  p.cell1Value,
		},
		{
			cell: cell2,
			val:  p.cell2Value,
		},
	} {
		for _, cond := range c.cell.applicableConditions {
			switch cond.expression {
			case conditionExpSumEquals:
				// fail early if cell by itself exceeds the collective sum
				if c.val > cond.operand {
					debugPrint(fmt.Printf, `Value %d in Cell %s violates "%s"`, c.val, c.cell.identifier(), cond.String())
					debugPrint(fmt.Println)
					return true
				}
			case conditionExpSumLessThan:
				// fail early if cell by itself meets or exceeds collective sum
				if c.val >= cond.operand {
					debugPrint(fmt.Printf, `Value %d in Cell %s violates "%s"`, c.val, c.cell.identifier(), cond.String())
					debugPrint(fmt.Println)
					return true
				}
			case conditionExpSumGreaterThan:
				// harder to check for without visiting other cells
			case conditionExpEquivalent:
				// harder to check for without visiting other cells
			case conditionExpDistinct:
				// harder to check for without visiting other cells
			default:
				panic("unhandled condition expression type")
			}
		}
	}

	return false
}

func branchUnplacedDominoes(in map[string]*domino) map[string]*domino {
	out := make(map[string]*domino)
	maps.Copy(out, in)
	return out
}

func branchPlacementsSoFar(in []DominoPlacement) []DominoPlacement {
	out := make([]DominoPlacement, len(in))
	copy(out, in)
	return out
}
