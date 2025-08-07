package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"tetris/game/grid"
)

type Next struct {
	OX, OY float64
	X, Y   float64
	W, H   float64
	Grid   *grid.Grid
	View   *grid.View
}

func NewNext(grid *grid.Grid, view *grid.View) *Next {
	return &Next{
		Grid: grid,
		View: view,
	}
}

func (n *Next) SetSize(width, height int) {
	n.W, n.H = float64(width), float64(height)
	n.computePos()
}

func (n *Next) SetPosition(x, y int) {
	n.OX, n.OY = float64(x), float64(y)
	n.computePos()
}

func (n *Next) computePos() {
	n.X = n.OX + n.W/2
	n.Y = n.OY + n.H/2
}

func (n *Next) Draw(screen *ebiten.Image) {
	n.View.DrawCenteredTetriminoAt(screen, n.Grid.Next, n.X, n.Y)
}
