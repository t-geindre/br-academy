package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"layout"
)

type Layout struct {
	layout *layout.Layout
	Root   *layout.Node

	Container *layout.Node
	Grid      *layout.Node

	NextTitle *layout.Node
	NextValue *layout.Node

	ScoreTitle *layout.Node
	ScoreValue *layout.Node

	LevelTitle *layout.Node
	LevelValue *layout.Node

	MinW, MinH int
}

func NewLayout(minW, minH int) *Layout {
	l := &Layout{
		MinW: minW,
		MinH: minH,
	}

	root := l.build()
	l.layout = layout.NewLayout(root)
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

func (l *Layout) build() *layout.Node {
	root := layout.NewNode(nil)
	root.ContentOrientation = layout.OrientationVertical

	vert := layout.NewNode(nil)
	vert.ContentOrientation = layout.OrientationHorizontal
	vert.Grow = 0.5
	vert.Size = l.MinH

	root.Append(l.getPusher(0, 1, 1))
	root.Append(vert)
	root.Append(l.getPusher(0, 1, 1))

	l.Container = layout.NewNode(nil)
	l.Container.ContentOrientation = layout.OrientationVertical
	l.Container.Size = l.MinW
	l.Container.Grow = .5
	l.Container.Name = "MAIN CONTAINER"

	vert.Append(l.getPusher(0, 1, 1))
	vert.Append(l.Container)
	vert.Append(l.getPusher(0, 1, 1))

	innerContainer := layout.NewNode(nil)
	l.Container.Append(l.getPusher(1, 1, 1))
	l.Container.Append(innerContainer)
	l.Container.Append(l.getPusher(1, 1, 1))

	innerContainer.Size = l.MinH - 40
	innerContainer.ContentOrientation = layout.OrientationHorizontal
	innerContainer.ContentSpacing = 10
	innerContainer.ContentSpacingUnit = layout.UnitPercentage
	innerContainer.Padding = [4]int{0, 0, 10, 10}
	innerContainer.PaddingUnit = layout.UnitPercentage

	// 360*720
	l.Grid = layout.NewNode(nil)
	l.Grid.Name = "GRID"
	l.Grid.Size = 360

	stats := layout.NewNode(nil)
	stats.ContentOrientation = layout.OrientationVertical
	stats.Grow = 1

	l.NextTitle = layout.NewNode(nil)
	l.NextTitle.Name = "NEXT TITLE"
	l.NextTitle.Size = 100
	stats.Append(l.NextTitle)

	l.NextValue = layout.NewNode(nil)
	l.NextValue.Name = "NEXT VALUE"
	l.NextValue.Size = 100
	stats.Append(l.NextValue)

	stats.Append(l.getPusher(1, 1, 1))

	l.ScoreTitle = layout.NewNode(nil)
	l.ScoreTitle.Name = "SCORE TITLE"
	l.ScoreTitle.Size = 100
	stats.Append(l.ScoreTitle)

	l.ScoreValue = layout.NewNode(nil)
	l.ScoreValue.Name = "SCORE VALUE"
	l.ScoreValue.Size = 100
	stats.Append(l.ScoreValue)

	stats.Append(l.getPusher(1, 1, 1))

	l.LevelTitle = layout.NewNode(nil)
	l.LevelTitle.Name = "LEVEL TITLE"
	l.LevelTitle.Size = 100
	stats.Append(l.LevelTitle)

	l.LevelValue = layout.NewNode(nil)
	l.LevelValue.Name = "LEVEL VALUE"
	l.LevelValue.Size = 100
	stats.Append(l.LevelValue)

	innerContainer.Append(l.Grid)
	innerContainer.Append(stats)

	return root
}

func (l *Layout) getPusher(push int, grow, shrink float64) *layout.Node {
	pusher := layout.NewNode(nil)
	pusher.Size = push
	pusher.Grow = grow
	pusher.Shrink = shrink
	return pusher
}
