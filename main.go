package main

import (
	"djlovell/nyt_pips_solver/solver"
	"fmt"
	"strings"
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
	fmt.Println("Testing out domino placements...")
	fmt.Println()
	for _, a := range dominoArrangements {
		possibleSolutions, err := solver.GetPossibleSolutionsForArrangement(board, &a)
		if err != nil {
			panic(err)
		}
		for _, s := range possibleSolutions {
			correct, err := solver.CheckSolution(board, &s)
			if err != nil {
				panic(err)
			}
			if correct {
				fmt.Println(strings.Repeat("*", 64))
				defer fmt.Println(strings.Repeat("*", 64) + "\n")
				fmt.Println("FOUND A VALID SOLUTION! GO TRY IT ON THE NYT GAMES APP/SITE.")
				fmt.Println(s.String())
				return
			}
		}
	}
	fmt.Println("FAILED TO SOLVE...DO IT YOURSELF!")
}
