package main

var board = InitializeBoard([][]CellType{
	{CellInPlay, CellInPlay, CellInPlay, CellInPlay},
	{CellEmpty, CellInPlay, CellInPlay, CellEmpty},
	{CellEmpty, CellInPlay, CellInPlay, CellEmpty},
})

func main() {
	board.Print()

	board.WillDominoesFit()
}
