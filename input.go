package main

/*
Input Specification

When describing game boards, a grid system is used.
The locations of cells, dominoes, conditions, etc. are consistently described and will be referred to
by their (x/column, y/row) position, zero-indexed and increasing down and to the right when visualizing the game board.

When inputting the board, a uniform grid (consistent row widths and column heights) is expected, and shall
be a [][]string.

Cells used in the game (aka that will have a domino placed on them) will be described as:
- "O" for in-play
- "X" if they are not part of the game

"Conditions" in the game apply to groups of cells and can be any of the following expressions:
- "N", where N is a number that domino faces covering the group of cells must add up to
- "=", where the domino faces must all be equivalent
- "!=", where the domino faces must all be distinct
- ">N", where the domino faces must be cumulatively greater than N
- "<N", where the domino faces must be cumulatively less than N

When inputting, conditions shall be a [][]string, each []string condition following the format:
{"{{expression from above}}", "{{integer operand for expressions that have N in them}}", "X1:Y1", "X2:Y2", etc. for any other cells the condition covers}
Example: {"N", "12", "0:0", "0:1"} implies that the first two cells in the top left corner of the board must add up to 12

Not every cell needs to be covered by a condition, and conditions may overlap across cells

Dominoes shall be input as a []string, each string of the format "N|M", with no particular polarity, but where 0 <= N,M <= 6.
*/

// in-game tutorial (at least as of 8/22) - confirmed working!
// var inputCells = [][]string{
// 	{"X", "O", "O"},
// 	{"X", "O", "O"},
// 	{"O", "O", "X"},
// }

// var inputConditions = [][]string{
// 	{"N", "0", "0:2"},
// 	{"=", "1:1", "1:2"},
// 	{"N", "10", "2:0", "2:1"},
// }

// var inputDominoes = []string{
// 	"5|5", "0|2", "2|3",
// }

// medium puzzle 8/26/2025 - confirmed working!
// var inputCells = [][]string{
// 	{"O", "X", "X", "X"},
// 	{"O", "X", "X", "X"},
// 	{"O", "X", "X", "X"},
// 	{"O", "X", "X", "X"},
// 	{"O", "O", "O", "O"},
// 	{"O", "X", "X", "O"},
// 	{"O", "X", "X", "O"},
// 	{"O", "X", "X", "O"},
// }

// var inputConditions = [][]string{
// 	{"N", "10", "0:1", "0:2"},
// 	{"N", "1", "0:3", "0:4"},
// 	{"=", "0:6", "0:7"},
// 	{"N", "12", "2:4", "3:4"},
// 	{"N", "8", "3:5", "3:6"},
// }

// var inputDominoes = []string{
// 	"3|5", "5|1", "0|3", "2|2", "2|6", "6|4", "4|0",
// }

// hard puzzle 8/26/2025 - confirmed working!
var inputCells = [][]string{
	{"O", "O", "X", "X", "X", "X"},
	{"O", "O", "X", "X", "X", "X"},
	{"X", "O", "O", "O", "O", "X"},
	{"X", "O", "O", "O", "O", "X"},
	{"X", "O", "X", "X", "O", "X"},
	{"X", "X", "X", "X", "O", "O"},
	{"X", "X", "O", "O", "O", "O"},
}

var inputConditions = [][]string{
	{"N", "10", "0:0", "0:1"},
	{"N", "1", "1:0", "1:1"},
	{"=", "1:2", "1:3", "1:4"},
	{"N", "12", "2:2", "2:3"},
	{"!=", "3:2", "3:3", "4:3"},
	{"=", "4:4", "4:5", "4:6"},
	{"N", "4", "5:5"},
	{"N", "0", "5:6"},
	{"N", "6", "3:6"},
	{"<N", "2", "2:6"},
}

var inputDominoes = []string{
	"6|0", "4|1", "3|6", "5|4", "3|3", "6|4", "2|3", "3|4", "3|0", "1|6",
}
