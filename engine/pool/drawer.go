package pool

import "github.com/hajimehoshi/ebiten/v2"

type Drawer struct {
	draw func(image *ebiten.Image)
}

func NewDrawer(drawFunc func(image *ebiten.Image)) *Drawer {
	return &Drawer{
		draw: drawFunc,
	}
}

func (d *Drawer) Draw(image *ebiten.Image) {
	d.draw(image)
}
