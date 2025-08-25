package main

import "djlovell/nyt_pips_solver/solver"

func main() {
	// load the board (cells + conditions) from input
	board, err := solver.InitializeBoard(inputCells, inputConditions)
	if err != nil {
		panic(err)
	}
	board.Print()

	_ = inputDominoes // TODO: load dominoes from input

	// calculate possible ways dominoes can fit on the board
	dominoArrangements, err := solver.GetDominoArrangements(board)
	if err != nil {
		panic(err)
	}

	_ = dominoArrangements
}
