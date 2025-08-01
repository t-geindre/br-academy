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

func (w *View) Draw(screen *ebiten.Image) {
	// Full grid
	for y := 0; y < w.Grid.H; y++ {
		for x := 0; x < w.Grid.W; x++ {
			i := y*w.Grid.W + x
			if w.Grid.Bricks[i] != nil {
				x, y := float64(x*w.BrickW+w.OffsetX), float64(y*w.BrickH+w.OffsetY)
				opts := w.Grid.Bricks[i]
				opts.GeoM.Reset()
				opts.GeoM.Translate(x, y)
				screen.DrawImage(w.Brick, opts)
			}
		}
	}

	// Active tetrimino
	w.DrawTetriminoAt(screen, w.Grid.Active, float64(w.Grid.Active.X*w.BrickW+w.OffsetX), float64(w.Grid.Active.Y*w.BrickH+w.OffsetY))
}

func (w *View) DrawTetriminoAt(dest *ebiten.Image, tetrimino *FallingTetrimino, x, y float64) {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if tetrimino.Shape.Shapes[tetrimino.RotIdx].S[i*4+j] != 0 {

				bx, by := x+float64(j*w.BrickW), y+float64(i*w.BrickH)
				if by < float64(w.OffsetY) {
					continue
				}

				opts := tetrimino.Shape.Color
				opts.GeoM.Reset()
				opts.GeoM.Translate(bx, by)
				dest.DrawImage(w.Brick, opts)
			}
		}
	}
}
