package grid

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
)

type Grid struct {
	W, H   int
	Bricks []*ebiten.DrawImageOptions // Nil is empty

	Active, Next *FallingTetrimino
	Tetriminos   [7]Tetrimino
	Bag7         []Tetrimino

	Level int
	Score int
	Lines int

	Ticks int
}

func NewGrid(w, h int) *Grid {
	g := &Grid{
		W:          w,
		H:          h,
		Tetriminos: GetTetriminos(),
	}

	g.Reset()

	return g
}

// CORE

func (g *Grid) Update() {
	g.FixTetrimino()

	if g.Ticks == 0 {
		g.Fall()
		g.ClearLines()
		if g.Level == 0 {
			g.Ticks = int(48.0 / 60.0 * float64(ebiten.TPS()))
			return
		}
		if g.Level == 1 {
			g.Ticks = int(43.0 / 60.0 * float64(ebiten.TPS()))
			return
		}
		if g.Level < 10 {
			g.Ticks = int(28.0 / 60.0 * float64(ebiten.TPS()))
			return
		}
		if g.Level < 20 {
			g.Ticks = int(18.0 / 60.0 * float64(ebiten.TPS()))
			return
		}
		if g.Level < 30 {
			g.Ticks = int(8.0 / 60.0 * float64(ebiten.TPS()))
			return
		}
		g.Ticks = int(1.0 / 60.0 * float64(ebiten.TPS()))
	}

	g.Ticks--
	/*
		Cool effect when losing a game
		for i := range g.Bricks {
			g.Bricks[i] = g.Tetriminos[rand.Intn(len(g.Tetriminos))].Color
		}
	*/
}

// CONTROLS

func (g *Grid) MoveLeft() {
	if !g.CollidesAt(g.Active.X-1, g.Active.Y, g.Active.RotIdx) && g.ResetFixDelay() {
		g.Active.X--
	}
}

func (g *Grid) MoveRight() {
	if !g.CollidesAt(g.Active.X+1, g.Active.Y, g.Active.RotIdx) && g.ResetFixDelay() {
		g.Active.X++

	}
}

func (g *Grid) MoveDown() {
	g.Fall()
}

func (g *Grid) Rotate() {
	newRot := (g.Active.RotIdx + 1) % len(g.Active.Shape.Shapes)

	// Wall kicks
	kicks := [][2]int{
		{0, 0},
		{-1, 0},
		{1, 0},
		{0, -1},
		{-2, 0},
		{2, 0},
	}

	for _, k := range kicks {
		dx, dy := k[0], k[1]
		newX := g.Active.X + dx
		newY := g.Active.Y + dy

		if !g.CollidesAt(newX, newY, newRot) && g.ResetFixDelay() {
			g.Active.X = newX
			g.Active.Y = newY
			g.Active.RotIdx = newRot
			return
		}
	}

	// Rotation canceled
}

// TETRIMINO LOGIC

func (g *Grid) GetNextTetrimino() Tetrimino {
	if len(g.Bag7) == 0 {
		g.Bag7 = append(g.Bag7, g.Tetriminos[:]...)
		rand.Shuffle(len(g.Bag7), func(i, j int) {
			g.Bag7[i], g.Bag7[j] = g.Bag7[j], g.Bag7[i]
		})
	}

	t := g.Bag7[0]
	g.Bag7 = g.Bag7[1:]

	return t
}

func (g *Grid) SpawnTetrimino() {
	if g.Next == nil {
		s := g.GetNextTetrimino()
		g.Next = &FallingTetrimino{
			Shape:  s,
			RotIdx: 0,
			X:      g.W/2 - 2,
			Y:      -s.Shapes[0].B, // Start above the grid
		}
	}

	g.Active = g.Next
	g.Next = nil
}

func (g *Grid) FixTetrimino() {
	if !g.ShouldFix() {
		return
	}

	if !g.Active.IsFixing {
		g.Active.IsFixing = true
		g.Active.FReset = 16 // Will be 15
		g.ResetFixDelay()
		return
	}

	g.Active.Fixing--
	if g.Active.Fixing > 0 {
		return
	}

	rot := g.Active.Shape.Shapes[g.Active.RotIdx]
	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if rot.S[y*4+x] != 0 {
				gridX := g.Active.X + x
				gridY := g.Active.Y + y

				if gridX >= 0 && gridX < g.W && gridY >= 0 && gridY < g.H {
					g.Bricks[gridY*g.W+gridX] = g.Active.Shape.Color
				}
			}
		}
	}

	g.SpawnTetrimino()
}

func (g *Grid) ResetFixDelay() bool {
	if !g.Active.IsFixing {
		return true
	}

	g.Active.FReset--
	if g.Active.FReset <= 0 {
		return false
	}

	g.Active.Fixing = ebiten.TPS() / 2
	return true
}

func (g *Grid) ShouldFix() bool {
	rot := g.Active.Shape.Shapes[g.Active.RotIdx]

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			if rot.S[y*4+x] == 0 {
				continue
			}

			gridX := g.Active.X + x
			gridY := g.Active.Y + y + 1

			if gridY >= g.H {
				return true
			}

			if gridX >= 0 && gridX < g.W && gridY >= 0 {
				if g.Bricks[gridY*g.W+gridX] != nil {
					return true
				}
			}
		}
	}

	return false
}

func (g *Grid) Fall() {
	if !g.CollidesAt(g.Active.X, g.Active.Y+1, g.Active.RotIdx) {
		g.Active.Y++
	}
}

func (g *Grid) CollidesAt(x, y, rotIdx int) bool {
	shape := g.Active.Shape.Shapes[rotIdx]

	for dy := 0; dy < 4; dy++ {
		for dx := 0; dx < 4; dx++ {
			if shape.S[dy*4+dx] == 0 {
				continue
			}

			gridX := x + dx
			gridY := y + dy

			if gridX < 0 || gridX >= g.W || gridY >= g.H {
				return true
			}
			if gridY >= 0 && g.Bricks[gridY*g.W+gridX] != nil {
				return true
			}
		}
	}

	return false
}

func (g *Grid) ClearLines() {
	linesCleared := 0

	for y := g.H - 1; y >= 0; y-- {
		full := true
		for x := 0; x < g.W; x++ {
			if g.Bricks[y*g.W+x] == nil {
				full = false
				break
			}
		}

		if full {
			linesCleared++

			// Move down
			for ty := y; ty > 0; ty-- {
				for x := 0; x < g.W; x++ {
					g.Bricks[ty*g.W+x] = g.Bricks[(ty-1)*g.W+x]
				}
			}

			// Clear the top line
			for x := 0; x < g.W; x++ {
				g.Bricks[x] = nil
			}

			// Stay on the same line
			y++
		}
	}

	if linesCleared == 0 {
		return
	}

	g.Score += []int{0, 100, 300, 500, 800}[linesCleared] * (g.Level + 1)
	g.Lines += linesCleared
	g.Level = g.Lines / 10
}

func (g *Grid) Reset() {
	g.Bricks = make([]*ebiten.DrawImageOptions, g.W*g.H)
	g.Active = nil
	g.Next = nil
	g.Bag7 = make([]Tetrimino, 0, 7)
	g.SpawnTetrimino()
	g.Level = 0
	g.Score = 0
	g.Lines = 0
	g.Ticks = 0
}
