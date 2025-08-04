package grid

import "github.com/hajimehoshi/ebiten/v2"

type Clearing struct {
	Line     int
	TicksEnd int
	Ticks    int
}

type View struct {
	Grid       *Grid
	Brick      *ebiten.Image
	BrickSize  int
	BrickPad   int
	BrickSpace int
	Offset     int
	Clearings  []*Clearing
	TempLine   *ebiten.Image
	Time       float32

	// Shaders
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

		Offset:   0,
		TempLine: ebiten.NewImage(bSize, bSize*g.W+bSpace*(g.W-1)),

		DisappearShd: dsh,
		GridShd:      gsh,
	}
}

func (v *View) Draw(screen *ebiten.Image) {
	for y := 0; y < v.Grid.H; y++ {
		if v.isClearing(y) {
			continue
		}

		for x := 0; x < v.Grid.W; x++ {
			opts := v.Grid.Bricks[y*v.Grid.W+x]
			if opts == nil {
				continue
			}
			v.DrawAtGridPos(screen, x, y, opts)
		}
	}

	// Draw active one
	if v.Grid.Active != nil {
		act := v.Grid.Active
		y := 0
		for x, p := range act.Shape.Shapes[act.RotIdx].S {
			if p != 0 {
				v.DrawAtGridPos(screen, act.X+x%4, act.Y+y, v.Grid.Active.Shape.Color)
			}
			if x%4 == 3 {
				y++
			}
		}
	}
}

func (v *View) DrawAtGridPos(screen *ebiten.Image, x, y int, opts *ebiten.DrawImageOptions) {
	opts.GeoM.Reset()
	opts.GeoM.Translate(
		float64(v.Offset-v.BrickPad+x*(v.BrickPad+v.BrickSpace)),
		float64(v.Offset-v.BrickPad+y*(v.BrickPad+v.BrickSpace)),
	)
	screen.DrawImage(v.Brick, opts)
}

func (v *View) isClearing(y int) bool {
	for _, c := range v.Clearings {
		if c.Line == y {
			c.Ticks++
			if c.Ticks > c.TicksEnd {
				v.Clearings = append(v.Clearings[:0], v.Clearings[1:]...)
				return false
			}
			return true
		}
	}

	if v.Grid.IsClearedLine(y) {
		v.Clearings = append(v.Clearings, &Clearing{
			Line:     y,
			Ticks:    0,
			TicksEnd: v.Grid.Stats.GetTickRate(),
		})

		return true
	}

	return false
}
