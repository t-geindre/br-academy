package game

import (
	"bytes"
	"engine/asset"
	"engine/control"
	"engine/debug"
	"engine/dsp"
	"engine/pool"
	"engine/ui"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"time"
)

type Game struct {
	pool *pool.Pool
}

func NewGame(loader *asset.Loader, settings *Settings) *Game {
	g := &Game{}
	g.init(loader, settings)
	return g
}

func (g *Game) Update() error {
	g.pool.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.pool.Draw(screen)
}

func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}

func (g *Game) init(loader *asset.Loader, settings *Settings) {
	g.pool = pool.NewPool()

	// Filter/pulse configuration
	const SampleRate = 44100

	// Stream
	rawStr, err := mp3.DecodeWithSampleRate(SampleRate, bytes.NewReader(loader.GetRaw("ost")))
	if err != nil {
		panic(err)
	}

	filter := dsp.NewFilterBandPass(SampleRate, settings.Band[0], settings.Band[1])
	stream := dsp.NewStreamPulser(rawStr, filter, settings.Threshold, settings.Release)

	actx := audio.NewContext(SampleRate)
	player, err := actx.NewPlayer(stream)
	if err != nil {
		panic(err)
	}

	// Player
	player.SetBufferSize(time.Millisecond * 20)
	player.SetVolume(1)
	player.Play()

	// pulse visualization
	shd := NewShader(loader.GetShader("pulse"))
	g.pool.Add(shd)

	// Parameters
	multiplier := 1.0
	parameters := NewParameters(
		Param(
			func() string { return fmt.Sprintf("Threshold: %.2f", settings.Threshold) },
			func() { settings.Threshold += 0.01 * multiplier },
			func() { settings.Threshold -= 0.01 * multiplier },
		),
		Param(
			func() string { return fmt.Sprintf("release: %s", settings.Release) },
			func() { settings.Release += time.Millisecond * time.Duration(multiplier) },
			func() { settings.Release -= time.Millisecond * time.Duration(multiplier) },
		),
		Param(
			func() string { return fmt.Sprintf("Band left: %.0f Hz", settings.Band[0]) },
			func() { settings.Band[0] += 1 * float64(multiplier) },
			func() { settings.Band[0] -= 1 * float64(multiplier) },
		),
		Param(
			func() string { return fmt.Sprintf("Band right: %.0f Hz", settings.Band[1]) },
			func() { settings.Band[1] += 1 * float64(multiplier) },
			func() { settings.Band[1] -= 1 * float64(multiplier) },
		),
	)

	// Controls for parameters
	das, arr := time.Millisecond*500, time.Millisecond*60
	repeaters := map[ebiten.Key]*control.Repeater{
		ebiten.KeyArrowUp: control.NewRepeater(das, arr, func() {
			parameters.Add()
		}),
		ebiten.KeyArrowDown: control.NewRepeater(das, arr, func() {
			parameters.Sub()
		}),
	}
	for _, r := range repeaters {
		g.pool.Add(r)
	}

	// Fullscreen toggle
	g.pool.Add(control.NewFsToggle())

	// Debug info
	debugToggle := control.NewToggle(ebiten.KeyF1)
	g.pool.Add(pool.NewDrawer(func(image *ebiten.Image) {
		if !debugToggle.IsOn() {
			return
		}
		debug.DrawFTPS(image)
		ui.DrawPanel(image, ui.TopRight, parameters.Display())
	}), debugToggle)

	// Global updater
	g.pool.Add(pool.NewUpdater(func() {
		// Control multiplier
		multiplier = 1.0
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			multiplier = 10.0
		}

		// Update repeaters
		for k, r := range repeaters {
			if ebiten.IsKeyPressed(k) {
				r.Start()
				continue
			}
			r.Stop()
		}

		// Update parameters
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			parameters.Prev()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			parameters.Next()
		}

		// Apply changes to the stream
		stream.SetRelease(settings.Release)
		stream.SetThreshold(settings.Threshold)

		// Simple loop
		if !player.IsPlaying() {
			player.SetPosition(0)
			player.Play()
		}

		// Forward the gate state to the shader
		if stream.Gate() {
			shd.Pulse()
		}

		// Watch window size
		ww, wh := ebiten.WindowSize()
		settings.WinSize[0] = ww
		settings.WinSize[1] = wh
	}))
}
