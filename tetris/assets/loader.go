package assets

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	TypeImage = iota
	TypeShader
	TypeFont
)

type toLoad struct {
	Path, Name string
	Type       int
}

type Loader struct {
	Shaders map[string]*ebiten.Shader
	Images  map[string]*ebiten.Image
	Fonts   map[string]*text.GoTextFaceSource
	ToLoad  []toLoad
	Loaded  bool
}

func NewLoader() *Loader {
	return &Loader{
		Shaders: make(map[string]*ebiten.Shader),
		Images:  make(map[string]*ebiten.Image),
		Fonts:   make(map[string]*text.GoTextFaceSource),
	}
}

func (l *Loader) AddImage(name, path string) {
	l.Add(name, path, TypeImage)
}

func (l *Loader) AddShader(name, path string) {
	l.Add(name, path, TypeShader)
}

func (l *Loader) AddFont(name, path string) {
	l.Add(name, path, TypeFont)
}

func (l *Loader) Add(name, path string, t int) {
	l.ToLoad = append(l.ToLoad, toLoad{
		Name: name,
		Path: path,
		Type: t,
	})
}

func (l *Loader) MustLoad() {
	if l.Loaded {
		panic("Loader already loaded")
	}

	for _, item := range l.ToLoad {
		switch item.Type {
		case TypeImage:
			l.Images[item.Name] = MustLoadImage(item.Path)
		case TypeShader:
			l.Shaders[item.Name] = MustLoadShader(item.Path)
		case TypeFont:
			l.Fonts[item.Name] = MustLoadFont(item.Path)
		default:
			panic("Unknown asset type")
		}
	}

	l.Loaded = true
}

func (l *Loader) GetImage(name string) *ebiten.Image {
	img, ok := l.Images[name]
	if !ok {
		panic("Image not found: " + name)
	}

	return img
}

func (l *Loader) GetShader(name string) *ebiten.Shader {
	shd, ok := l.Shaders[name]
	if !ok {
		panic("Shader not found: " + name)
	}

	return shd
}

func (l *Loader) GetFont(name string) *text.GoTextFaceSource {
	font, ok := l.Fonts[name]
	if !ok {
		panic("Font not found: " + name)
	}

	return font
}
