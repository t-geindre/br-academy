package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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
	state         int
	width, height int
	grid          *grid.Grid
	gridView      *grid.View
	controls      *Controls
	background    *ui.Background
	loader        *assets.Loader
	layout        *Layout
	TitleNext     *ui.Text
}

func NewGame(loader *assets.Loader) *Game {
	gr := grid.NewGrid(10, 20)

	g := &Game{
		state:    StateInit,
		grid:     gr,
		controls: NewControls(gr),
		loader:   loader,
		layout:   NewLayout(1024, 768),
	}

	g.Init()

	return g
}

func (g *Game) Update() error {
	// Assets loading
	if g.state == StateInit {
		return nil
	}

	// Running
	g.controls.Update()
	g.grid.Update()
	g.gridView.Update()
	g.background.Update()
	g.layout.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.state == StateInit {
		debug.PanelPrintf(screen, debug.BottomRight, "Loading...")
		return
	}

	g.background.Draw(screen)
	g.gridView.Draw(screen)
	g.gridView.DrawCenteredTetriminoAt(screen, g.grid.Next, 525.0, 130.0)
	g.TitleNext.Draw(screen)

	debug.DrawFTPS(screen)
}

func (g *Game) Layout(x, y int) (int, int) {
	return x, y
}

func (g *Game) Init() {
	go func() {
		g.loader.MustLoad()

		g.gridView = grid.NewView(
			g.grid, 4, 32,
			g.loader.GetImage("brick"),
			g.loader.GetShader("disappear"),
			g.loader.GetShader("grid"),
		)
		g.layout.Grid.Component = g.gridView

		g.background = ui.NewBackground(
			g.loader.GetShader("background"),
		)
		g.layout.Container.Component = g.background

		g.TitleNext = ui.NewText(
			"NEXT",
			500, 60,
			&text.GoTextFace{
				Source: g.loader.GetFont("bold"),
				Size:   40,
			},
		)
		g.layout.NextTitle.Component = g.TitleNext

		g.width, g.height = 1024, 768 // todo fixme

		ebiten.SetWindowSize(g.width, g.height)

		g.state = StateRunning
	}()
}
