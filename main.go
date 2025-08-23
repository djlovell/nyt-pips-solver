package main

func main() {
	board := InitializeBoard(inputCells)
	_ = inputConditions // TODO: parse these, not sure about whether to make them board attributes, or an overlaid object
	_ = inputDominoes   // TODO: parse these
	board.Print()
	board.DominoFitCheck()
}
