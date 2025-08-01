package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(1024, 1344)
	err := ebiten.RunGame(NewGame())
	if err != nil {
		panic(err)
	}
}
