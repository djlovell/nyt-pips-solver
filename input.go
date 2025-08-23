package main

/*
Input Specification

When describing game boards, a grid system is used.
The locations of cells, dominoes, conditions, etc. are consistently described and will be referred to
by their (x/column, y/row) position, zero-indexed and increasing down and to the right when visualizing the game board.

When inputting the board, a uniform grid (consistent row widths and column heights) is expected.
Cells used in the game (aka that will have a domino placed on them) will be described as:
- "O" for in-play
- "X" if they are not part of the game

"Conditions" in the game apply to groups of cells and can be any of the below:
- "sum_equal"
- "N", where N is a number that domino faces covering the group of cells must add up to
- "=", where the domino faces must all be equivalent
- "!=", where the domino faces must all be distinct
- ">N", where the domino faces must be cumulatively greater than N
- "<N", where the domino faces must be cumulatively less than N



Not every cell needs to be covered by a condition, and conditions may overlap across cells
*/

// in-game tutorial (at least as of 8/22)
var inputCells = [][]string{
	{"X", "O", "O"},
	{"X", "O", "O"},
	{"O", "O", "X"},
}

// 8/18/25 easy puzzle
// var inputCells = [][]string{
// 	{"O", "O", "O", "O"},
// 	{"X", "O", "O", "X"},
// 	{"X", "O", "O", "X"},
// }

// // 8/19/25 hard puzzle
// var inputCells = [][]string{
// 	{"O", "O", "X", "X", "X", "X"},
// 	{"O", "O", "X", "X", "X", "X"},
// 	{"O", "O", "O", "O", "X", "O"},
// 	{"X", "O", "O", "O", "O", "O"},
// 	{"X", "X", "O", "O", "O", "O"},
// }

// 8/20/25 hard puzzle
// var inputCells = [][]string{
// 	{"O", "X", "X", "X", "X", "X"},
// 	{"O", "X", "X", "X", "X", "X"},
// 	{"O", "X", "X", "X", "X", "X"},
// 	{"O", "X", "X", "X", "X", "X"},
// 	{"O", "O", "O", "O", "O", "O"},
// 	{"O", "X", "X", "X", "X", "O"},
// 	{"O", "X", "X", "X", "X", "O"},
// 	{"O", "X", "X", "X", "X", "O"},
// 	{"O", "X", "X", "X", "X", "O"},
// }
