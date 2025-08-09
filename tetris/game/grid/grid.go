package grid

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math/rand"
)

type Grid struct {
	W, H   int
	Bricks []*ebiten.DrawImageOptions // Nil is empty

	// Todo there is no need for a next, Bag7[0] should be enough
	Active, Next *FallingTetrimino
	Tetriminos   [7]Tetrimino
	Bag7         []Tetrimino

	Ticks int
	Stats *Stats

	ToClear []int

	Highest int // Highest brick
}

func NewGrid(w, h int) *Grid {
	g := &Grid{
		W:          w,
		H:          h,
		Tetriminos: GetTetriminos(),
		Stats:      NewStats(),
	}

	g.Reset()

	return g
}

// CORE

func (g *Grid) Update() {
	g.FixTetrimino()

	if g.Ticks == 0 {
		g.ClearLines()
		g.ComputeFullLines()
		g.Fall()
		g.Ticks = g.Stats.GetTickRate()
		return
	}

	g.Ticks--
}

func (g *Grid) Reset() {
	g.Active = nil
	g.Next = nil
	g.ToClear = nil

	g.Bag7 = make([]Tetrimino, 0, 7)
	g.Bricks = make([]*ebiten.DrawImageOptions, g.W*g.H)
	g.SpawnTetrimino()

	g.Stats.Reset()
	g.Ticks = g.Stats.GetTickRate()
	g.Highest = 0
}

// TETRIMINO LOGIC

func (g *Grid) GetNextTetrimino() *FallingTetrimino {
	if len(g.Bag7) == 0 {
		g.Bag7 = append(g.Bag7, g.Tetriminos[:]...)
		rand.Shuffle(len(g.Bag7), func(i, j int) {
			g.Bag7[i], g.Bag7[j] = g.Bag7[j], g.Bag7[i]
		})
	}

	t := g.Bag7[0]
	g.Bag7 = g.Bag7[1:]

	return &FallingTetrimino{
		Shape:  t,
		RotIdx: 0,
		X:      g.W/2 - 2,
		Y:      -t.Shapes[0].B, // Start above the grid
	}
}

func (g *Grid) SpawnTetrimino() {
	if g.Next == nil {
		g.Next = g.GetNextTetrimino()
	}

	g.Active = g.Next
	g.Next = g.GetNextTetrimino()
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

	g.computeHighest()
	g.SpawnTetrimino()

	// Force grid ticks, re-sync grid with tetrimino spawn
	g.Ticks = 0
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

func (g *Grid) ComputeFullLines() {
	g.ToClear = nil

	for y := g.H - 1; y >= 0; y-- {
		full := true
		for x := 0; x < g.W; x++ {
			if g.Bricks[y*g.W+x] == nil {
				full = false
				break
			}
		}

		if full {
			g.ToClear = append(g.ToClear, y)
		}
	}
}

func (g *Grid) ClearLines() {
	if len(g.ToClear) == 0 {
		return
	}

	g.Stats.AddLines(len(g.ToClear))

	// Shift lines down
	for y := len(g.ToClear) - 1; y >= 0; y-- {
		for shiftY := g.ToClear[y]; shiftY > 0; shiftY-- {
			for x := 0; x < g.W; x++ {
				g.Bricks[shiftY*g.W+x] = g.Bricks[(shiftY-1)*g.W+x]
			}
		}
	}

	g.ToClear = nil
	g.computeHighest()
}

func (g *Grid) computeHighest() {
	g.Highest = 0
	for y := 0; y < g.H; y++ {
		for x := 0; x < g.W; x++ {
			v := g.H - y
			if g.Bricks[y*g.W+x] != nil && v > g.Highest {
				g.Highest = v
			}
		}
	}
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
