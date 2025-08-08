package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Background struct {
	Shader *ebiten.Shader
	Opts   *ebiten.DrawRectShaderOptions
	Time   float32
}

func NewBackground(shader *ebiten.Shader) *Background {
	b := &Background{}
	b.setShader(shader)
	return b
}

func (b *Background) Draw(screen *ebiten.Image) {
	screen.DrawRectShader(screen.Bounds().Dx(), screen.Bounds().Dy(), b.Shader, b.Opts)
}

func (b *Background) Update() {
	b.Time++
	b.Opts.Uniforms["Time"] = b.Time
}

func (b *Background) SetSize(width, height int) {
	b.Opts.Uniforms["BoardSize"] = [2]float32{float32(width), float32(height)}
}

func (b *Background) SetPosition(x, y int) {
	b.Opts.Uniforms["BoardPos"] = [2]float32{float32(x), float32(y)}
}

func (b *Background) setShader(shader *ebiten.Shader) {
	b.Shader = shader
	b.Opts = &ebiten.DrawRectShaderOptions{
		Uniforms: map[string]any{
			"BaseColor":  [4]float32{0.047, 0.051, 0.271, 1.0},  // #0C0D45
			"GlowLeft":   [4]float32{00.098, 0.623, 0.863, 1.0}, // #199FDC
			"GlowRight":  [4]float32{0.510, 0.176, 0.592, 1.0},  // #822D97
			"Intensity":  float32(.9),
			"Spread":     float32(.28),
			"VertSpread": float32(.45),
			"Time":       float32(0),
		},
	}
}
