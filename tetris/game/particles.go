package game

import (
	"engine/math"
	"github.com/fogleman/ease"
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type Particles struct {
	shader *ebiten.Shader
	opts   *ebiten.DrawRectShaderOptions
	time   float32
	pulse  *math.PingPong
}

func NewParticles(shader *ebiten.Shader) *Particles {
	p := &Particles{
		pulse: math.NewPingPong(
			ease.InSine, ease.OutSine,
			time.Millisecond*100, time.Millisecond*300,
		),
	}
	p.setShader(shader)
	return p
}

func (p *Particles) Draw(screen *ebiten.Image) {
	screen.DrawRectShader(screen.Bounds().Dx(), screen.Bounds().Dy(), p.shader, p.opts)
}

func (p *Particles) Update() {
	p.time++
	p.opts.Uniforms["Time"] = p.time
	p.opts.Uniforms["Pulse"] = float32(p.pulse.Value())
}

func (p *Particles) setShader(shader *ebiten.Shader) {
	p.shader = shader
	p.opts = &ebiten.DrawRectShaderOptions{
		Uniforms: map[string]any{
			"time": float32(0),
		},
	}
}

func (p *Particles) Pulse() {
	p.pulse.Start()
}
