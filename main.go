package main

import (
	"djlovell/nyt_pips_solver/solver"
)

func main() {
	// load the game board from input
	board, err := solver.LoadGame(inputCells, inputConditions, inputDominoes)
	if err != nil {
		panic(err)
	}
	board.Print()

	// calculate possible ways dominoes can fit on the board
	dominoArrangements, err := solver.GetDominoArrangements(board)
	if err != nil {
		panic(err)
	}

	// try to solve each arrangement
	for _, a := range dominoArrangements {
		solver.TryArrangement(board, &a)
	}
}
