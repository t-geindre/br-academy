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

	LinesTitle *layout.Node
	LinesValue *layout.Node

	MinW, MinH int
}

func NewLayout() *Layout {
	l := &Layout{
		MinW: 600,
		MinH: 760,
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
	root := layout.NewNode()
	root.ContentOrientation = layout.OrientationVertical
	root.Name = "ROOT"

	vert := layout.NewNode()
	vert.ContentOrientation = layout.OrientationHorizontal
	vert.Grow = 0.3
	vert.Size = 800
	vert.Name = "HORIZONTAL CONTAINER"

	root.Append(l.getPusher(1, 1, "TOP PUSHER"))
	root.Append(vert)
	root.Append(l.getPusher(1, 1, "BOTTOM PUSHER"))

	l.Container = layout.NewNode()
	l.Container.ContentOrientation = layout.OrientationVertical
	l.Container.Size = 700
	l.Container.Grow = .3
	l.Container.Name = "MAIN CONTAINER"

	vert.Append(l.getPusher(1, 1, "LEFT PUSHER"))
	vert.Append(l.Container)
	vert.Append(l.getPusher(1, 1, "RIGHT PUSHER"))

	innerContainer := layout.NewNode()
	l.Container.Append(l.getPusher(1, 1, "TOP INNER PUSHER"))
	l.Container.Append(innerContainer)
	l.Container.Append(l.getPusher(1, 1, "BOTTOM INNER PUSHER"))

	innerContainer.Size = l.MinH - 40
	innerContainer.ContentOrientation = layout.OrientationHorizontal

	// 360*720
	l.Grid = layout.NewNode()
	l.Grid.Name = "GRID"
	l.Grid.Size = 360

	stats := layout.NewNode()
	stats.ContentOrientation = layout.OrientationVertical
	stats.Size = 200
	stats.Name = "STATS CONTAINER"

	l.NextTitle = layout.NewNode()
	l.NextTitle.Name = "NEXT TITLE"
	l.NextTitle.Size = 80

	l.NextValue = layout.NewNode()
	l.NextValue.Name = "NEXT VALUE"
	l.NextValue.Size = 80

	l.ScoreTitle = layout.NewNode()
	l.ScoreTitle.Name = "SCORE TITLE"
	l.ScoreTitle.Size = 80

	l.ScoreValue = layout.NewNode()
	l.ScoreValue.Name = "SCORE VALUE"
	l.ScoreValue.Size = 80

	l.LevelTitle = layout.NewNode()
	l.LevelTitle.Name = "LEVEL TITLE"
	l.LevelTitle.Size = 80

	l.LevelValue = layout.NewNode()
	l.LevelValue.Name = "LEVEL VALUE"
	l.LevelValue.Size = 80

	l.LinesTitle = layout.NewNode()
	l.LinesTitle.Name = "LINES TITLE"
	l.LinesTitle.Size = 80

	l.LinesValue = layout.NewNode()
	l.LinesValue.Name = "LINES VALUE"
	l.LinesValue.Size = 80

	stats.Append(l.NextTitle)
	stats.Append(l.NextValue)
	stats.Append(l.getPusher(1, 1, "STATS PUSHER"))
	stats.Append(l.LevelTitle)
	stats.Append(l.LevelValue)
	stats.Append(l.getPusher(1, 1, "STATS PUSHER"))
	stats.Append(l.LinesTitle)
	stats.Append(l.LinesValue)
	stats.Append(l.getPusher(1, 1, "STATS PUSHER"))
	stats.Append(l.ScoreTitle)
	stats.Append(l.ScoreValue)

	innerContainer.Append(l.getPusher(1, 1, "INNER PUSHER"))
	innerContainer.Append(l.Grid)
	innerContainer.Append(l.getPusher(1, 1, "INNER PUSHER"))
	innerContainer.Append(stats)
	innerContainer.Append(l.getPusher(1, 1, "INNER PUSHER"))

	return root
}

func (l *Layout) getPusher(grow, shrink float64, n string) *layout.Node {
	pusher := layout.NewNode()
	pusher.Grow = grow
	pusher.Shrink = shrink
	pusher.Name = n
	return pusher
}
