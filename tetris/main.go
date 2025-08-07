package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"tetris/game"
)

func main() {
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	err := ebiten.RunGame(
		game.NewGame(game.GetAssetsLoader()),
	)
	if err != nil {
		panic(err)
	}
}
