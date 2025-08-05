package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err := ebiten.RunGame(
		NewGame(GetAssetsLoader()),
	)
	if err != nil {
		panic(err)
	}
}
