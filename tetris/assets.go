package main

import "tetris/assets"

func GetAssetsLoader() *assets.Loader {
	loader := assets.NewLoader()

	// Images
	loader.AddImage("background", "assets/imgs/ui.png")
	loader.AddImage("brick", "assets/imgs/brick_glow.png")

	// Shaders
	loader.AddShader("disappear", "assets/shaders/disappear.kage")
	loader.AddShader("grid", "assets/shaders/grid.kage")
	loader.AddShader("background", "assets/shaders/background.kage")

	// Fonts
	loader.AddFont("normal", "assets/fonts/normal.ttf")
	loader.AddFont("bold", "assets/fonts/bold.ttf")

	return loader
}
