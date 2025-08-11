package dsp

import (
	"encoding/binary"
	"io"
	"time"
)

type StreamPulser struct {
	src       io.ReadSeeker
	filter    *FilterBandPass
	roll      []byte
	gate      bool
	gateStart time.Time
	release   time.Duration
	threshold float64
}

func NewStreamPulser(src io.ReadSeeker, bp *FilterBandPass, t float64, r time.Duration) *StreamPulser {
	return &StreamPulser{src: src, filter: bp, threshold: t, release: r}
}

func (t *StreamPulser) Read(p []byte) (int, error) {
	n, err := t.src.Read(p)
	if n <= 0 {
		return n, err
	}

	// Concat rollover + le chunk
	chunk := append(t.roll[:0:0], t.roll...)
	chunk = append(chunk, p[:n]...)

	const frameBytes = 4 // 2 octets L + 2 octets R
	full := (len(chunk) / frameBytes) * frameBytes
	data := chunk[:full]
	t.roll = chunk[full:] // garde l'éventuel reste < 1 frame

	// int16 LE stereo -> mono
	const norm = 0.5 / 32768.0
	for i := 0; i+3 < len(data); i += frameBytes {
		l := int16(binary.LittleEndian.Uint16(data[i:]))
		r := int16(binary.LittleEndian.Uint16(data[i+2:]))
		mono := (float64(l) + float64(r)) * norm
		y := t.filter.ProcessSample(mono)

		if !t.gate &&
			(y > t.threshold || y < -t.threshold) &&
			(t.gateStart.IsZero() || time.Since(t.gateStart) > t.release) {
			t.gate = true
			t.gateStart = time.Now()
		}
	}

	return n, err
}

func (t *StreamPulser) Seek(offset int64, whence int) (int64, error) {
	return t.src.Seek(offset, whence)
}

func (t *StreamPulser) Gate() bool {
	v := t.gate
	t.gate = false

	return v
}

func (t *StreamPulser) SetRelease(r time.Duration) {
	t.release = r
}

func (t *StreamPulser) SetThreshold(threshold float64) {
	t.threshold = threshold
}
