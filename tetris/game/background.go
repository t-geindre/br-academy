package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Background struct {
	Shader *ebiten.Shader
	Opts   *ebiten.DrawRectShaderOptions
	Time   float32

	// Glow colors
	glowLeft   [4]float32
	glowRight  [4]float32
	glowDanger [4]float32

	danger float32 // 0.0 - 1.0
	dEase  float32 // 0.0 - 1.0
}

func NewBackground(shader *ebiten.Shader) *Background {
	b := &Background{
		glowLeft:   [4]float32{00.098, 0.623, 0.863, 1.0}, // #199FDC
		glowRight:  [4]float32{0.510, 0.176, 0.592, 1.0},  // #822D97
		glowDanger: [4]float32{1, 0, 0, 1.0},              // RED
	}
	b.setShader(shader)
	return b
}

func (b *Background) Draw(screen *ebiten.Image) {
	screen.DrawRectShader(screen.Bounds().Dx(), screen.Bounds().Dy(), b.Shader, b.Opts)
}

func (b *Background) Update() {
	glowLeft := b.glowLeft
	glowRight := b.glowRight

	if b.dEase != b.danger {
		delta := b.dEase - b.danger
		if delta < 0.01 && delta > -0.01 {
			b.dEase = b.danger
		} else {
			if delta < 0 {
				b.dEase += 0.002
			} else {
				b.dEase -= 0.002
			}
		}
	}

	if b.dEase != 0 {
		for _, glow := range []*[4]float32{&glowLeft, &glowRight} {
			for i := 0; i < 3; i++ { // R, G, B
				glow[i] = glow[i]*(1-b.dEase) + b.glowDanger[i]*b.dEase
			}
		}
	}

	b.Time++
	b.Opts.Uniforms["Time"] = b.Time
	b.Opts.Uniforms["GlowLeft"] = glowLeft
	b.Opts.Uniforms["GlowRight"] = glowRight
}

func (b *Background) setShader(shader *ebiten.Shader) {
	b.Shader = shader
	b.Opts = &ebiten.DrawRectShaderOptions{
		Uniforms: map[string]any{
			"BaseColor":  [4]float32{0.047, 0.051, 0.271, 1.0}, // #0C0D45
			"Intensity":  float32(.9),
			"Spread":     float32(.28),
			"VertSpread": float32(.45),
			"Time":       float32(0),
		},
	}
}

func (b *Background) SetDanger(danger float32) {
	b.danger = danger
}
