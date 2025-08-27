package solver

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type cell struct {
	// whether or not the grid cell is part of the game board
	inPlay bool
	// the cell's position in the grid (calculated upon board initialization)
	posX, posY int
	// neighbor cells - will be nil if neighbor is unused
	neighborLeft, neighborAbove, neighborRight, neighborBelow *cell
	// applicable conditions - this might be useful later for preemptively eliminating invalid solutions early in domino placement
	// e.g. if the first domino is placed such that a "5" is on a cell with condition "sum < 5" or "4", that solution-path can
	// be early terminated
	applicableConditions []*condition
}

// unique identifier for a cell based on its position
func (c cell) identifier() string {
	return boardPosToCellIdentifier(c.posX, c.posY)
}

// recovers a cell identifier from x/y positions on a board
func boardPosToCellIdentifier(posX, posY int) string {
	return strconv.Itoa(posX) + ":" + strconv.Itoa(posY)
}

// unpacks a cell identifier to x/y positions, also can be used to validate an alleged cell identifier "X:Y"
func cellIdentifierToBoardPos(s string) (int, int, error) {
	coords := strings.Split(s, ":")
	if len(coords) != 2 {
		return 0, 0, fmt.Errorf("failed to parse cell identifier - %s is not a valid format", s)
	}
	xPos, err := strconv.Atoi(coords[0])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse cell identifier x position - %w", err)
	}
	yPos, err := strconv.Atoi(coords[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse cell identifier y position - %w", err)
	}
	if xPos < 0 || yPos < 0 {
		return 0, 0, errors.New("failed to parse cell identifier - positions cannot be negative")
	}
	return xPos, yPos, nil
}

// parses a Cell from an input specification
func parseInputCell(s string) (*cell, error) {
	switch s {
	case "X":
		return &cell{inPlay: false}, nil
	case "O":
		return &cell{inPlay: true}, nil
	default:
		return nil, fmt.Errorf("%s is an unknown input cell type", s)
	}
}
