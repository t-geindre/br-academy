package game

import (
	"engine/math"
	"github.com/fogleman/ease"
	"github.com/hajimehoshi/ebiten/v2"
	"time"
)

type Shader struct {
	shader *ebiten.Shader
	opts   *ebiten.DrawRectShaderOptions
	start  time.Time
	pulse  *math.PingPong
}

func NewShader(shader *ebiten.Shader) *Shader {
	return &Shader{
		shader: shader,
		opts: &ebiten.DrawRectShaderOptions{
			Uniforms: map[string]any{
				"Time":  float32(0),
				"Pulse": float32(0),
			},
		},
		pulse: math.NewPingPong(
			ease.Linear, ease.Linear,
			time.Millisecond*50, time.Millisecond*350,
		),
		start: time.Now(),
	}
}

func (s *Shader) Draw(screen *ebiten.Image) {
	bds := screen.Bounds()
	screen.DrawRectShader(bds.Dx(), bds.Dy(), s.shader, s.opts)
}

func (s *Shader) Update() {
	s.opts.Uniforms["Pulse"] = float32(s.pulse.Value())
	s.opts.Uniforms["Time"] = float32(time.Since(s.start).Seconds())
}

func (s *Shader) Pulse() {
	s.pulse.Start()
}
