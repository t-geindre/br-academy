package grid

import "github.com/hajimehoshi/ebiten/v2"

type Shape struct {
	S       [4 * 4]uint8
	L, R, B int // empty space, Left, Right, Bottom
}
type Shapes [4]Shape

type Tetrimino struct {
	Shapes Shapes
	Color  *ebiten.DrawImageOptions
}

type FallingTetrimino struct {
	Shape    Tetrimino
	RotIdx   int
	X, Y     int // grid pos (top left 4x4)
	IsFixing bool
	Fixing   int // ticks before fixing
	FReset   int // reset count
}

func GetTetriminos() [7]Tetrimino {
	ts := [7]Tetrimino{}

	// I (Cyan)
	ts[0] = Tetrimino{
		Shapes: Shapes{
			toShape(
				"....",
				"XXXX",
				"....",
				"....",
			),
			toShape(
				"..X.",
				"..X.",
				"..X.",
				"..X.",
			),
			toShape(
				"....",
				"....",
				"XXXX",
				"....",
			),
			toShape(
				".X..",
				".X..",
				".X..",
				".X..",
			),
		},
		Color: colorOpts(0.00, 1.00, 1.00, 1.00), // Cyan
	}

	// J (Blue)
	ts[1] = Tetrimino{
		Shapes: Shapes{
			toShape(
				"X..",
				"XXX",
				"...",
				"...",
			),
			toShape(
				".XX",
				".X.",
				".X.",
				"...",
			),
			toShape(
				"...",
				"XXX",
				"..X",
				"...",
			),
			toShape(
				".X.",
				".X.",
				"XX.",
				"...",
			),
		},
		Color: colorOpts(0.00, 0.00, 1.00, 1.00), // Blue
	}

	// L (Orange)
	ts[2] = Tetrimino{
		Shapes: Shapes{
			toShape(
				"..X",
				"XXX",
				"...",
				"...",
			),
			toShape(
				".X.",
				".X.",
				".XX",
				"...",
			),
			toShape(
				"...",
				"XXX",
				"X..",
				"...",
			),
			toShape(
				"XX.",
				".X.",
				".X.",
				"...",
			),
		},
		Color: colorOpts(1.00, 0.65, 0.00, 1.00), // Orange
	}

	// O (Yellow)
	ts[3] = Tetrimino{
		Shapes: Shapes{
			toShape(
				".XX.",
				".XX.",
				"....",
				"....",
			),
			toShape(
				".XX.",
				".XX.",
				"....",
				"....",
			),
			toShape(
				".XX.",
				".XX.",
				"....",
				"....",
			),
			toShape(
				".XX.",
				".XX.",
				"....",
				"....",
			),
		},
		Color: colorOpts(1.00, 1.00, 0.00, 1.00), // Yellow
	}

	// S (Green)
	ts[4] = Tetrimino{
		Shapes: Shapes{
			toShape(
				".XX",
				"XX.",
				"...",
				"...",
			),
			toShape(
				".X.",
				".XX",
				"..X",
				"...",
			),
			toShape(
				"...",
				".XX",
				"XX.",
				"...",
			),
			toShape(
				"X..",
				"XX.",
				".X.",
				"...",
			),
		},
		Color: colorOpts(0.00, 1.00, 0.00, 1.00), // Green
	}

	// T (Purple)
	ts[5] = Tetrimino{
		Shapes: Shapes{
			toShape(
				".X.",
				"XXX",
				"...",
				"...",
			),
			toShape(
				".X.",
				".XX",
				".X.",
				"...",
			),
			toShape(
				"...",
				"XXX",
				".X.",
				"...",
			),
			toShape(
				".X.",
				"XX.",
				".X.",
				"...",
			),
		},
		Color: colorOpts(0.50, 0.00, 0.50, 1.00), // Purple
	}

	// Z (Red)
	ts[6] = Tetrimino{
		Shapes: Shapes{
			toShape(
				"XX.",
				".XX",
				"...",
				"...",
			),
			toShape(
				"..X",
				".XX",
				".X.",
				"...",
			),
			toShape(
				"...",
				"XX.",
				".XX",
				"...",
			),
			toShape(
				".X.",
				"XX.",
				"X..",
				"...",
			),
		},
		Color: colorOpts(1.00, 0.00, 0.00, 1.00), // Red
	}

	return ts
}

// helper to create color
func colorOpts(r, g, b, a float32) *ebiten.DrawImageOptions {
	opt := &ebiten.DrawImageOptions{}
	opt.ColorScale.Scale(r, g, b, a)
	return opt
}

func toShape(rows ...string) Shape {
	var s Shape
	minX, maxX := 4, -1
	maxY := -1

	for y, row := range rows {
		for x, c := range row {
			if c == 'X' {
				s.S[y*4+x] = 1

				// Mises à jour des bornes utiles
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	if maxX >= minX {
		s.L = minX
		s.R = 3 - maxX
		s.B = 3 - maxY
	} else {
		// Forme vide (juste au cas où)
		s.L, s.R, s.B = 4, 4, 4
	}

	return s
}
