package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"os"
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
}

func NewGame() *Game {
	gr := grid.NewGrid(10, 20)

	g := &Game{
		State:    StateInit,
		Grid:     gr,
		Controls: NewControls(gr),
	}

	go g.Init()

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
	// Init
	if g.State == StateInit {
		ui.PanelPrintf(screen, ui.BottomRight, "Loading...")
		return
	}

	// Running
	screen.DrawImage(g.Background, nil)
	g.GridView.Draw(screen)
	g.GridView.DrawCenteredTetriminoAt(screen, g.Grid.Next, 868.0, 210.0)

	ui.PanelPrintf(
		screen, ui.BottomRight,
		"[F1] Reset\nLevel: %d\nScore: %d\nLines: %d",
		g.Grid.Stats.Level, g.Grid.Stats.Score, g.Grid.Stats.Lines,
	)

	ui.DrawFTPS(screen)
}

func (g *Game) Layout(x, y int) (int, int) {
	if g.State == StateRunning {
		return g.Width, g.Height
	}

	return x, y
}

func (g *Game) Init() {
	// Load UI
	bg, _, err := ebitenutil.NewImageFromFile("assets/ui.png")
	if err != nil {
		panic(err)
	}
	g.Background = bg

	// Load brick
	brick, _, err := ebitenutil.NewImageFromFile("assets/brick.png")
	if err != nil {
		panic(err)
	}

	// Load disappearing shader
	rawShd, err := os.ReadFile("assets/shaders/disappear.kage")
	if err != nil {
		panic(err)
	}

	// Compile shader
	dShd, err := ebiten.NewShader(rawShd)
	if err != nil {
		panic(err)
	}

	// Loading done
	g.GridView = grid.NewView(g.Grid, 64, 32, brick, dShd)

	bds := g.Background.Bounds()
	g.Width, g.Height = bds.Dx(), bds.Dy()

	ebiten.SetWindowSize(g.Width, g.Height)

	g.State = StateRunning
}
