package main

import (
	"engine/asset"
	"engine/debug"
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	_ = ebiten.RunGame(NewGame())
}

type Game struct {
	shd   *ebiten.Shader
	opt   *ebiten.DrawRectShaderOptions
	start time.Time
}

func NewGame() *Game {
	return &Game{
		shd: asset.MustLoadShader("fractal.kage"),
		opt: &ebiten.DrawRectShaderOptions{
			Uniforms: map[string]any{
				"Time": float32(0),
			},
		},
		start: time.Now(),
	}
}

func (g *Game) Update() error {
	g.opt.Uniforms["Time"] = float32(time.Since(g.start).Seconds())
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	bds := screen.Bounds()
	screen.DrawRectShader(bds.Dx(), bds.Dy(), g.shd, g.opt)
	debug.DrawFTPS(screen)
}

func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}
