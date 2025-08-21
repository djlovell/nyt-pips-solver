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

type CellType int

const (
	CellUnused CellType = iota
	CellInPlay
)

func NewCell(t CellType) *Cell {
	switch t {
	case CellUnused:
		return &Cell{inPlay: false}
	case CellInPlay:
		return &Cell{inPlay: true}
	default:
		panic("unknown cell type")
	}
}

// unique identifier for a cell based on its position
func (c Cell) Identifier() string {
	return strconv.Itoa(c.posX) + ":" + strconv.Itoa(c.posY)
}
