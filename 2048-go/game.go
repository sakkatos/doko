package main

import (
	"math/rand"
)

type Game struct {
	gridSize  int
	score     int
	highScore int
	over      bool
	won       bool
	grid      *Grid
	drawer    *Drawer
}

type Vector struct {
	x int
	y int
}

type PositionTraversal struct {
	x []int
	y []int
}

func (g *Game) setup(gameInfo GameInfo) {
	g.grid = &Grid{size: g.gridSize}
	g.grid.setup(gameInfo.getTiles())

	g.highScore = gameInfo.HighScore
	g.score = gameInfo.CurrentScore

	if len(gameInfo.getTiles()) == 0 {
		g.addStartTiles(2)
	}

	add("up", func(message *Message) {
		g.move(0)
	})
	add("right", func(message *Message) {
		g.move(1)
	})
	add("down", func(message *Message) {
		g.move(2)
	})
	add("left", func(message *Message) {
		g.move(3)
	})
}

func (g *Game) addStartTiles(startTileNum int) {
	for i := 0; i < startTileNum; i++ {
		g.addRandomTile()
	}
}

func (g *Game) addRandomTile() {
	if g.grid.cellsAvailable() {
		value := 2
		if rand.Float32() < 0.9 {
			value = 4
		}
		tile := g.grid.randomAvailableCell()
		newTile := Tile{x: tile.x, y: tile.y, value: value, isEmpty: false}

		g.grid.insertTile(newTile)
	}
}

func (g *Game) getVector(direction int) Vector {
	res := make(map[int]Vector)

	res[0] = Vector{x: 0, y: -1} // Up
	res[1] = Vector{x: 1, y: 0}  // Right
	res[2] = Vector{x: 0, y: 1}  // Down
	res[3] = Vector{x: -1, y: 0} // Left

	return res[direction]
}

func (g *Game) moveTile(tile *Tile, farPos *Tile) Tile {
	g.grid.removeTile(tile)
	g.grid.cells[farPos.x][farPos.y] = Tile{x: farPos.x, y: farPos.y, value: tile.value, mergedFrom: tile.mergedFrom, isEmpty: false}
	return g.grid.cells[farPos.x][farPos.y]
}

func (g *Game) IsGameTerminated() bool {
	return false
}

func (g *Game) BuildTraversals(vec Vector) PositionTraversal {
	traversals := PositionTraversal{x: make([]int, g.gridSize), y: make([]int, g.gridSize)}

	for i := 0; i < g.gridSize; i++ {
		traversals.x[i] = i
		traversals.y[i] = i
	}

	if vec.x == 1 {
		ReverseList(traversals.x)
	}

	if vec.y == 1 {
		ReverseList(traversals.y)
	}

	return traversals
}

func (g *Game) FindFarthestPosition(cell Tile, vector Vector) (*Tile, *Tile) {
	previous := cell
	isFirst := true

	for isFirst || (g.grid.WithinBounds(&cell) && g.grid.CellAvailable(&cell)) {
		previous = cell
		cell = Tile{x: previous.x + vector.x, y: previous.y + vector.y}

		isFirst = false
	}

	return &previous, &cell
}

func (g *Game) positionsEqual(first *Tile, second *Tile) bool {
	return first.x == second.x && first.y == second.y
}

func (g *Game) tileMatchesAvailable() bool {
	for y := 0; y < g.grid.size; y++ {
		for x := 0; x < g.grid.size; x++ {
			t, error := g.grid.CellContent(&Tile{x: x, y: y})

			if error == nil && !t.isEmpty {
				for d := 0; d < 4; d++ {
					vec := g.getVector(d)
					cell := Tile{x: x + vec.x, y: y + vec.y}

					o, error := g.grid.CellContent(&cell)

					if error == nil && !o.isEmpty && o.value == t.value {
						return true
					}
				}
			}
		}
	}

	return false
}

func (g *Game) movesAvailable() bool {
	return g.grid.cellsAvailable() || g.tileMatchesAvailable()
}

func (g *Game) move(direction int) {

	if g.IsGameTerminated() {
		return
	}

	moved := false
	vector := g.getVector(direction)
	traversals := g.BuildTraversals(vector)

	for i := 0; i < len(traversals.x); i++ {
		for j := 0; j < len(traversals.y); j++ {

			cell := Tile{x: traversals.x[i], y: traversals.y[j]}
			tile, error := g.grid.CellContent(&cell)

			if error == nil && !tile.isEmpty {
				farPos, nextPos := g.FindFarthestPosition(cell, vector)
				next, error := g.grid.CellContent(nextPos)
				if error == nil && next.value == tile.value {
					merged := Tile{x: nextPos.x, y: nextPos.y, value: tile.value * 2}
					tiles := make([]Tile, 2)
					tiles[0] = g.copyTile(tile)
					tiles[1] = g.copyTile(nextPos)
					merged.mergedFrom = tiles

					g.grid.insertTile(merged)
					g.grid.removeTile(tile)

					temp := g.copyTile(tile)
					tile = &temp
					tile.updatePosition(nextPos)

					g.score += merged.value
				} else {
					temp := g.moveTile(tile, farPos)
					tile = &temp
				}

				if !g.positionsEqual(&cell, tile) {
					moved = true
				}
			}
		}
	}

	if moved {
		g.addRandomTile()
		if !g.movesAvailable() {
			g.over = true
		}
		g.actuate()
	}

	if g.score > g.highScore {
		g.highScore = g.score
	}
}

func (g *Game) copyTile(tile *Tile) Tile {
	return Tile{x: tile.x, y: tile.y, value: tile.value, isEmpty: tile.isEmpty}
}

func (g *Game) actuate() {
	g.drawer.redraw(g.grid, g.score, g.highScore, g.over)
}
