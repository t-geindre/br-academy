package game

import (
	"bytes"
	"engine/asset"
	"engine/component"
	"engine/control"
	"engine/debug"
	"engine/dsp"
	"engine/pool"
	"engine/ui"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"image/color"
	"tetris/game/grid"
	"time"
)

const (
	StateInit = iota
	StateRunning
)

type Game struct {
	state int
	pool  *pool.Pool
}

func NewGame(loader *asset.Loader) *Game {
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
	layout := NewLayout()
	g.pool.Add(layout)

	// Background
	bg := NewBackground(loader.GetShader("background"))
	g.pool.Add(bg)

	// Background particles
	particles := NewParticles(loader.GetShader("particles"))
	g.pool.Add(particles)

	// Main container
	container := NewBox(loader.GetShader("box"), color.White, 2, 10, 0.5)
	g.pool.Add(container)
	layout.Container.Attach(container)

	// Grid
	gr := grid.NewGrid(10, 20)
	g.pool.Add(gr)

	// Controls
	g.pool.Add(NewControls(gr))

	// Grid box
	grBox := NewBox(loader.GetShader("box"), color.White, 2, 10, 0.6)
	grBox.Padding = 5
	g.pool.Add(grBox)
	layout.Grid.Attach(grBox)

	// Grid view
	grView := grid.NewView(
		gr, 4, 32,
		loader.GetImage("brick"),
		loader.GetShader("disappear"),
		loader.GetShader("grid"),
	)
	g.pool.Add(grView)
	layout.Grid.Attach(grView)

	// Prepare font faces
	titleFont := &text.GoTextFace{Source: loader.GetFont("bold"), Size: 40}
	normalFont := &text.GoTextFace{Source: loader.GetFont("normal"), Size: 40}

	// Next piece
	nextTitle := component.NewText("NEXT", 0, 0, titleFont)
	g.pool.Add(nextTitle)
	layout.NextTitle.Attach(nextTitle)

	nextValue := NewNext(gr, grView)
	g.pool.Add(nextValue)
	layout.NextValue.Attach(nextValue)

	// Score
	scoreTitle := component.NewText("SCORE", 0, 0, titleFont)
	g.pool.Add(scoreTitle)
	layout.ScoreTitle.Attach(scoreTitle)

	scoreValue := component.NewUpdatableText(func() string {
		return fmt.Sprintf("%d", gr.Stats.Score)
	}, 0, 0, normalFont)
	g.pool.Add(scoreValue)
	layout.ScoreValue.Attach(scoreValue)

	// Level
	levelTitle := component.NewText("LEVEL", 0, 0, titleFont)
	g.pool.Add(levelTitle)
	layout.LevelTitle.Attach(levelTitle)

	levelValue := component.NewUpdatableText(func() string {
		return fmt.Sprintf("%d", gr.Stats.Level)
	}, 0, 0, normalFont)
	g.pool.Add(levelValue)
	layout.LevelValue.Attach(levelValue)

	// Lines
	linesTitle := component.NewText("LINES", 0, 0, titleFont)
	g.pool.Add(linesTitle)
	layout.LinesTitle.Attach(linesTitle)

	linesValue := component.NewUpdatableText(func() string {
		return fmt.Sprintf("%d", gr.Stats.Lines)
	}, 0, 0, normalFont)
	g.pool.Add(linesValue)
	layout.LinesValue.Attach(linesValue)

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

	// Audio
	const SampleRate = 44100

	rawStr, err := mp3.DecodeWithSampleRate(44100, bytes.NewReader(loader.GetRaw("audio-st")))
	if err != nil {
		panic(err)
	}

	// Main theme time.Millisecond*41142, time.Millisecond*109714
	// Old theme 0, time.Millisecond*13714
	loopStream := dsp.NewLooper(44100, rawStr)
	loopStream.Play(loopStream.AddLoop(time.Millisecond*41142, time.Millisecond*109714))

	pulserStr := dsp.NewStreamPulser(
		loopStream,
		dsp.NewFilterBandPass(SampleRate, 50, 100),
		0.20,                 // threshold
		time.Millisecond*300, // release
	)

	ctx := audio.NewContext(44100)
	player, err := ctx.NewPlayer(pulserStr)
	if err != nil {
		panic(err)
	}

	player.SetVolume(1)
	player.SetBufferSize(time.Millisecond * 20) // Small buffer for low latency sync
	player.SetPosition(time.Second * 35)        // todo remove me
	player.Play()

	// Updates

	// Grid danger effect
	g.pool.Add(pool.NewUpdater(func() {
		// Grid danger effect
		v := float32(gr.Highest)
		if v > 15 {
			bg.SetDanger((v - 15) / 5)
		} else {
			bg.SetDanger(0)
		}

		// Pulser
		if pulserStr.Gate() {
			particles.Pulse()
			bg.Pulse()
		}
	}))
}
