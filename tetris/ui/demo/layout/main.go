package main

import (
	"debug"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"layout"
	"math/rand"
	"tetris/game"
)

type Game struct {
	layout  *game.Layout
	colors  map[*layout.Node]color.Color
	drawAll bool
}

func NewGame() *Game {

	return &Game{
		layout: game.NewLayout(),
		colors: make(map[*layout.Node]color.Color),
	}
}

func (g *Game) Update() error {
	g.layout.Update()

	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.drawAll = !g.drawAll
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawNode(g.layout.Root, screen)
	if g.drawAll {
		debug.All(screen)
	}
}

func (g *Game) DrawNode(node *layout.Node, screen *ebiten.Image) {
	if len(node.Name) > 0 || g.drawAll {
		if g.colors[node] == nil {
			g.colors[node] = color.RGBA{
				R: uint8(rand.Intn(100) + 55),
				G: uint8(rand.Intn(100) + 55),
				B: uint8(rand.Intn(100) + 55),
				A: 255,
			}
		}

		vector.DrawFilledRect(
			screen,
			float32(node.X), float32(node.Y),
			float32(node.W), float32(node.H),
			g.colors[node], false,
		)

		str := fmt.Sprintf(
			"Grow: %.2f, Shrink: %.2f\nSize: %d (%d,%d)\nPos: %d,%d",
			node.Grow, node.Shrink, node.Size, node.W, node.H, node.X, node.Y,
		)
		if len(node.Name) > 0 {
			str = fmt.Sprintf("%s\n%s", node.Name, str)
		}

		ebitenutil.DebugPrintAt(screen, str, node.X+5, node.Y+5)
	}

	for _, child := range node.Children {
		g.DrawNode(child, screen)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}

func main() {
	ebiten.SetWindowSize(1024, 768)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	_ = ebiten.RunGame(NewGame())
}
