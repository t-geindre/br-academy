package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"tetris/game"
)

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err := ebiten.RunGame(
		game.NewGame(GetAssetsLoader()),
	)
	if err != nil {
		panic(err)
	}
}
