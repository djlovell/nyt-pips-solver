package main

// 8/18/25 easy puzzle
// var boardSpec = [][]CellType{
// 	{CellInPlay, CellInPlay, CellInPlay, CellInPlay},
// 	{CellEmpty, CellInPlay, CellInPlay, CellEmpty},
// 	{CellEmpty, CellInPlay, CellInPlay, CellEmpty},
// }

// // 8/19/25 hard puzzle
// var boardSpec = [][]CellType{
// 	{CellInPlay, CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty},
// 	{CellInPlay, CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty},
// 	{CellInPlay, CellInPlay, CellInPlay, CellInPlay, CellEmpty, CellInPlay},
// 	{CellEmpty, CellInPlay, CellInPlay, CellInPlay, CellInPlay, CellInPlay},
// 	{CellEmpty, CellEmpty, CellInPlay, CellInPlay, CellInPlay, CellInPlay},
// }

// 8/20/25 hard puzzle
var boardSpec = [][]CellType{
	{CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty, CellEmpty},
	{CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty, CellEmpty},
	{CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty, CellEmpty},
	{CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty, CellEmpty},
	{CellInPlay, CellInPlay, CellInPlay, CellInPlay, CellInPlay, CellInPlay},
	{CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty, CellInPlay},
	{CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty, CellInPlay},
	{CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty, CellInPlay},
	{CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty, CellInPlay},
}

func main() {
	board := InitializeBoard(boardSpec)
	board.Print()
	board.DominoFitCheck()
}
