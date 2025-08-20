package main

/*
Board Specifications

When specifying game boards, a grid system is used.
The locations of cells, dominoes, conditions, etc. are consistently described and will be referred to
by their (x/column, y/row) position, zero-indexed and increasing down and to the right when visualizing the game board.

Cells used in the game (aka that will have a domino placed on them) will be described as "in play", and those not
part of the game "empty".

"Conditions" in the game apply to groups of cells and can be any of the following:
- "N", where N is a number that domino faces covering the group of cells must add up to
- "=", where the domino faces must all be equivalent
- "!=", where the domino faces must all be distinct
- ">N", where the domino faces must be cumulatively greater than N
- "<N", where the domino faces must be cumulatively less than N

Not every cell needs to be covered by a condition, and conditions may overlap across cells
*/

// 8/18/25 easy puzzle
// var boardSpec = [][]CellType{
// 	{CellInPlay, CellInPlay, CellInPlay, CellInPlay},
// 	{CellEmpty, CellInPlay, CellInPlay, CellEmpty},
// 	{CellEmpty, CellInPlay, CellInPlay, CellEmpty},
// }

// // 8/19/25 hard puzzle
var boardSpec = [][]CellType{
	{CellInPlay, CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty},
	{CellInPlay, CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty},
	{CellInPlay, CellInPlay, CellInPlay, CellInPlay, CellEmpty, CellInPlay},
	{CellEmpty, CellInPlay, CellInPlay, CellInPlay, CellInPlay, CellInPlay},
	{CellEmpty, CellEmpty, CellInPlay, CellInPlay, CellInPlay, CellInPlay},
}

// 8/20/25 hard puzzle
// var boardSpec = [][]CellType{
// 	{CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty, CellEmpty},
// 	{CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty, CellEmpty},
// 	{CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty, CellEmpty},
// 	{CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty, CellEmpty},
// 	{CellInPlay, CellInPlay, CellInPlay, CellInPlay, CellInPlay, CellInPlay},
// 	{CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty, CellInPlay},
// 	{CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty, CellInPlay},
// 	{CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty, CellInPlay},
// 	{CellInPlay, CellEmpty, CellEmpty, CellEmpty, CellEmpty, CellInPlay},
// }

func main() {
	board := InitializeBoard(boardSpec)
	board.Print()
	board.DominoFitCheck()
}
