package main

import "tetris/assets"

func GetAssetsLoader() *assets.Loader {
	loader := assets.NewLoader()
	loader.AddImage("background", "assets/imgs/ui.png")
	loader.AddImage("brick", "assets/imgs/brick.png")
	loader.AddShader("disappear", "assets/shaders/disappear.kage")

	return loader
}
