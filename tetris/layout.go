package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"tetris/ui"
)

type Layout struct {
	layout     *ui.Layout
	Container  *ui.Node
	Grid       *ui.Node
	NextTitle  *ui.Node
	MinW, MinH int
}

func NewLayout(minW, minH int) *Layout {
	l := &Layout{
		MinW: minW,
		MinH: minH,
	}
	l.layout = ui.NewLayout(l.build())

	return l
}

func (l *Layout) Update() {
	ww, wh := ebiten.WindowSize()

	if ww < l.MinW || wh < l.MinH {
		ww, wh = l.MinW, l.MinH
	}

	l.layout.Root.W = ww
	l.layout.Root.H = wh
	l.layout.Root.Size = wh
	l.layout.Apply()
}

// Minimal target is 720, 820

func (l *Layout) build() *ui.Node {
	l.Container = ui.NewNode(nil)
	l.Container.ContentOrientation = ui.OrientationHorizontal
	l.Container.Padding = [4]int{2, 2, 2, 2}
	l.Container.PaddingUnit = ui.UnitPercentage
	l.Container.Grow = 15
	l.Container.Size = 680

	h := l.constrainNode(1, 5, 1, l.Container, ui.OrientationHorizontal)
	h.Grow = 10
	h.Size = 780

	v := l.constrainNode(1, 10, 1, h, ui.OrientationVertical)

	l.Container.Append(l.getPusher(0, 1, 0))

	l.Grid = ui.NewNode(nil)
	l.Grid.Grow = 5
	l.Grid.Size = 588
	l.Container.Append(l.Grid)

	l.Container.Append(l.getPusher(0, 2, 1))

	stats := ui.NewNode(nil)
	stats.ContentOrientation = ui.OrientationVertical
	stats.Grow = 5
	stats.Size = 200
	l.Container.Append(stats)

	next := ui.NewNode(nil)
	next.ContentOrientation = ui.OrientationVertical
	next.Grow = 1
	stats.Append(next)

	l.NextTitle = ui.NewNode(nil)
	l.NextTitle.Size = 100
	next.Append(l.NextTitle)

	nextValue := ui.NewNode(nil)
	nextValue.Size = 100
	next.Append(nextValue)

	score := ui.NewNode(nil)
	score.ContentOrientation = ui.OrientationVertical
	score.Grow = 1
	stats.Append(score)

	scoreTitle := ui.NewNode(nil)
	scoreTitle.Size = 100
	score.Append(scoreTitle)

	scoreValue := ui.NewNode(nil)
	scoreValue.Size = 100
	score.Append(scoreValue)

	lines := ui.NewNode(nil)
	lines.ContentOrientation = ui.OrientationVertical
	lines.Grow = 1
	stats.Append(lines)

	linesTitle := ui.NewNode(nil)
	linesTitle.Size = 100
	lines.Append(linesTitle)

	linesValue := ui.NewNode(nil)
	linesValue.Size = 100
	lines.Append(linesValue)

	l.Container.Append(l.getPusher(0, 1, 0))

	return v
}

func (l *Layout) constrainNode(push int, grow, shrink float64, node *ui.Node, orient uint8) *ui.Node {
	container := ui.NewNode(nil)
	container.ContentOrientation = orient
	container.Append(l.getPusher(push, grow, shrink))
	container.Append(node)
	container.Append(l.getPusher(push, grow, shrink))
	return container
}

func (l *Layout) getPusher(push int, grow, shrink float64) *ui.Node {
	pusher := ui.NewNode(nil)
	pusher.Size = push
	pusher.Grow = grow
	pusher.Shrink = shrink
	return pusher
}
