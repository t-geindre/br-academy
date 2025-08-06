package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Text struct {
	Content string
	X, Y    float64
	Size    float64
	Face    text.Face
}

func NewText(content string, x, y float64, face text.Face) *Text {
	return &Text{
		Content: content,
		X:       x,
		Y:       y,
		Face:    face,
	}
}

func (t *Text) Draw(screen *ebiten.Image) {
	opts := &text.DrawOptions{}
	opts.GeoM.Translate(t.X, t.Y)
	text.Draw(screen, t.Content, t.Face, opts)
}
