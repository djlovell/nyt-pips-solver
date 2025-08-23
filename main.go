package main

func main() {
	board := InitializeBoard(inputCells)
	board.Print()
	board.DominoFitCheck()
}
