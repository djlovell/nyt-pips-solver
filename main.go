package main

import (
	"djlovell/nyt_pips_solver/solver"
	"fmt"
	"strings"
	"sync"
)

func main() {
	// load the game board from input
	board, err := solver.LoadGame(inputCells, inputConditions, inputDominoes)
	if err != nil {
		panic(err)
	}
	board.Print()

	// calculate possible ways dominoes can fit on the board
	fmt.Println("Calculating possible domino arrangements...")
	fmt.Println()
	dominoArrangements, err := solver.GetDominoArrangements(board)
	if err != nil {
		panic(err)
	}

	// get possible solutions for each arrangement
	fmt.Println("Calculating possible solutions...")
	fmt.Println()
	possibleSolutionChan := make(chan solver.Solution)
	{
		// calculate possible solutions for each arrangement in parallel
		wg := new(sync.WaitGroup)
		for _, a := range dominoArrangements {
			wg.Add(1)
			go func() {
				defer wg.Done()
				solver.GetPossibleSolutionsForArrangement(board, &a, possibleSolutionChan)
			}()
		}
		go func() {
			wg.Wait()
			close(possibleSolutionChan)
		}()
	}

	// find valid solutions
	fmt.Println("Testing possible solutions...")
	fmt.Println()
	validSolutionChan := make(chan solver.Solution)
	{
		// use a worker pool to check solutions in parallel
		numCheckers := 100
		wg := new(sync.WaitGroup)
		for range numCheckers {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for s := range possibleSolutionChan {
					correct, err := solver.CheckSolution(board, &s)
					if err != nil {
						panic(err)
					}
					if correct {
						validSolutionChan <- s
					}
				}
			}()
		}
		go func() {
			wg.Wait()
			close(validSolutionChan)
		}()
	}

	validSolutions := make([]solver.Solution, 0)
	for s := range validSolutionChan {
		validSolutions = append(validSolutions, s)
	}

	fmt.Println(strings.Repeat("*", 64))
	defer fmt.Println(strings.Repeat("*", 64))
	fmt.Printf("NYT Pips Solver Completed. ")
	switch l := len(validSolutions); l {
	case 0:
		fmt.Println("No valid soltions found (RIP).")
	case 1:
		fmt.Println("Found a valid solution.\n\nGo try it on the NYT Games app/site!")
		fmt.Println()
	default:
		fmt.Printf("Found %d valid solutions.\n\nGo try them on the NYT Games app/site!\n\n", l)
	}
	for _, s := range validSolutions {
		fmt.Println(s.String())
	}
}
