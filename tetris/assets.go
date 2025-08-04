package main

import "tetris/assets"

func GetAssetsLoader() *assets.Loader {
	loader := assets.NewLoader()
	loader.AddImage("background", "assets/imgs/ui.png")
	loader.AddImage("brick", "assets/imgs/brick_glow.png")
	loader.AddShader("disappear", "assets/shaders/disappear.kage")
	loader.AddShader("grid", "assets/shaders/grid.kage")

	return loader
}
