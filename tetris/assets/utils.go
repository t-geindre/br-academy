package assets

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"os"
)

func MustLoadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)

	if err != nil {
		panic(err)
	}

	return img
}

func MustLoadShader(path string) *ebiten.Shader {
	raw, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	shader, err := ebiten.NewShader(raw)
	if err != nil {
		panic(err)
	}

	return shader
}
