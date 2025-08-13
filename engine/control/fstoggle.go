package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type FsToggle struct {
}

func NewFsToggle() *FsToggle {
	return &FsToggle{}
}

func (f *FsToggle) Update() {
	if !inpututil.IsKeyJustPressed(ebiten.KeyEnter) || !ebiten.IsKeyPressed(ebiten.KeyAlt) {
		return
	}

	if ebiten.IsFullscreen() {
		ebiten.SetFullscreen(false)
		return
	}

	ebiten.SetFullscreen(true)
}
