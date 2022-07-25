package main

import (
	"errors"
	"math/rand"
	"time"
)

type Grid struct {
	size  int
	cells [][]Tile
}

func (g *Grid) setup(preTiles [][]Tile) {
	if len(preTiles) > 0 {
		g.cells = preTiles
	} else {
		g.cells = make([][]Tile, g.size)
		for lx := 0; lx < g.size; lx++ {
			g.cells[lx] = make([]Tile, g.size)

			for ly := 0; ly < g.size; ly++ {
				g.cells[lx][ly] = Tile{x: lx, y: ly, isEmpty: true}
			}
		}
	}
}

func (g *Grid) randomAvailableCell() Tile {
	cells := g.AvailableCells()

	rand.Seed(time.Now().UnixNano())
	if len(cells) > 0 {
		return cells[rand.Int31n(int32(len(cells)))]
	} else {
		return Tile{isEmpty: true}
	}
}

func (g *Grid) insertTile(tile Tile) {
	g.cells[tile.x][tile.y] = tile
}

func (g *Grid) cellsAvailable() bool {
	isAvailable := false
	for _, row := range g.cells {
		for _, item := range row {
			if item.isEmpty {
				isAvailable = item.isEmpty
			}
		}
	}

	return isAvailable
}

func (g *Grid) eachCell(callback func(x int, y int, tile Tile)) {
	for x := 0; x < g.size; x++ {
		for y := 0; y < g.size; y++ {
			callback(x, y, g.cells[x][y])
		}
	}
}

func (g *Grid) CellAvailable(tile *Tile) bool {
	r, error := g.CellContent(tile)

	return error == nil && r.isEmpty
}

func (g *Grid) WithinBounds(position *Tile) bool {
	return position.x >= 0 && position.x < g.size &&
		position.y >= 0 && position.y < g.size
}

func (g *Grid) AvailableCells() []Tile {
	var cells []Tile

	i := 0
	g.eachCell(func(x int, y int, tile Tile) {
		if tile.isEmpty {
			cells = append(cells, tile)
			i += 1
		}
	})

	return cells
}

func (g *Grid) removeTile(tile *Tile) {
	g.cells[tile.x][tile.y].isEmpty = true
}

func (g *Grid) CellContent(cell *Tile) (*Tile, error) {
	if g.WithinBounds(cell) {
		return &g.cells[cell.x][cell.y], nil
	} else {
		return nil, errors.New("The cells doesn't exist")
	}
}
