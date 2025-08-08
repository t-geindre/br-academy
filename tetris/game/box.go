package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Box struct {
	Shader     *ebiten.Shader
	Opts       *ebiten.DrawRectShaderOptions
	X, Y, W, H int
	Padding    int
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
	b.W, b.H = width+b.Padding*2, height+b.Padding*2
}

func (b *Box) SetPosition(x, y int) {
	b.X, b.Y = x-b.Padding, y-b.Padding
}
