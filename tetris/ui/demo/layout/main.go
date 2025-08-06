package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"math/rand"
	"tetris/ui"
	ui2 "ui"
)

type Game struct {
	layout  *ui.Layout
	colors  map[*ui.Node]color.Color
	drawAll bool
}

func NewGame() *Game {

	return &Game{
		layout:  ui.NewLayout(GetLayout()),
		colors:  make(map[*ui.Node]color.Color),
		drawAll: true,
	}
}

func (g *Game) Update() error {
	ww, wh := ebiten.WindowSize()
	g.layout.Root.W = ww
	g.layout.Root.H = wh
	g.layout.Root.Size = wh
	g.layout.Apply()

	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.drawAll = !g.drawAll
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawNode(g.layout.Root, screen)
	if g.drawAll {
		ui2.DrawFTPS(screen)
	}
}

func (g *Game) DrawNode(node *ui.Node, screen *ebiten.Image) {
	if node.Component != nil || g.drawAll {
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
			"Grow: %.2f\nShrink: %.2f\nSize: %d\nW: %d\nH: %d",
			node.Grow, node.Shrink, node.Size, node.W, node.H,
		)
		if node.Component != nil {
			str = fmt.Sprintf("%s\n%s", node.Component.(Component).Name, str)
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
	ebiten.SetWindowSize(720, 820)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	_ = ebiten.RunGame(NewGame())
}
