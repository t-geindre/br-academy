package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"math"
	"ui"
)

const (
	RacketMargin = 10
	RacketWidth  = 10
	RacketHeight = 100
	RacketSpeed  = 5
	BallSize     = 10
	BallAngle    = 1.0472 // rad(100)
	BallSpeed    = 5
)

const (
	StateDead = iota
	StateStarting
	StateAlive
)

const (
	ModeIa = iota
	ModePlayer
)

type Game struct {
	Objects      [3]*Object // L/R racket; ball
	BallSpeed    float64
	State        int
	StateElapsed int // ticks
	GameMode     int // Ia/Player
}

func NewGame() *Game {
	return &Game{
		Objects: [3]*Object{
			// Left racket
			{
				Hitbox: image.Rect(
					RacketMargin, RacketMargin,
					RacketWidth+RacketMargin,
					RacketHeight+RacketMargin,
				),
			},
			// Right racket
			{
				Hitbox: image.Rect(
					RacketMargin, RacketMargin,
					RacketWidth+RacketMargin,
					RacketHeight+RacketMargin,
				),
			},
			// Ball
			{
				// Starts off-screen
				Hitbox:   image.Rect(-100, 50, -100+BallSize, 50+BallSize),
				Velocity: image.Pt(0, 0),
			},
		},
		BallSpeed: 0,
		State:     StateDead,
		GameMode:  ModeIa,
	}
}

func (g *Game) Update() error {
	g.UpdateInput()
	g.UpdatePhysics()
	g.UpdateLogic()
	return nil
}

func (g *Game) UpdateLogic() {
	ww, wh := ebiten.WindowSize()

	if g.State == StateAlive {
		// Check if ball is out of bounds
		if g.Objects[2].Hitbox.Max.X < 0 || g.Objects[2].Hitbox.Min.X > ww {
			g.State = StateDead
			g.StateElapsed = 0
			return
		}
	}

	if g.State == StateDead {
		g.StateElapsed++
		if g.StateElapsed > 60 { // 1 sec
			g.State = StateStarting
			g.StateElapsed = 0
			g.BallSpeed = BallSpeed
			// Reset ball position + velocity
			g.Objects[2].Hitbox = image.Rect(
				ww/2-BallSize/2, wh/2-BallSize/2,
				ww/2+BallSize/2, wh/2+BallSize/2,
			)
			g.Objects[2].Velocity = image.Pt(0, 0)
		}
	}

	if g.State == StateStarting {
		g.StateElapsed++
		if g.StateElapsed > 60 { // 1 sec
			g.State = StateAlive
			g.Objects[2].Velocity = image.Pt(int(g.BallSpeed), 0)
		}
	}
}

func (g *Game) UpdateInput() {
	rr, lr := g.Objects[0], g.Objects[1]

	// Left racket
	if g.GameMode == ModeIa {
		// Ai driven
		rr.Velocity = image.Pt(0, 0)
		if g.State == StateAlive {
			bl := g.Objects[2]
			rrc := float64(rr.Hitbox.Min.Y+rr.Hitbox.Max.Y) / 2

			var target float64
			if bl.Velocity.X > 0 {
				_, wh := ebiten.WindowSize()
				target = float64(wh) / 2
			} else {
				target = float64(bl.Hitbox.Min.Y+bl.Hitbox.Max.Y) / 2
			}

			diff := math.Abs(target - rrc)
			if diff > 10 {
				if target < rrc {
					// Ball is above racket
					rr.Velocity = image.Pt(0, -RacketSpeed)
				} else if target > rrc {
					// Ball is below racket
					rr.Velocity = image.Pt(0, RacketSpeed)
				}
			}
		}
	} else {
		// Player driven
		rrv := image.Pt(0, 0)
		if ebiten.IsKeyPressed(ebiten.KeyQ) {
			rrv = rrv.Add(image.Pt(0, -1*RacketSpeed))
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) {
			rrv = rrv.Add(image.Pt(0, 1*RacketSpeed))
		}
		rr.Velocity = rrv
	}

	// Right racket, player driven
	llv := image.Pt(0, 0)
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		llv = llv.Add(image.Pt(0, -1*RacketSpeed))
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		llv = llv.Add(image.Pt(0, 1*RacketSpeed))
	}
	lr.Velocity = llv

	// Toggle game mode
	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		if g.GameMode == ModeIa {
			g.GameMode = ModePlayer
		} else {
			g.GameMode = ModeIa
		}
	}
}

