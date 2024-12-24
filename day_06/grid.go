package main

import (
	"slices"
)

type CellType rune

const (
	CellTypeOpen    CellType = '.'
	CellTypeBlocked          = '#'
	CellTypeGuard            = '^'
)

type Grid struct {
	width  int
	height int
	cells  []CellType
}

func (g *Grid) Clone() Grid {
	return Grid{
		width:  g.width,
		height: g.height,
		cells:  slices.Clone(g.cells),
	}
}

func (g *Grid) WithinBounds(loc Vector) bool {
	return loc.X < g.width &&
		loc.Y < g.height &&
		loc.X >= 0 &&
		loc.Y >= 0
}

func (g *Grid) CellAt(pos Vector) CellType {
	return g.cells[pos.Y*g.width+pos.X]
}

func (g *Grid) SetCellAt(pos Vector, typ CellType) {
	g.cells[pos.Y*g.width+pos.X] = typ
}

func (g *Grid) CellAtPositionWalkable(pos Vector) bool {
	return g.CellAt(pos) == CellTypeOpen
}
