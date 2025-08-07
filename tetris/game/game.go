package game

import (
	"component"
	"debug"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"tetris/assets"
	"tetris/game/grid"
	"ui"
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
	background    *Background
	loader        *assets.Loader
	layout        *Layout

	TitleNext *component.Text
	ValueNext *Next

	TitleScore *component.Text
	ValueScore *component.UpdatableText

	TitleLevel *component.Text
	ValueLevel *component.UpdatableText

	TitleLines *component.Text
	ValueLines *component.UpdatableText
}

func NewGame(loader *assets.Loader) *Game {
	gr := grid.NewGrid(10, 20)

	g := &Game{
		state:    StateInit,
		grid:     gr,
		controls: NewControls(gr),
		loader:   loader,
		layout:   NewLayout(600, 760),
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
	g.ValueScore.Update()
	g.ValueLevel.Update()
	g.ValueLines.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.state == StateInit {
		ui.DrawPanel(screen, ui.BottomRight, "Loading...")
		return
	}

	g.background.Draw(screen)
	g.gridView.Draw(screen)

	g.TitleNext.Draw(screen)
	g.ValueNext.Draw(screen)

	g.TitleScore.Draw(screen)
	g.ValueScore.Draw(screen)

	g.TitleLevel.Draw(screen)
	g.ValueLevel.Draw(screen)

	g.TitleLines.Draw(screen)
	g.ValueLines.Draw(screen)

	debug.DrawLayoutNode(screen, g.layout.Root)
	debug.All(screen)
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

		g.background = NewBackground(
			g.loader.GetShader("background"),
		)
		g.layout.Container.Component = g.background

		titleFont := &text.GoTextFace{
			Source: g.loader.GetFont("bold"),
			Size:   40,
		}

		g.TitleNext = component.NewText("NEXT", 500, 60, titleFont)
		g.layout.NextTitle.Component = g.TitleNext

		g.ValueNext = NewNext(g.grid, g.gridView)
		g.layout.NextValue.Component = g.ValueNext

		g.TitleScore = component.NewText("SCORE", 500, 120, titleFont)
		g.layout.ScoreTitle.Component = g.TitleScore

		g.ValueScore = component.NewUpdatableText(
			func() string {
				return fmt.Sprintf("%d", g.grid.Stats.Score)
			}, 500, 150, titleFont)
		g.layout.ScoreValue.Component = g.ValueScore

		g.TitleLevel = component.NewText("LEVEL", 500, 180, titleFont)
		g.layout.LevelTitle.Component = g.TitleLevel

		g.ValueLevel = component.NewUpdatableText(
			func() string {
				return fmt.Sprintf("%d", g.grid.Stats.Level)
			}, 500, 210, titleFont)
		g.layout.LevelValue.Component = g.ValueLevel

		g.TitleLines = component.NewText("LINES", 500, 240, titleFont)
		g.layout.LinesTitle.Component = g.TitleLines

		g.ValueLines = component.NewUpdatableText(
			func() string {
				return fmt.Sprintf("%d", g.grid.Stats.Lines)
			},
			500, 270, titleFont)
		g.layout.LinesValue.Component = g.ValueLines

		g.width, g.height = 1024, 768 // todo fixme

		ebiten.SetWindowSize(g.width, g.height)

		g.state = StateRunning
	}()
}
