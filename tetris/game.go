package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"tetris/assets"
	"tetris/grid"
	"ui"
)

const (
	StateInit = iota
	StateRunning
)

type Game struct {
	State         int
	Width, Height int
	Background    *ebiten.Image
	Grid          *grid.Grid
	GridView      *grid.View
	Controls      *Controls
	Loader        *assets.Loader
}

func NewGame(loader *assets.Loader) *Game {
	gr := grid.NewGrid(10, 20)

	g := &Game{
		State:    StateInit,
		Grid:     gr,
		Controls: NewControls(gr),
		Loader:   loader,
	}

	g.Init()

	return g
}

func (g *Game) Update() error {
	// Assets loading
	if g.State == StateInit {
		return nil
	}

	// Running
	g.Controls.Update()
	g.Grid.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.State == StateInit {
		ui.PanelPrintf(screen, ui.BottomRight, "Loading...")
		return
	}

	screen.DrawImage(g.Background, nil)
	g.GridView.Draw(screen)
	g.GridView.DrawCenteredTetriminoAt(screen, g.Grid.Next, 868.0, 210.0)

	ui.DrawFTPS(screen)
}

func (g *Game) Layout(x, y int) (int, int) {
	if g.State == StateRunning {
		return g.Width, g.Height
	}

	return x, y
}

func (g *Game) Init() {
	go func() {
		g.Loader.MustLoad()

		g.Background = g.Loader.GetImage("background")
		g.GridView = grid.NewView(
			g.Grid, 64, 32,
			g.Loader.GetImage("brick"),
			g.Loader.GetShader("disappear"),
		)

		bds := g.Background.Bounds()
		g.Width, g.Height = bds.Dx(), bds.Dy()

		ebiten.SetWindowSize(g.Width, g.Height)

		g.State = StateRunning
	}()
}
