package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Particles struct {
	Shader *ebiten.Shader
	Opts   *ebiten.DrawRectShaderOptions
	Time   float32
}

func NewParticles(shader *ebiten.Shader) *Particles {
	p := &Particles{}
	p.setShader(shader)
	return p
}

func (p *Particles) Draw(screen *ebiten.Image) {
	screen.DrawRectShader(screen.Bounds().Dx(), screen.Bounds().Dy(), p.Shader, p.Opts)
}

func (p *Particles) Update() {
	p.Time++
	p.Opts.Uniforms["Time"] = p.Time
}

func (p *Particles) setShader(shader *ebiten.Shader) {
	p.Shader = shader
	p.Opts = &ebiten.DrawRectShaderOptions{
		Uniforms: map[string]any{
			"Time": float32(0),
		},
	}
}
