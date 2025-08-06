package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"tetris/ui"
)

type Game struct {
	layout *ui.Layout
	colors map[*ui.Node]color.Color
}

func NewGame() *Game {
	return &Game{
		layout: ui.NewLayout(ui.NewNode(ui.NewBox())),
		colors: make(map[*ui.Node]color.Color),
	}
}

func (g *Game) Update() error {
	ww, wh := ebiten.WindowSize()
	g.layout.Root.W = ww
	g.layout.Root.H = wh

	g.layout.Apply()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.DrawNode(g.layout.Root, screen)
}

func (g *Game) DrawNode(node *ui.Node, screen *ebiten.Image) {
	if g.colors[node] == nil {
		g.colors[node] = color.RGBA{
			R: uint8(len(g.colors) * 50 % 255),
			G: uint8(len(g.colors) * 100 % 255),
			B: uint8(len(g.colors) * 150 % 255),
			A: 255,
		}
	}

	vector.DrawFilledRect(
		screen,
		float32(node.X), float32(node.Y),
		float32(node.W), float32(node.H),
		g.colors[node], false,
	)

	for _, child := range node.Children {
		g.DrawNode(child, screen)
	}
}

func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	_ = ebiten.RunGame(NewGame())
}
