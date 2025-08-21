package main

func main() {
	board := InitializeBoard(boardSpec)
	board.Print()
	board.DominoFitCheck()
}
