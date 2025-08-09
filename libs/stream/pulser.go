package stream

/*

import (
	"github.com/fogleman/ease"
	"io"
	"maths"
)

type Pulser struct {
	in, out ease.Function
	pulses  []int64
	src     io.ReadSeeker
}

func NewPulser(src io.ReadSeeker, in, out ease.Function) *Pulser {
	return &Pulser{
		src:    src,
		in:     in,
		out:    out,
		pulses: make([]int64, 0),
	}
}

func (p *Pulser) AddPules(t int64) {
	p.pulses = append(p.pulses, t)
}

func (p *Pulser) Pulse() float64 {
	return math.PingPong(0, 10, ease.InBack, ease.InOutElastic)
}

func (p *Pulser) Read(b []byte) (n int, err error) {
	return p.src.Read(b)
}

func (p *Pulser) Seek(offset int64, whence int) (int64, error) {
	return p.src.Seek(offset, whence)
}
*/
