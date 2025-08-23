package main

import (
	"strconv"
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
	return strconv.Itoa(c.posX) + ":" + strconv.Itoa(c.posY)
}

// parses a Cell from an input specification
func parseInputCell(s string) *Cell {
	switch s {
	case "X":
		return &Cell{inPlay: false}
	case "O":
		return &Cell{inPlay: true}
	default:
		panic("unknown cell type")
	}
}
