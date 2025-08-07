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

// PanelPrintf draws a formatted string on the image at the specified position.
// img is the image to draw on,
// pos is one of TopLeft, TopRight, BottomLeft, BottomRight,
// format and are used as in fmt.Printf.
func PanelPrintf(img *ebiten.Image, pos Position, format string, a ...interface{}) {
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

	var x, y int
	switch pos {
	default:
		x, y = 0, 0
	case TopRight:
		x = img.Bounds().Dx() - int(w)
	case BottomLeft:
		y = img.Bounds().Dy() - int(h)
	case BottomRight:
		y = img.Bounds().Dy() - int(h)
		x = img.Bounds().Dx() - int(w)
	}

	vector.DrawFilledRect(
		img, float32(x), float32(y), w, h, BackgroundColor, false,
	)
	ebitenutil.DebugPrintAt(img, str, x+PaddingH, y+PaddingV)
}

func DrawFTPS(img *ebiten.Image) {
	PanelPrintf(img, TopLeft, "FPS %.0f TPS %.0f", ebiten.ActualFPS(), ebiten.ActualTPS())
}
