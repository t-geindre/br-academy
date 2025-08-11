package main

import (
	"engine/asset"
	"github.com/hajimehoshi/ebiten/v2"
	"pulse/game"
)

func main() {
	loader := asset.NewLoader()
	loader.AddShader("pulse", "assets/pulse.kage")
	loader.AddRaw("ost", "assets/st.mp3")
	loader.MustLoad()

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Kick detector")

	_ = ebiten.RunGame(game.NewGame(loader))
}
