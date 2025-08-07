package debug

import (
	"github.com/hajimehoshi/ebiten/v2"
	"ui"
)

func DrawFTPS(img *ebiten.Image) {
	DrawFtpsPos(img, ui.TopLeft)
}

func DrawFtpsPos(img *ebiten.Image, pos ui.Position) {
	ui.DrawPanel(img, pos, "FPS %.0f TPS %.0f", ebiten.ActualFPS(), ebiten.ActualTPS())
}

func DrawResolution(screen *ebiten.Image) {
	DrawResolutionPos(screen, ui.TopRight)
}

func DrawResolutionPos(screen *ebiten.Image, position ui.Position) {
	ww, wh := ebiten.WindowSize()
	mw, mh := ebiten.Monitor().Size()
	ui.DrawPanel(screen, position, "Window: %d x %d\nMonitor: %d x %d", ww, wh, mw, mh)
}

func DrawAll(screen *ebiten.Image) {
	DrawFTPS(screen)
	DrawResolution(screen)
}
