package ui

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"image/color"
	"strings"
)

type Position int

var BackgroundColor = color.RGBA{A: 0x88}
var PaddingH = 10
var PaddingV = 5

const (
	cw               = 6
	ch               = 16
	TopLeft Position = iota
	TopRight
	BottomLeft
	BottomRight
)

// DrawPanel draws a formatted string on the image at the specified position.
// img is the image to draw on,
// pos is one of TopLeft, TopRight, BottomLeft, BottomRight,
// format and are used as in fmt.Printf.
func DrawPanel(img *ebiten.Image, pos Position, format string, a ...interface{}) {
	w, h, str := computePanel(format, a...)

	var x, y float32
	switch pos {
	default:
		x, y = 0, 0
	case TopRight:
		x = float32(img.Bounds().Dx()) - w
	case BottomLeft:
		y = float32(img.Bounds().Dy()) - h
	case BottomRight:
		y = float32(img.Bounds().Dy()) - h
		x = float32(img.Bounds().Dx()) - w
	}

	drawPanel(img, x, y, w, h, str)
}

func DrawPanelAt(img *ebiten.Image, x, y float32, format string, a ...interface{}) {
	w, h, str := computePanel(format, a...)
	drawPanel(img, x, y, w, h, str)
}

func computePanel(format string, a ...interface{}) (float32, float32, string) {
	str := fmt.Sprintf(format, a...)
	h, w := float32(0), float32(0)
	for _, l := range strings.Split(str, "\n") {
		h += ch
		ln := float32(len(l)) * cw
		if ln > w {
			w = ln
		}
	}

	w += float32(PaddingH) * 2
	h += float32(PaddingV) * 2

	return w, h, str
}

func drawPanel(img *ebiten.Image, x, y float32, w, h float32, str string) {
	vector.DrawFilledRect(
		img, x, y, w, h, BackgroundColor, false,
	)
	ebitenutil.DebugPrintAt(img, str, int(x)+PaddingH, int(y)+PaddingV)
}
