package grid

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
)

type Clearing struct {
	Line  int
	Ticks int
}

type View struct {
	Grid *Grid

	Brick      *ebiten.Image
	BrickSize  int
	BrickPad   int
	BrickSpace int

	Offset int

	Clearings []*Clearing
	TempLine  *ebiten.Image

	DisappearShd *ebiten.Shader
	GridShd      *ebiten.Shader
}

func NewView(g *Grid, ox, oy int, b *ebiten.Image, dsh, gsh *ebiten.Shader) *View {
	bPad := 32
	bSize := 96
	bSpace := 4

	return &View{
		Grid:       g,
		Brick:      b,
		BrickSize:  bSize,
		BrickSpace: bSpace,
		BrickPad:   bPad,

		Offset: 64,
		TempLine: ebiten.NewImage(
			bSize*g.W-bPad*(g.W*2-1)+bSpace*(g.W-1),
			bSize,
		),

		DisappearShd: dsh,
		GridShd:      gsh,
	}
}

func (v *View) Draw(screen *ebiten.Image) {
	v.DrawGrid(screen)
	v.DrawActive(screen)
	v.DrawClearing(screen)
}

func (v *View) Update() {
	v.UpdateClearing()
}

func (v *View) DrawGrid(screen *ebiten.Image) {
	for y := 0; y < v.Grid.H; y++ {
		if v.IsClearing(y) {
			continue
		}

		for x := 0; x < v.Grid.W; x++ {
			opts := v.Grid.Bricks[y*v.Grid.W+x]
			if opts == nil {
				continue
			}
			v.DrawBrickAtGridPos(screen, x, y, opts)
		}
	}
}

func (v *View) DrawActive(screen *ebiten.Image) {
	act := v.Grid.Active
	for i, p := range act.Shape.Shapes[act.RotIdx].S {
		if p != 0 {
			v.DrawBrickAtGridPos(screen, act.X+i%4, act.Y+i/4, v.Grid.Active.Shape.Color)
		}
	}
}

func (v *View) DrawClearing(screen *ebiten.Image) {
	for _, c := range v.Clearings {
		v.TempLine.Clear()
		for x := 0; x < v.Grid.W; x++ {
			opts := v.Grid.Bricks[c.Line*v.Grid.W+x]
			opts.GeoM.Reset()
			opts.GeoM.Translate(float64(x*(v.BrickPad+v.BrickSpace)), 0)
			v.TempLine.DrawImage(v.Brick, opts)
		}

		bds := v.TempLine.Bounds()
		rx, ry := float64(v.Offset-v.BrickPad), float64(v.Offset-v.BrickPad+c.Line*(v.BrickPad+v.BrickSpace))

		thr := float64(c.Ticks) / float64(v.Grid.Stats.GetTickRate()) * float64(bds.Dx())
		thr = math.Max(math.Min(thr, float64(bds.Dx())), 0) // Clamp, [0, bds.Dx()]

		opts := &ebiten.DrawRectShaderOptions{
			Uniforms: map[string]interface{}{
				"OffsetX":   rx,
				"OffsetY":   ry,
				"Threshold": thr,
			},
			Images: [4]*ebiten.Image{
				v.TempLine,
			},
		}
		opts.GeoM.Translate(rx, ry)

		screen.DrawRectShader(bds.Dx(), bds.Dy(), v.DisappearShd, opts)
	}
}

func (v *View) DrawCenteredTetriminoAt(screen *ebiten.Image, t *FallingTetrimino, x, y float64) {
	opts := t.Shape.Color
	for i, p := range t.Shape.Shapes[t.RotIdx].S {
		if p != 0 {
			v.DrawBrickAt(
				screen,
				x+float64(i%4*(v.BrickPad+v.BrickSpace)),
				y+float64(i/4*(v.BrickPad+v.BrickSpace)),
				opts,
			)
		}
	}
}

func (v *View) DrawBrickAtGridPos(screen *ebiten.Image, x, y int, opts *ebiten.DrawImageOptions) {
	v.DrawBrickAt(
		screen,
		float64(v.Offset-v.BrickPad+x*(v.BrickPad+v.BrickSpace)),
		float64(v.Offset-v.BrickPad+y*(v.BrickPad+v.BrickSpace)),
		opts,
	)
}

func (v *View) DrawBrickAt(screen *ebiten.Image, x, y float64, opts *ebiten.DrawImageOptions) {
	opts.GeoM.Reset()
	opts.GeoM.Translate(x, y)
	screen.DrawImage(v.Brick, opts)
}

func (v *View) IsClearing(y int) bool {
	for _, c := range v.Clearings {
		if c.Line == y {
			return true
		}
	}

	return false
}

func (v *View) UpdateClearing() {
	for _, c := range v.Clearings {
		c.Ticks++
		if c.Ticks > v.Grid.Stats.GetTickRate() {
			v.Clearings = append(v.Clearings[:0], v.Clearings[1:]...)
			return
		}
	}

toClear:
	for _, i := range v.Grid.ToClear {
		for _, c := range v.Clearings {
			if c.Line == i {
				continue toClear
			}
		}

		v.Clearings = append(v.Clearings, &Clearing{
			Line:  i,
			Ticks: 0,
		})
	}
}
