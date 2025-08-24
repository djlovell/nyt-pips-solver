package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Cell struct {
	// whether or not the grid cell is part of the game board
	inPlay bool
	// the cell's position in the grid (calculated upon board initialization)
	posX, posY int
	// neighbor cells - will be nil if neighbor is unused
	neighborLeft, neighborAbove, neighborRight, neighborBelow *Cell
}

// unique identifier for a cell based on its position
func (c Cell) Identifier() string {
	return BoardPosToCellIdentifier(c.posX, c.posY)
}

// parses a Cell from an input specification
func parseInputCell(s string) (*Cell, error) {
	switch s {
	case "X":
		return &Cell{inPlay: false}, nil
	case "O":
		return &Cell{inPlay: true}, nil
	default:
		return nil, fmt.Errorf("%s is an unknown input cell type", s)
	}
}

// recovers a cell identifier from x/y positions on a board
func BoardPosToCellIdentifier(posX, posY int) string {
	return strconv.Itoa(posX) + ":" + strconv.Itoa(posY)
}

// unpacks a cell identifier to x/y positions, also can be used to validate an alleged cell identifier "X:Y"
func CellIdentifierToBoardPos(s string) (int, int, error) {
	coords := strings.Split(s, ":")
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
