package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"tetris/ui"
)

type Layout struct {
	layout *ui.Layout
	Root   *ui.Node

	Container *ui.Node
	Grid      *ui.Node

	NextTitle *ui.Node
	NextValue *ui.Node

	ScoreTitle *ui.Node
	ScoreValue *ui.Node

	LevelTitle *ui.Node
	LevelValue *ui.Node

	MinW, MinH int
}

func NewLayout(minW, minH int) *Layout {
	l := &Layout{
		MinW: minW,
		MinH: minH,
	}

	root := l.build()
	l.layout = ui.NewLayout(root)
	l.Root = root

	return l
}

func (l *Layout) Update() {
	ww, wh := ebiten.WindowSize()
	l.layout.Root.W = ww
	l.layout.Root.H = wh
	l.layout.Root.Size = wh
	l.layout.Apply()
}

// Minimal target is 760, 820

func (l *Layout) build() *ui.Node {
	root := ui.NewNode(nil)
	root.ContentOrientation = ui.OrientationVertical

	vert := ui.NewNode(nil)
	vert.ContentOrientation = ui.OrientationHorizontal
	vert.Grow = 0.5
	vert.Size = l.MinH

	root.Append(l.getPusher(0, 1, 1))
	root.Append(vert)
	root.Append(l.getPusher(0, 1, 1))

	l.Container = ui.NewNode(nil)
	l.Container.ContentOrientation = ui.OrientationVertical
	l.Container.Size = l.MinW
	l.Container.Grow = .5
	l.Container.Name = "MAIN CONTAINER"

	vert.Append(l.getPusher(0, 1, 1))
	vert.Append(l.Container)
	vert.Append(l.getPusher(0, 1, 1))

	innerContainer := ui.NewNode(nil)
	l.Container.Append(l.getPusher(1, 1, 1))
	l.Container.Append(innerContainer)
	l.Container.Append(l.getPusher(1, 1, 1))

	innerContainer.Size = l.MinH - 40
	innerContainer.ContentOrientation = ui.OrientationHorizontal
	innerContainer.ContentSpacing = 10
	innerContainer.ContentSpacingUnit = ui.UnitPercentage
	innerContainer.Padding = [4]int{0, 0, 10, 10}
	innerContainer.PaddingUnit = ui.UnitPercentage

	// 360*720
	l.Grid = ui.NewNode(nil)
	l.Grid.Name = "GRID"
	l.Grid.Size = 360

	stats := ui.NewNode(nil)
	stats.ContentOrientation = ui.OrientationVertical
	stats.Grow = 1

	l.NextTitle = ui.NewNode(nil)
	l.NextTitle.Name = "NEXT TITLE"
	l.NextTitle.Size = 100
	stats.Append(l.NextTitle)

	l.NextValue = ui.NewNode(nil)
	l.NextValue.Name = "NEXT VALUE"
	l.NextValue.Size = 100
	stats.Append(l.NextValue)

	stats.Append(l.getPusher(1, 1, 1))

	l.ScoreTitle = ui.NewNode(nil)
	l.ScoreTitle.Name = "SCORE TITLE"
	l.ScoreTitle.Size = 100
	stats.Append(l.ScoreTitle)

	l.ScoreValue = ui.NewNode(nil)
	l.ScoreValue.Name = "SCORE VALUE"
	l.ScoreValue.Size = 100
	stats.Append(l.ScoreValue)

	stats.Append(l.getPusher(1, 1, 1))

	l.LevelTitle = ui.NewNode(nil)
	l.LevelTitle.Name = "LEVEL TITLE"
	l.LevelTitle.Size = 100
	stats.Append(l.LevelTitle)

	l.LevelValue = ui.NewNode(nil)
	l.LevelValue.Name = "LEVEL VALUE"
	l.LevelValue.Size = 100
	stats.Append(l.LevelValue)

	innerContainer.Append(l.Grid)
	innerContainer.Append(stats)

	return root
}

func (l *Layout) getPusher(push int, grow, shrink float64) *ui.Node {
	pusher := ui.NewNode(nil)
	pusher.Size = push
	pusher.Grow = grow
	pusher.Shrink = shrink
	return pusher
}