func (g *Game) UpdatePhysics() {
	ww, wh := ebiten.WindowSize()
	rr, lr, bl := g.Objects[0], g.Objects[1], g.Objects[2]

	// Fix right racket position
	if lr.Hitbox.Max.X != ww-RacketMargin {
		lr.Hitbox = image.Rect(ww-RacketMargin-RacketWidth, RacketMargin, ww-RacketMargin, RacketHeight)
	}

	// Make objects move
	for _, obj := range g.Objects {
		obj.Hitbox = obj.Hitbox.Add(obj.Velocity)
	}

	for _, r := range []*Object{rr, lr} {
		// Keep rackets within window bounds
		if r.Hitbox.Min.Y < RacketMargin {
			r.Hitbox.Min.Y = RacketMargin
			r.Hitbox.Max.Y = RacketHeight + RacketMargin
		} else if r.Hitbox.Max.Y > wh-RacketMargin {
			r.Hitbox.Min.Y = wh - RacketHeight - RacketMargin
			r.Hitbox.Max.Y = wh - RacketMargin
		}

		// Make ball bounce on rackets
		inter := bl.Hitbox.Intersect(r.Hitbox)
		if !inter.Empty() {
			// Racket & ball center Y
			bc := float64(bl.Hitbox.Min.Y+bl.Hitbox.Max.Y) / 2
			rc := float64(r.Hitbox.Min.Y+r.Hitbox.Max.Y) / 2

			// Dist, normalized
			relative := (bc - rc) / (float64(r.Hitbox.Dy()) / 2)

			// Angle of the ball relative to the racket
			angle := relative * BallAngle

			// New velocity
			newVX := math.Cos(angle) * g.BallSpeed
			newVY := math.Sin(angle) * g.BallSpeed

			// Apply new velocity + ball position fix
			if bl.Velocity.X > 0 {
				bl.Hitbox = bl.Hitbox.Sub(image.Pt(inter.Dx(), 0))
				newVX = -math.Abs(newVX)
			} else {
				bl.Hitbox = bl.Hitbox.Add(image.Pt(inter.Dx(), 0))
				newVX = math.Abs(newVX)
			}

			bl.Velocity = image.Pt(int(newVX), int(newVY))

			// Increase ball speed
			g.BallSpeed += 0.5
		}
	}

	// Make ball bounce on top/bottom
	if bl.Hitbox.Min.Y <= RacketMargin || bl.Hitbox.Min.Y >= wh-RacketMargin {
		bl.Velocity = image.Pt(bl.Velocity.X, -bl.Velocity.Y)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	cols := []color.Color{
		colornames.Red,
		colornames.Green,
		colornames.Blue,
	}

	for i, obj := range g.Objects {
		vector.DrawFilledRect(
			screen,
			float32(obj.Hitbox.Min.X), float32(obj.Hitbox.Min.Y),
			float32(obj.Hitbox.Dx()), float32(obj.Hitbox.Dy()),
			cols[i],
			false,
		)
	}

	mode := "IA"
	if g.GameMode == ModePlayer {
		mode = "Player"
	}

	ui.PanelPrintf(screen, ui.TopRight, "[F1] Mode: %s", mode)
	ui.DrawFTPS(screen)
}

func (g *Game) Layout(w, h int) (int, int) {
	return w, h
}
