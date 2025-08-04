package control

import (
	"github.com/hajimehoshi/ebiten/v2"
	"math"
	"time"
)

const (
	StateOff = iota
	StateOn
	StateRepeating
)

type Repeater struct {
	das   int // Delay After Start
	arr   int // Auto Repeat Rate
	do    func()
	state int
	ticks int
}

func NewRepeater(das, arr time.Duration, do func()) *Repeater {
	return &Repeater{
		das:   int(math.Ceil(das.Seconds() * float64(ebiten.TPS()))),
		arr:   int(math.Round(arr.Seconds() * float64(ebiten.TPS()))),
		do:    do,
		state: StateOff,
	}
}

func (r *Repeater) Update() {
	if r.state == StateOff {
		return
	}

	r.ticks--
	if r.ticks > 0 {
		return
	}

	if r.state == StateOn {
		r.state = StateRepeating
		r.ticks = r.arr
		r.do()
		return
	}

	// StateRepeating
	r.ticks = r.arr
	r.do()
}

func (r *Repeater) Start() {
	if r.state != StateOff {
		return
	}

	r.do()
	r.state = StateOn
	r.ticks = r.das
}

func (r *Repeater) Stop() {
	r.state = StateOff
}

func (r *Repeater) SetActive(a bool) {
	if a {
		r.Start()
		return
	}
	r.Stop()
}
