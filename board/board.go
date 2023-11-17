package board

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math/rand"
)

const pGenerate = 0.5

type Board struct {
	Grid  Grid
	Rules Ruler
}

type Grid [][]bool

type Ruler interface {
	UnderPopulation(grid Grid, i, j int) bool
	OverPopulation(grid Grid, i, j int) bool
	Reproduce(grid Grid, i, j int) bool
}

const (
	Rows = 50
	Cols = 100
	Size = 10
)

func New(r Ruler) *Board {
	return &Board{
		Grid:  GenerateCell(),
		Rules: r,
	}
}

func GenerateCell() Grid {
	grid := make(Grid, Rows)
	for i := range grid {
		grid[i] = make([]bool, Cols)
		for j := range grid[i] {
			grid[i][j] = rand.Intn(100) <= pGenerate*100
		}
	}
	return grid
}

func (b *Board) Draw(win *pixelgl.Window) {
	win.Clear(pixel.RGB(1, 1, 1)) // Blanc

	for i, row := range b.Grid {
		for j, alive := range row {
			if alive {
				b.DrawCell(win, j*Size, (Rows-i-1)*Size)
			}
		}
	}

	win.Update()
}

func (b *Board) DrawCell(win *pixelgl.Window, x, y int) {
	rect := pixel.R(float64(x), float64(y), float64(x+Size), float64(y+Size))
	sprite := pixel.NewSprite(nil, rect)
	sprite.Draw(win, pixel.IM.Moved(rect.Center()).Scaled(rect.Center(), 0.5))
}

func (b *Board) Update() {
	newGrid := make(Grid, Rows)
	for i := range b.Grid {
		newGrid[i] = make([]bool, Cols)
		copy(newGrid[i], b.Grid[i])
	}

	for i, row := range b.Grid {
		for j := range row {
			if isAlive(b.Grid, i, j) {
				if b.Rules.UnderPopulation(b.Grid, i, j) {
					newGrid[i][j] = false
				}
				if b.Rules.OverPopulation(b.Grid, i, j) {
					newGrid[i][j] = false
				}
			} else {
				if b.Rules.Reproduce(b.Grid, i, j) {
					newGrid[i][j] = true
				}
			}
		}
	}

	b.Grid = newGrid
}
