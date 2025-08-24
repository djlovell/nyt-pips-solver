package main

import (
	"fmt"
)

func main() {
	board, err := InitializeBoard(inputCells, inputConditions)
	if err != nil {
		panic(err)
	}
	_ = inputDominoes // TODO: parse these
	board.Print()

	// debug - verify that dominoes can fit on the board, and possible orientations
	if err := GetDominoArrangements(board); err != nil {
		fmt.Println(err)
		return
	}
}
