package game

import (
	"component"
	"control"
	"debug"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"pool"
	"tetris/assets"
	"tetris/game/grid"
	"ui"
)

const (
	StateInit = iota
	StateRunning
)

type Game struct {
	state int
	pool  *pool.Pool
}

func NewGame(loader *assets.Loader) *Game {
	g := &Game{
		state: StateInit,
		pool:  pool.NewPool(),
	}

	go g.Init()

	return g
}

func (g *Game) Update() error {
	if g.state == StateRunning {
		g.pool.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.state != StateRunning {
		ui.DrawPanel(screen, ui.TopLeft, "Loading...")
		return
	}

	g.pool.Draw(screen)
}

func (g *Game) Layout(x, y int) (int, int) {
	return x, y
}

func (g *Game) Init() {
	// Load assets
	loader := GetAssetsLoader()
	loader.MustLoad()

	// Build layout
	layout := NewLayout(600, 760)
	g.pool.Add(layout)

	// Background
	bg := NewBackground(loader.GetShader("background"))
	g.pool.Add(bg)
	layout.Container.Component = bg

	// Grid
	gr := grid.NewGrid(10, 20)
	g.pool.Add(gr)

	// Controls
	g.pool.Add(NewControls(gr))

	// Grid view
	grView := grid.NewView(
		gr, 4, 32,
		loader.GetImage("brick"),
		loader.GetShader("disappear"),
		loader.GetShader("grid"),
	)
	g.pool.Add(grView)
	layout.Grid.Component = grView

	// Prepare font faces
	titleFont := &text.GoTextFace{Source: loader.GetFont("bold"), Size: 40}
	normalFont := &text.GoTextFace{Source: loader.GetFont("normal"), Size: 40}

	// Next piece
	nextTitle := component.NewText("NEXT", 0, 0, titleFont)
	g.pool.Add(nextTitle)
	layout.NextTitle.Component = nextTitle

	nextValue := NewNext(gr, grView)
	g.pool.Add(nextValue)
	layout.NextValue.Component = nextValue

	// Score
	scoreTitle := component.NewText("SCORE", 0, 0, titleFont)
	g.pool.Add(scoreTitle)
	layout.ScoreTitle.Component = scoreTitle

	scoreValue := component.NewUpdatableText(func() string {
		return fmt.Sprintf("%d", gr.Stats.Score)
	}, 0, 0, normalFont)
	g.pool.Add(scoreValue)
	layout.ScoreValue.Component = scoreValue

	// Level
	levelTitle := component.NewText("LEVEL", 0, 0, titleFont)
	g.pool.Add(levelTitle)
	layout.LevelTitle.Component = levelTitle

	levelValue := component.NewUpdatableText(func() string {
		return fmt.Sprintf("%d", gr.Stats.Level)
	}, 0, 0, normalFont)
	g.pool.Add(levelValue)
	layout.LevelValue.Component = levelValue

	// Lines
	linesTitle := component.NewText("LINES", 0, 0, titleFont)
	g.pool.Add(linesTitle)
	layout.LinesTitle.Component = linesTitle

	linesValue := component.NewUpdatableText(func() string {
		return fmt.Sprintf("%d", gr.Stats.Lines)
	}, 0, 0, normalFont)
	g.pool.Add(linesValue)
	layout.LinesValue.Component = linesValue

	// Debug panels
	dbgOverlayCtrl := control.NewToggle(ebiten.KeyF2)
	g.pool.Add(dbgOverlayCtrl)
	g.pool.Add(pool.NewDrawer(func(image *ebiten.Image) {
		if dbgOverlayCtrl.IsOn() {
			debug.DrawAll(image)
		}
	}))

	// Debug layout
	dbgLayoutCtrl := control.NewToggle(ebiten.KeyF3)
	g.pool.Add(dbgLayoutCtrl)
	g.pool.Add(pool.NewDrawer(func(image *ebiten.Image) {
		if dbgLayoutCtrl.IsOn() {
			debug.DrawLayoutNode(image, layout.Root)
		}
	}))

	// Run the game
	g.state = StateRunning
}
