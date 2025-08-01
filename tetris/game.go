package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"os"
	"tetris/grid"
	"ui"
)

const (
	StateInit = iota
	StateLoading
	StateLoaded
	StateRunning
)

type Game struct {
	State         int
	Width, Height int
	Background    *ebiten.Image
	Grid          *grid.Grid
	GridView      *grid.View
}

func NewGame() *Game {
	return &Game{
		State: StateInit,
		Grid:  grid.NewGrid(10, 20),
	}
}

func (g *Game) Update() error {
	// Load assets
	if g.State == StateInit {
		g.State = StateLoading
		go g.Init()
		return nil
	}

	// Setup
	if g.State == StateLoaded {
		g.State = StateRunning
		return nil
	}

	// Running
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.Grid.Rotate()
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.Grid.MoveDown()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.Grid.MoveLeft()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.Grid.MoveRight()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.Grid.Reset()
	}
	g.Grid.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Init
	if g.State == StateInit || g.State == StateLoading {
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

	g.State = StateLoaded
}
