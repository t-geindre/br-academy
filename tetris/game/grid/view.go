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

	OffsetX, OffsetY int

	Clearings []*Clearing
	TempLine  *ebiten.Image

	DisappearShd *ebiten.Shader
	GridShd      *ebiten.Shader
}

func NewView(g *Grid, spacing, padding int, b *ebiten.Image, dsh, gsh *ebiten.Shader) *View {
	bds := b.Bounds()

	return &View{
		Grid:       g,
		Brick:      b,
		BrickSize:  bds.Dx(),
		BrickSpace: spacing,
		BrickPad:   padding,

		TempLine: ebiten.NewImage(
			bds.Dx()*g.W-padding*(g.W*2-1)+spacing*(g.W-1),
			bds.Dx(),
		),

		DisappearShd: dsh,
		GridShd:      gsh,
	}
}

func (v *View) Draw(screen *ebiten.Image) {
	v.DrawBackground(screen)
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
			if opts == nil {
				continue
			}
			opts.GeoM.Reset()
			opts.GeoM.Translate(float64(x*(v.BrickPad+v.BrickSpace)), 0)
			v.TempLine.DrawImage(v.Brick, opts)
		}

		bds := v.TempLine.Bounds()
		rx, ry := float64(v.OffsetX-v.BrickPad), float64(v.OffsetY-v.BrickPad+c.Line*(v.BrickPad+v.BrickSpace))

		thr := float64(c.Ticks) / float64(v.Grid.Stats.GetTickRate())
		thr = math.Max(math.Min(thr, float64(bds.Dx())), 0) // Clamp, [0, bds.Dx()]

		opts := &ebiten.DrawRectShaderOptions{
			Uniforms: map[string]interface{}{
				"Threshold": float32(thr),
			},
			Images: [4]*ebiten.Image{
				v.TempLine,
			},
		}
		opts.GeoM.Translate(rx, ry)

		screen.DrawRectShader(bds.Dx(), bds.Dy(), v.DisappearShd, opts)
	}
}

func (v *View) DrawCenteredTetriminoAt(screen *ebiten.Image, t *FallingTetrimino, cx, cy float64) {
	shape := t.Shape.Shapes[t.RotIdx]
	opts := t.Shape.Color
	step := float64(v.BrickPad + v.BrickSpace)

	for i, p := range shape.S {
		if p == 0 {
			continue
		}

		col := float64(i % 4)
		row := float64(i / 4)

		// Position centrée, en tenant compte du barycentre (grille logique)
		x := cx + (col-shape.CX)*step
		y := cy + (row-shape.CY)*step

		// Décalage unique pour corriger la marge graphique du sprite
		x -= float64(v.BrickPad)
		y -= float64(v.BrickPad)

		v.DrawBrickAt(screen, x, y, opts)
	}
}

func (v *View) DrawBrickAtGridPos(screen *ebiten.Image, x, y int, opts *ebiten.DrawImageOptions) {
	if y < 0 {
		return
	}

	v.DrawBrickAt(
		screen,
		float64(v.OffsetX-v.BrickPad+x*(v.BrickPad+v.BrickSpace)),
		float64(v.OffsetY-v.BrickPad+y*(v.BrickPad+v.BrickSpace)),
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

func (v *View) DrawBackground(screen *ebiten.Image) {
	offsetX := float64(v.OffsetX - v.BrickSpace/2)
	offsetY := float64(v.OffsetY - v.BrickSpace/2)

	w := (v.BrickSize - (v.BrickPad * 2) + v.BrickSpace) * v.Grid.W
	h := (v.BrickSize - (v.BrickPad * 2) + v.BrickSpace) * v.Grid.H

	opts := &ebiten.DrawRectShaderOptions{
		Uniforms: map[string]interface{}{
			"CellSize":  float32(v.BrickSize - (v.BrickPad * 2) + v.BrickSpace),
			"LineWidth": float32(1.0),
			"LineColor": [4]float32{.08, .08, .08, .01},
			"Offset":    [2]float64{offsetX, offsetY},
		},
	}
	opts.GeoM.Translate(offsetX, offsetY)

	screen.DrawRectShader(w, h, v.GridShd, opts)
}

func (v *View) SetSize(width, height int) {
}

func (v *View) SetPosition(x, y int) {
	v.OffsetX, v.OffsetY = x, y
}
