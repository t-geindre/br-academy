package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"tetris/assets"
)

func main() {
	loader := assets.NewLoader()
	loader.AddImage("background", "assets/imgs/ui.png")
	loader.AddImage("brick", "assets/imgs/brick.png")
	loader.AddShader("disappear", "assets/shaders/disappear.kage")

	err := ebiten.RunGame(NewGame(loader))
	if err != nil {
		panic(err)
	}
}
