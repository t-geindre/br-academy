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
	loader.AddRaw("settings", "assets/settings.json")
	loader.MustLoad()

	settings := &game.Settings{}
	setLoader := asset.NewSettings(loader, "settings")
	setLoader.Load(settings)
	defer setLoader.Persist()

	ebiten.SetWindowSize(settings.WinSize[0], settings.WinSize[1])
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Kick detector")

	_ = ebiten.RunGame(game.NewGame(loader, settings))
}
