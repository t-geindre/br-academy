package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Box struct {
	Shader     *ebiten.Shader
	Opts       *ebiten.DrawRectShaderOptions
	X, Y, W, H int
}

func NewBox(shader *ebiten.Shader, col color.Color, thick, rad, dark float32) *Box {
	r, g, b, a := col.RGBA()
	colF := []float32{
		float32(r) / 0xffff, float32(g) / 0xffff, float32(b) / 0xffff, float32(a) / 0xffff,
	}

	return &Box{
		Shader: shader,
		Opts: &ebiten.DrawRectShaderOptions{
			Uniforms: map[string]any{
				"BoxSize":         [2]float32{1, 1},
				"BoxDarken":       [4]float32{0, 0, 0, dark},
				"CornerRadius":    rad,
				"BorderThickness": thick,
				"BorderColor":     colF,
			},
		},
	}
}

func (b *Box) Draw(screen *ebiten.Image) {
	b.Opts.GeoM.Reset()
	b.Opts.GeoM.Translate(float64(b.X), float64(b.Y))
	screen.DrawRectShader(b.W, b.H, b.Shader, b.Opts)
}

func (b *Box) SetSize(width, height int) {
	b.Opts.Uniforms["BoxSize"] = [2]float32{float32(width), float32(height)}
	b.W, b.H = width, height
}

func (b *Box) SetPosition(x, y int) {
	b.X, b.Y = x, y
}
