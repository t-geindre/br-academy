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
	v.Padding = [4]int{5, 5, 5, 5}

	m.Append(GetPusher(0, 1, 0))

	grid := ui.NewNode(nil)
	grid.Component = NewComponent("GRID")
	grid.Grow = 5
	grid.Size = 588
	m.Append(grid)

	m.Append(GetPusher(0, 5, 0))

	score := ui.NewNode(nil)
	score.Component = NewComponent("SCORE")
	score.Grow = 5
	score.Size = 292
	m.Append(score)

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
