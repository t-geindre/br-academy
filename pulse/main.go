package main

import (
	"engine/asset"
	"github.com/hajimehoshi/ebiten/v2"
	"pulse/game"
)

func main() {
	loader := asset.NewLoader()
	loader.AddShader("pulse", "assets/pulse.kage")
	loader.AddRaw("ost", "assets/trance.mp3")
	loader.MustLoad()

	ebiten.SetWindowSize(300, 300)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Kick detector")

	_ = ebiten.RunGame(game.NewGame(loader))
}
