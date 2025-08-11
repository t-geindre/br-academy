package pool

import "github.com/hajimehoshi/ebiten/v2"

type Drawable interface {
	Draw(image *ebiten.Image)
}

type Updatable interface {
	Update()
}

type Pool struct {
	drawables  []Drawable
	updatables []Updatable
}

func NewPool() *Pool {
	return &Pool{
		drawables:  make([]Drawable, 0),
		updatables: make([]Updatable, 0),
	}
}

func (p *Pool) Add(items ...any) {
	for _, item := range items {
		if drawable, ok := item.(Drawable); ok {
			p.drawables = append(p.drawables, drawable)
		}
		if updatable, ok := item.(Updatable); ok {
			p.updatables = append(p.updatables, updatable)
		}
	}
}

func (p *Pool) Remove(item any) {
	if drawable, ok := item.(Drawable); ok {
		for i, d := range p.drawables {
			if d == drawable {
				p.drawables = append(p.drawables[:i], p.drawables[i+1:]...)
				break
			}
		}
	}
	if updatable, ok := item.(Updatable); ok {
		for i, u := range p.updatables {
			if u == updatable {
				p.updatables = append(p.updatables[:i], p.updatables[i+1:]...)
				break
			}
		}
	}
}

func (p *Pool) Update() {
	for _, updatable := range p.updatables {
		updatable.Update()
	}
}

func (p *Pool) Draw(image *ebiten.Image) {
	for _, drawable := range p.drawables {
		drawable.Draw(image)
	}
}
