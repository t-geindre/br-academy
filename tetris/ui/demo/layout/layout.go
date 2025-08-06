package main

import "tetris/ui"

// Minimal target is 720, 820

func GetLayout() *ui.Node {
	m := ui.NewNode(NewComponent("MAIN CONTAINER"))
	m.ContentOrientation = ui.OrientationHorizontal
	m.Padding = [4]int{2, 2, 2, 2}
	m.PaddingUnit = ui.UnitPercentage
	m.Grow = 15
	m.Size = 680

	h := ConstrainNode(1, 5, 1, m, ui.OrientationHorizontal)
	h.Grow = 10
	h.Size = 780

	v := ConstrainNode(1, 10, 1, h, ui.OrientationVertical)

	m.Append(GetPusher(0, 1, 0))

	grid := ui.NewNode(nil)
	grid.Component = NewComponent("GRID")
	grid.Grow = 5
	grid.Size = 588
	m.Append(grid)

	m.Append(GetPusher(0, 2, 1))

	stats := ui.NewNode(nil)
	stats.ContentOrientation = ui.OrientationVertical
	stats.Grow = 5
	stats.Size = 200
	m.Append(stats)

	next := ui.NewNode(nil)
	next.ContentOrientation = ui.OrientationVertical
	next.Grow = 1
	stats.Append(next)

	nextTitle := ui.NewNode(nil)
	nextTitle.Component = NewComponent("Next (T)")
	nextTitle.Size = 100
	next.Append(nextTitle)

	nextValue := ui.NewNode(nil)
	nextValue.Component = NewComponent("Next (V)")
	nextValue.Size = 100
	next.Append(nextValue)

	score := ui.NewNode(nil)
	score.ContentOrientation = ui.OrientationVertical
	score.Grow = 1
	stats.Append(score)

	scoreTitle := ui.NewNode(nil)
	scoreTitle.Component = NewComponent("Score (T)")
	scoreTitle.Size = 100
	score.Append(scoreTitle)

	scoreValue := ui.NewNode(nil)
	scoreValue.Component = NewComponent("Score (V)")
	scoreValue.Size = 100
	score.Append(scoreValue)

	lines := ui.NewNode(nil)
	lines.ContentOrientation = ui.OrientationVertical
	lines.Grow = 1
	stats.Append(lines)

	linesTitle := ui.NewNode(nil)
	linesTitle.Component = NewComponent("Lines (T)")
	linesTitle.Size = 100
	lines.Append(linesTitle)

	linesValue := ui.NewNode(nil)
	linesValue.Component = NewComponent("Lines (V)")
	linesValue.Size = 100
	lines.Append(linesValue)

	m.Append(GetPusher(0, 1, 0))

	return v
}

func ConstrainNode(push int, grow, shrink float64, node *ui.Node, orient uint8) *ui.Node {
	container := ui.NewNode(nil)
	container.ContentOrientation = orient

	container.Append(GetPusher(push, grow, shrink))

	container.Append(node)

	container.Append(GetPusher(push, grow, shrink))

	return container
}

func GetPusher(push int, grow, shrink float64) *ui.Node {
	pusher := ui.NewNode(nil)
	pusher.Size = push
	pusher.Grow = grow
	pusher.Shrink = shrink
	return pusher
}
