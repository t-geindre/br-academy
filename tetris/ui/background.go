package ui

import "github.com/hajimehoshi/ebiten/v2"

type Background struct {
	Shader *ebiten.Shader
	Time   float32
}

func NewBackground(shader *ebiten.Shader) *Background {
	return &Background{
		Shader: shader,
	}
}

func (b *Background) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawRectShaderOptions{
		Uniforms: map[string]interface{}{
			"BaseColor": [4]float32{0.047, 0.051, 0.271, 1.0},  // #0C0D45
			"GlowLeft":  [4]float32{00.098, 0.623, 0.863, 1.0}, // #199FDC
			"GlowRight": [4]float32{0.510, 0.176, 0.592, 1.0},  // #822D97
			"Intensity": float32(0.6),                          // Glow boost
			"Spread":    float32(900.0),                        // Plus → glow large

			"Margin":          float32(32.0),
			"CornerRadius":    float32(4.0),
			"BorderThickness": float32(2.0),
			"BorderColor":     [4]float32{.9, .9, .9, 1}, // White border

			"Time": b.Time,
		},
	}

	screen.DrawRectShader(screen.Bounds().Dx(), screen.Bounds().Dy(), b.Shader, opts)
}

func (b *Background) Update() {
	b.Time++
}
