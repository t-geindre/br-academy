package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Toggle struct {
	On  bool
	Key ebiten.Key
}

func NewToggle(key ebiten.Key) *Toggle {
	return &Toggle{
		On:  false,
		Key: key,
	}
}

func (t *Toggle) Update() {
	if inpututil.IsKeyJustPressed(t.Key) {
		t.On = !t.On
	}
}

func (t *Toggle) IsOn() bool {
	return t.On
}
