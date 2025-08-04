package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"tetris/control"
	"tetris/grid"
	"time"
)

type Controls struct {
	Left  *control.Repeater
	Right *control.Repeater
	Down  func()
	Up    func()
	Reset func()
}

func NewControls(grid *grid.Grid) *Controls {
	return &Controls{
		Left: control.NewRepeater(
			time.Millisecond*170, // DAS
			time.Millisecond*40,  // ARR
			grid.MoveLeft,
		),
		Right: control.NewRepeater(
			time.Millisecond*170, // DAS
			time.Millisecond*40,  // ARR
			grid.MoveRight,
		),
		Down: func() {
			grid.MoveDown()
		},
		Up: func() {
			grid.Rotate()
		},
		Reset: func() {
			grid.Reset()
		},
	}
}

func (c *Controls) Update() {
	c.Left.SetActive(ebiten.IsKeyPressed(ebiten.KeyArrowLeft))
	c.Left.Update()

	c.Right.SetActive(ebiten.IsKeyPressed(ebiten.KeyArrowRight))
	c.Right.Update()

	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		c.Down()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		c.Up()
	}

	// Debug
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		c.Reset()
	}
}
