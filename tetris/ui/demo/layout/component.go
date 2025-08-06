package main

type Component struct {
	Name string
}

func NewComponent(name string) Component {
	return Component{
		Name: name,
	}
}

func (c Component) SetSize(width, height int) {
}

func (c Component) SetPosition(x, y int) {
}
