package grid

import "github.com/hajimehoshi/ebiten/v2"

type View struct {
	Grid             *Grid
	Brick            *ebiten.Image
	BrickW, BrickH   int
	OffsetX, OffsetY int
}

func NewView(g *Grid, ox, oy int, b *ebiten.Image) *View {
	bds := b.Bounds()
	return &View{
		Grid:    g,
		Brick:   b,
		BrickW:  bds.Dx(),
		BrickH:  bds.Dy(),
		OffsetX: ox,
		OffsetY: oy,
	}
}

func (v *View) Draw(screen *ebiten.Image) {
	for y := 0; y < v.Grid.H; y++ {
		for x := 0; x < v.Grid.W; x++ {
			i := y*v.Grid.W + x
			opts := v.Grid.Bricks[i]
			if opts != nil {
				v.DrawBrickAtGridPos(screen, x, y, opts)
			}
		}
	}

	if v.Grid.Active != nil {
		v.DrawTetriminoAtGridPos(screen, v.Grid.Active, v.Grid.Active.X, v.Grid.Active.Y)
	}
}

func (v *View) DrawTetriminoAt(dest *ebiten.Image, tetrimino *FallingTetrimino, x, y float64) {
	shape := tetrimino.Shape.Shapes[tetrimino.RotIdx]

	for row := 0; row < 4; row++ {
		for col := 0; col < 4; col++ {
			if shape.S[row*4+col] != 0 {
				bx := x + float64(col*v.BrickW)
				by := y + float64(row*v.BrickH)

				if by < float64(v.OffsetY) {
					continue
				}

				v.drawBrick(dest, bx, by, tetrimino.Shape.Color)
			}
		}
	}
}

func (v *View) DrawCenteredTetriminoAt(dest *ebiten.Image, t *FallingTetrimino, centerX, centerY float64) {
	shape := t.Shape.Shapes[t.RotIdx]

	w := 4 - shape.L - shape.R
	h := 4 - shape.T - shape.B

	offsetX := centerX - float64(w*v.BrickW)/2 - float64(shape.L*v.BrickW)
	offsetY := centerY - float64(h*v.BrickH)/2 - float64(shape.T*v.BrickH)

	v.DrawTetriminoAt(dest, t, offsetX, offsetY)
}

func (v *View) DrawTetriminoAtGridPos(dest *ebiten.Image, tetrimino *FallingTetrimino, gx, gy int) {
	x := float64(gx*v.BrickW + v.OffsetX)
	y := float64(gy*v.BrickH + v.OffsetY)
	v.DrawTetriminoAt(dest, tetrimino, x, y)
}

func (v *View) DrawBrickAtGridPos(dest *ebiten.Image, gx, gy int, opts *ebiten.DrawImageOptions) {
	x := float64(gx*v.BrickW + v.OffsetX)
	y := float64(gy*v.BrickH + v.OffsetY)
	v.drawBrick(dest, x, y, opts)
}

func (v *View) drawBrick(dest *ebiten.Image, x, y float64, opts *ebiten.DrawImageOptions) {
	opts.GeoM.Reset()
	opts.GeoM.Translate(x, y)
	dest.DrawImage(v.Brick, opts)
}
