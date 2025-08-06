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
	l.ApplyNode(l.Root)
}

func (l *Layout) ApplyNode(n *Node) {
	start, space, opStart, opSpace := l.getStartSpace(n)
	_, _ = opStart, opSpace // Unused, but kept for clarity

	cSpace, totalGrow, totalShrink := l.getContentSumGrowShrink(n)

	spacing := l.applyUnit(n.ContentSpacing, n.ContentSpacingUnit, space)
	cSpace += (len(n.Children) - 1) * spacing

	free := space - cSpace
	cursor := start

	if free > 0 && totalGrow == 0 {
		cursor += free / 2
	}

	for i, child := range n.Children {
		size := child.Size

		if free > 0 && child.Grow > 0 {
			// GROW
			size += int(float64(free) * (child.Grow / totalGrow))
		} else if free < 0 && child.Shrink > 0 {
			// SHRINK
			size += int(float64(free) * (child.Shrink / totalShrink)) // free est négatif
			if size < 0 {
				size = 0 // Évite les tailles négatives
			}
		}

		// Position
		if n.ContentOrientation == OrientationHorizontal {
			child.X = cursor
			child.Y = opStart
			child.W = size
			child.H = opSpace
			cursor += size
		} else {
			child.X = opStart
			child.Y = cursor
			child.W = opSpace
			child.H = size
			cursor += size
		}

		if child.Component != nil {
			child.Component.SetPosition(child.X, child.Y)
			child.Component.SetSize(child.W, child.H)
		}

		// Ajoute l'espacement entre les éléments sauf le dernier
		if i < len(n.Children)-1 {
			cursor += spacing
		}

		l.ApplyNode(child)
	}
}

func (l *Layout) applyUnit(p int, unit uint8, space int) int {
	if unit == UnitPixel {
		return p
	}

	return (p * space) / 100
}

func (l *Layout) getStartSpace(n *Node) (int, int, int, int) {
	space := n.H
	start := n.Y

	opStart := n.X
	opSpace := n.W

	lp := l.applyUnit(n.Padding[0], n.PaddingUnit, space)
	rp := l.applyUnit(n.Padding[1], n.PaddingUnit, space)

	otp := l.applyUnit(n.Padding[2], n.PaddingUnit, opSpace)
	obp := l.applyUnit(n.Padding[3], n.PaddingUnit, opSpace)

	if n.ContentOrientation == OrientationHorizontal {
		space = n.W
		start = n.X

		opStart = n.Y
		opSpace = n.H

		lp = l.applyUnit(n.Padding[2], n.PaddingUnit, space)
		rp = l.applyUnit(n.Padding[3], n.PaddingUnit, space)

		otp = l.applyUnit(n.Padding[0], n.PaddingUnit, opSpace)
		obp = l.applyUnit(n.Padding[1], n.PaddingUnit, opSpace)
	}

	start += lp
	space -= lp + rp

	opStart += otp
	opSpace -= otp + obp

	return start, space, opStart, opSpace
}

func (l *Layout) getContentSumGrowShrink(n *Node) (int, float64, float64) {
	cSizeSum := 0
	totalGrow := 0.0
	totalShrink := 0.0

	for _, child := range n.Children {
		cSizeSum += child.Size
		totalGrow += child.Grow
		totalShrink += child.Shrink
	}

	return cSizeSum, totalGrow, totalShrink
}
