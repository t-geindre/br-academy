package assets

import "github.com/hajimehoshi/ebiten/v2"

const (
	TypeImage = iota
	TypeShader
)

type toLoad struct {
	Path, Name string
	Type       int
}

type Loader struct {
	Shaders map[string]*ebiten.Shader
	Images  map[string]*ebiten.Image
	ToLoad  []toLoad
	Loaded  bool
}

func NewLoader() *Loader {
	return &Loader{
		Shaders: make(map[string]*ebiten.Shader),
		Images:  make(map[string]*ebiten.Image),
	}
}

func (l *Loader) AddImage(name, path string) {
	l.Add(name, path, TypeImage)
}

func (l *Loader) AddShader(name, path string) {
	l.Add(name, path, TypeShader)
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
