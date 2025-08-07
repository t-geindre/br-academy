package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Text struct {
	Content string
	OX, OY  float64
	X, Y    float64
	W, H    float64
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

func (t *Text) SetSize(width, height int) {
	t.W, t.H = float64(width), float64(height)
	t.computePos()
}

func (t *Text) SetPosition(x, y int) {
	t.OX, t.OY = float64(x), float64(y)
	t.computePos()
}

func (t *Text) SetContent(content string) {
	t.Content = content
	t.computePos()
}

func (t *Text) computePos() {
	if t.Content == "" {
		return
	}

	w, h := text.Measure(t.Content, t.Face, 0)

	// Centre dans la zone W,H
	t.X = t.OX + (t.W-w)/2
	t.Y = t.OY + (t.H-h)/2 // important pour la baseline
}
