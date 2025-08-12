package game

import (
	"bytes"
	"engine/asset"
	"engine/control"
	"engine/debug"
	"engine/dsp"
	"engine/pool"
	"engine/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/mp3"
	"time"
)

type Game struct {
	pool *pool.Pool
}

func NewGame(loader *asset.Loader) *Game {
	g := &Game{}
	g.init(loader)
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

func (g *Game) init(loader *asset.Loader) {
	g.pool = pool.NewPool()

	// Filter/pulse configuration
	threshold := 0.26
	release := time.Millisecond * 400
	band := [2]float64{50, 100}

	const SampleRate = 44100

	// Stream
	rawStr, err := mp3.DecodeWithSampleRate(SampleRate, bytes.NewReader(loader.GetRaw("ost")))
	if err != nil {
		panic(err)
	}

	filter := dsp.NewFilterBandPass(SampleRate, band[0], band[1])
	stream := dsp.NewStreamPulser(rawStr, filter, threshold, release)

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

	// Debug info
	g.pool.Add(pool.NewDrawer(func(image *ebiten.Image) {
		debug.DrawFTPS(image)
		ui.DrawPanel(image, ui.TopRight, "Threshold: %.2f", threshold)
		ui.DrawPanel(image, ui.BottomRight, "Release: %s", release)
		ui.DrawPanel(image, ui.BottomLeft, "Band: %.0f-%.0f Hz", band[0], band[1])
	}))

	// Controls
	mult, das, arr := 1.0, time.Millisecond*500, time.Millisecond*60
	bindings := map[ebiten.Key]*control.Repeater{
		ebiten.KeyArrowUp:    control.NewRepeater(das, arr, func() { threshold += 0.01 * mult }),
		ebiten.KeyArrowDown:  control.NewRepeater(das, arr, func() { threshold -= 0.01 * mult }),
		ebiten.KeyArrowRight: control.NewRepeater(das, arr, func() { release += time.Millisecond * time.Duration(mult) }),
		ebiten.KeyArrowLeft:  control.NewRepeater(das, arr, func() { release -= time.Millisecond * time.Duration(mult) }),
	}
	for _, r := range bindings {
		g.pool.Add(r)
	}

	// Global updater
	g.pool.Add(pool.NewUpdater(func() {
		// Control multiplier
		mult = 1.0
		if ebiten.IsKeyPressed(ebiten.KeyShift) {
			mult = 10.0
		}

		// Update repeaters
		for k, r := range bindings {
			if ebiten.IsKeyPressed(k) {
				r.Start()
				continue
			}
			r.Stop()
		}

		// Player seek
		p := time.Duration(0)
		if ebiten.IsKeyPressed(ebiten.KeyPageUp) {
			p += time.Millisecond * 100 * time.Duration(mult)
		}
		if ebiten.IsKeyPressed(ebiten.KeyPageDown) {
			p -= time.Millisecond * 100 * time.Duration(mult)
		}
		if p != 0 {
			player.SetPosition(p + player.Position())
		}

		// Apply changes to the stream
		stream.SetRelease(release)
		stream.SetThreshold(threshold)

		// Simple loop
		if !player.IsPlaying() {
			player.SetPosition(0)
			player.Play()
		}

		// Forward the gate state to the shader
		if stream.Gate() {
			shd.Pulse()
		}
	}))
}
