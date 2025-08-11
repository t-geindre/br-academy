package math

import (
	"github.com/fogleman/ease"
	"time"
)

type PingPong struct {
	in  ease.Function
	out ease.Function

	ind  time.Duration
	outd time.Duration

	start time.Time
}

func NewPingPong(in, out ease.Function, ind, outd time.Duration) *PingPong {
	return &PingPong{
		in:    in,
		out:   out,
		ind:   ind,
		outd:  outd,
		start: time.Time{},
	}
}

func (p *PingPong) Start() {
	p.start = time.Now()
}

func (p *PingPong) Value() float64 {
	if p.start.IsZero() {
		return 0
	}

	elapsed := time.Since(p.start)

	if elapsed < p.ind {
		return p.in(float64(elapsed) / float64(p.ind))
	}

	if elapsed < p.ind+p.outd {
		return 1 - p.out(float64(elapsed-p.ind)/float64(p.outd))
	}

	p.start = time.Time{} // reset

	return 0
}
