package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"tetris/assets"
	"tetris/grid"
	"tetris/ui"
	debug "ui"
)

const (
	StateInit = iota
	StateRunning
)

type Game struct {
	State         int
	Width, Height int
	Grid          *grid.Grid
	GridView      *grid.View
	Controls      *Controls
	Background    *ui.Background
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
	g.GridView.Update()
	g.Background.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.State == StateInit {
		debug.PanelPrintf(screen, debug.BottomRight, "Loading...")
		return
	}

	g.Background.Draw(screen)
	g.GridView.Draw(screen)
	g.GridView.DrawCenteredTetriminoAt(screen, g.Grid.Next, 500.0, 210.0)

	debug.DrawFTPS(screen)
}

func (g *Game) Layout(x, y int) (int, int) {
	return x, y
}

func (g *Game) Init() {
	go func() {
		g.Loader.MustLoad()

		g.GridView = grid.NewView(
			g.Grid, 48, 4,
			g.Loader.GetImage("brick"),
			g.Loader.GetShader("disappear"),
			g.Loader.GetShader("grid"),
		)

		g.Background = ui.NewBackground(
			g.Loader.GetShader("background"),
		)

		g.Width, g.Height = 720, 820 // todo fixme

		ebiten.SetWindowSize(g.Width, g.Height)

		g.State = StateRunning
	}()
}
