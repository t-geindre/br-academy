package layout

type Component interface {
	SetSize(width, height int)
	SetPosition(x, y int)
}
