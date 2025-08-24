package main

import (
	"fmt"
)

func main() {
	board, err := InitializeBoard(inputCells)
	if err != nil {
		panic(err)
	}

	_ = inputConditions // TODO: parse these, not sure about whether to make them board attributes, or an overlaid object
	_ = inputDominoes   // TODO: parse these
	board.Print()

	// debug - verify that dominoes can fit on the board, and possible orientations
	if err := GetDominoArrangements(board); err != nil {
		fmt.Println(err)
		return
	}
}
