package solver

import (
	"errors"
	"fmt"
	"strings"
)

type Game struct {
	board      [][]*cell
	conditions []*condition
	dominoes   []*domino
	// helpers for solving
	inPlayCellsByIdentifier map[string]*cell
}

// LoadGame - loads a game board from input
func LoadGame(
	inputCells [][]string,
	inputConditions [][]string,
	inputDominoes []string,
) (*Game, error) {
	game := new(Game)
	game.inPlayCellsByIdentifier = make(map[string]*cell)

	// cell initialization
	{
		board := make([][]*cell, 0)
		gridWidth := -1 // set by first row, then used to make sure rows are fixed-width
		for _, inputRow := range inputCells {
			if gridWidth == -1 {
				gridWidth = len(inputRow)
			}
			if len(inputRow) != gridWidth {
				return nil, errors.New("input cell grid is not a consistent width")
			}
			cellRow := make([]*cell, 0)
			for _, c := range inputRow {
				if p, err := parseInputCell(c); err != nil {
					return nil, err
				} else {
					cellRow = append(cellRow, p)
				}
			}
			board = append(board, cellRow)
		}

		// fill in cell positions and establish neighbors starting from the bottom/right of the board
		for yIdx := len(board) - 1; yIdx >= 0; yIdx-- {
			for xIdx := len(board[yIdx]) - 1; xIdx >= 0; xIdx-- {
				cell := board[yIdx][xIdx]
				if cell == nil {
					return nil, errors.New("nil cell found during initialization")
				}

				// set the position of the current cell
				cell.posX = xIdx
				cell.posY = yIdx

				// unused cells do not get neighbors nor get to be neighbors
				if !cell.inPlay {
					continue
				}

				// establish neighbor associations with the neighbor to the left
				if xIdx > 0 {
					neighborLeft := board[yIdx][xIdx-1]
					if neighborLeft == nil {
						return nil, errors.New("nil cell found during initialization")
					}
					if neighborLeft.inPlay {
						cell.neighborLeft = neighborLeft
						neighborLeft.neighborRight = cell
					}
				}

				// establish neighbor associations with the neighbor above
				if yIdx > 0 {
					neighborAbove := board[yIdx-1][xIdx]
					if neighborAbove == nil {
						return nil, errors.New("nil cell found during initialization")
					}
					if neighborAbove.inPlay {
						cell.neighborAbove = neighborAbove
						neighborAbove.neighborBelow = cell
					}
				}

				// make the cell easily accessible by its identifier
				game.inPlayCellsByIdentifier[cell.identifier()] = cell
			}
		}
		game.board = board
	}

	// condition initialization
	{
		conditions := make([]*condition, 0)
		for _, inputCond := range inputConditions {
			condition, err := parseInputCondition(inputCond)
			if err != nil {
				return nil, err
			}
			// fact check that all the locations are actually valid cells
			for _, conditionCellLoc := range condition.cellIdentifiers {
				xPos, yPos, err := cellIdentifierToBoardPos(conditionCellLoc)
				if err != nil {
					return nil, err
				}
				if yPos >= len(game.board) {
					return nil, errors.New("input condition cell location out of y range")
				}
				if xPos >= len(game.board[yPos]) {
					return nil, errors.New("input condition cell location out of x range")
				}
				cell := game.board[yPos][xPos]
				if !cell.inPlay {
					return nil, fmt.Errorf("input condition contains cell %s not in play", cell.identifier())
				}
				// link condition to cell (in case this is needed/useful)
				cell.applicableConditions = append(cell.applicableConditions, condition)
			}
			conditions = append(conditions, condition)
		}
		game.conditions = conditions
	}

	// domino initialization
	{
		dominoes := make([]*domino, 0)
		for _, inputDomino := range inputDominoes {
			domino, err := parseInputDomino(inputDomino)
			if err != nil {
				return nil, err
			}
			// make the domino easily accessible by its identifier
			dominoes = append(dominoes, domino)
		}
		game.dominoes = dominoes
	}

	return game, nil
}

// I hate it but this is my confirmation that input parsing worked for now
func (b Game) Print() {
	// pretty print the board
	xIdxMax := 0
	fmt.Println(strings.Repeat("*", 64))
	fmt.Println("This is kinda what the board looks like...")
	rowStrings := []string{}
	for yIdx, r := range b.board {
		rowCells := []string{fmt.Sprintf("%-3d", yIdx)} // TODO: make Y index pad up to 3 digits
		for xIdx, c := range r {
			cellStr := "["
			if c.inPlay {
				cellStr += " "
			} else {
				cellStr += "X"
			}
			cellStr += "]"
			rowCells = append(rowCells, cellStr)
			xIdxMax = max(xIdxMax, xIdx)
		}
		rowStrings = append(rowStrings, strings.Join(rowCells, " "))
	}
	headerRowValues := []string{"   "}
	for i := 0; i <= xIdxMax; i++ {
		headerRowValues = append(headerRowValues, fmt.Sprintf("%-3d", i))
	}
	rowStrings = append([]string{strings.Join(headerRowValues, " ")}, rowStrings...)
	fmt.Println(strings.Join(rowStrings, "\n") + "\n")

	// print out the conditions in english
	fmt.Println("Conditions:")
	for _, c := range b.conditions {
		fmt.Printf("  %s\n", c.String())
	}
	fmt.Println()

	// print out the dominoes
	fmt.Println("Dominoes:")
	for _, d := range b.dominoes {
		fmt.Printf("  %s\n", d.String())
	}

	fmt.Println(strings.Repeat("*", 64))
	fmt.Println()
}
