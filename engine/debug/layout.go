package debug

import (
	"engine/layout"
	"engine/ui"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"image/color"
	"math/rand"
	"time"
)

func DrawLayoutNode(screen *ebiten.Image, n *layout.Node) {
	col := getNodeColor(n)

	mx, my := ebiten.CursorPosition()
	if n.X <= mx && mx < n.X+n.W && n.Y <= my && my < n.Y+n.H {
		str := ""
		if n.Name != "" {
			str = fmt.Sprintf("%s\n", n.Name)
		}
		str = fmt.Sprintf(
			"%sG %.2f S %.2f\nSi %d (%d,%d) Po[%d,%d]",
			str, n.Grow, n.Shrink, n.Size, n.W, n.H, n.X, n.Y,
		)
		ui.DrawPanelAt(screen, float32(n.X)+5, float32(n.Y)+5, str)
		col = colornames.White
	}

	vector.StrokeRect(
		screen,
		float32(n.X+2), float32(n.Y+2),
		float32(n.W-3), float32(n.H-3),
		3, col, false,
	)

	for _, c := range n.Children {
		DrawLayoutNode(screen, c)
	}
}

type nodeColor struct {
	lastSeen time.Time
	color    color.Color
}

var nodeColors map[*layout.Node]*nodeColor

func getNodeColor(n *layout.Node) color.Color {
	if nodeColors == nil {
		nodeColors = make(map[*layout.Node]*nodeColor)
	} else {
		clearNodeColors()
	}

	if nodeColors[n] == nil {
		nodeColors[n] = &nodeColor{
			color: color.RGBA{
				R: uint8(rand.Intn(150) + 55),
				G: uint8(rand.Intn(150) + 55),
				B: uint8(rand.Intn(150) + 55),
				A: 255,
			},
		}
	}

	nodeColors[n].lastSeen = time.Now()

	return nodeColors[n].color
}

func clearNodeColors() {
	for i, c := range nodeColors {
		if time.Since(c.lastSeen) > time.Second*10 {
			delete(nodeColors, i)
		}
	}
}
