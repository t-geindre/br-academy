package ui

type Layout struct {
	Root     *Node
	Rlw, Rlh int // Root last known size
}

func NewLayout(root *Node) *Layout {
	return &Layout{
		Root: root,
	}
}

func (l *Layout) Apply() {
	if l.Root == nil {
		return
	}

	if l.Rlw == l.Root.W && l.Rlh == l.Root.H {
		return
	}

	l.Rlw, l.Rlh = l.Root.W, l.Root.H
	for _, child := range l.Root.Children {
		l.ApplyNode(child)
	}
}

func (l *Layout) ApplyNode(n *Node) {
	start, size := 0, 0

	space := l.Root.H
	if n.Parent.Orientation == OrientationHorizontal {
		space = l.Root.W
	}

	if n.Parent.Padding[0] > 0 {
		if n.Parent.PaddingUnit == UnitPixel {
			start = n.Parent.Padding[0]
			size -= n.Parent.Padding[0]
		} else {
			start = (n.Parent.Padding[0] * space) / 100
			size -= (n.Parent.Padding[0] * space) / 100
		}
	}

	if n.Parent.Padding[1] > 0 {
		if n.Parent.PaddingUnit == UnitPixel {
			size -= n.Parent.Padding[1]
		} else {
			size -= (n.Parent.Padding[1] * space) / 100
		}
	}

	if n.MinSize > 0 && size < n.MinSize {
		size = n.MinSize
	}

	if n.MaxSize > 0 && size > n.MaxSize {
		size = n.MaxSize
	}

	if n.Orientation == OrientationHorizontal {
		n.W = size
		n.H = n.Parent.H
		n.X = start
		n.Y = 0
	} else {
		n.W = n.Parent.W
		n.H = size
		n.X = 0
		n.Y = start
	}
}
