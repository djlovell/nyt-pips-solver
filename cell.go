package main

import (
	"fmt"
	"strconv"
)

type cell struct {
	// whether or not the grid cell is part of the game board
	inPlay bool
	// the cell's position in the grid (calculated upon board initialization)
	posX, posY int
	// neighbor cells - will be nil if neighbor is unused
	neighborLeft, neighborAbove, neighborRight, neighborBelow *cell
}

type CellType int

const (
	CellEmpty CellType = iota
	CellInPlay
)

func NewCell(t CellType) *cell {
	switch t {
	case CellEmpty:
		return &cell{inPlay: false}
	case CellInPlay:
		return &cell{inPlay: true}
	default:
		panic("unknown cell type")
	}
}

func (c cell) String() string {
	s := fmt.Sprintf("(%d,%d)", c.posX, c.posY)
	s += "["
	if !c.inPlay {
		s += "X"
	} else {
		s += " "
	}
	s += "]"
	return s
}

// unique identifier for a cell based on its position
func (c cell) Identifier() string {
	return strconv.Itoa(c.posX) + ":" + strconv.Itoa(c.posY)
}
