package main

import (
	"djlovell/nyt_pips_solver/input"
	"djlovell/nyt_pips_solver/solver"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	// define CL argument for specifying input json file
	inputFilename := flag.String("f", "", "Input file (JSON)")
	verbose := flag.Bool("v", false, "Enable debug output (it's not gonna be pretty...)")

	flag.Parse()
	if inputFilename == nil {
		panic("input file name flag should have at least defaulted to empty")
	}
	if verbose == nil {
		panic("verbose flag should have defaulted to something")
	}

	if !strings.HasSuffix(*inputFilename, ".json") {
		fmt.Println("Error: input file should be of the format *.json")
		fmt.Println(*inputFilename)
		return
	}
	solver.SetDebugPrint(*verbose)

	// load the game input file
	if _, err := os.Stat(*inputFilename); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Error: input file not found")
		return
	}
	inputGame, err := input.ReadFile(*inputFilename)
	if err != nil {
		fmt.Printf("Error: input file read failed with error - %s\n", err.Error())
		return
	}

	game, err := solver.ParseInputGame(inputGame)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}
	game.Print()

	// start a timer for solving
	startTime := time.Now()

	// calculate possible ways dominoes can fit on the game board
	fmt.Println("Calculating possible domino arrangements...")
	fmt.Println()
	dominoArrangementChan := make(chan solver.DominoArrangement)
	{
		wg := new(sync.WaitGroup)
		wg.Go(func() {
			solver.GetDominoArrangements(game, dominoArrangementChan)
		})
		go func() {
			wg.Wait()
			close(dominoArrangementChan)
		}()
	}

	// get possible solutions for each arrangement
	fmt.Println("Calculating possible solutions...")
	fmt.Println()
	possibleSolutionChan := make(chan solver.Solution)
	{
		// calculate possible solutions for each arrangement in parallel
		wg := new(sync.WaitGroup)
		for a := range dominoArrangementChan {
			wg.Go(func() {
				solver.GetPossibleSolutionsForArrangement(game, &a, possibleSolutionChan)
			})
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
			wg.Go(func() {
				for s := range possibleSolutionChan {
					if correct := solver.CheckSolution(game, &s); correct {
						validSolutionChan <- s
					}
				}
			})
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
	fmt.Printf("NYT Pips Solver Completed in %f seconds. ", time.Since(startTime).Seconds())
	switch l := len(validSolutions); l {
	case 0:
		fmt.Println("No valid solutions found (RIP).")
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
