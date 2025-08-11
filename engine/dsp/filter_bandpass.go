package dsp

import "math"

// Biquad passe-bande (RBJ), DF2T
type FilterBandPass struct {
	b0, b1, b2 float64
	a1, a2     float64
	z1, z2     float64
}

func NewFilterBandPass(sampleRate int, lowHz, highHz float64) *FilterBandPass {
	b := &FilterBandPass{}
	b.SetBand(sampleRate, lowHz, highHz)
	return b
}

func (b *FilterBandPass) SetBand(sampleRate int, lowHz, highHz float64) {
	if highHz <= lowHz {
		highHz = lowHz + 1 // évite division par zéro
	}
	f0 := math.Sqrt(lowHz * highHz)
	bw := highHz - lowHz
	if bw <= 0 {
		bw = 1
	}
	Q := f0 / bw
	if Q < 1e-6 {
		Q = 1e-6
	}

	w0 := 2 * math.Pi * f0 / float64(sampleRate)
	cosw0 := math.Cos(w0)
	sinw0 := math.Sin(w0)
	alpha := sinw0 / (2 * Q)

	// RBJ: "band-pass (constant skirt gain, peak gain = Q)"
	b0 := alpha
	b1 := 0.0
	b2 := -alpha
	a0 := 1 + alpha
	a1 := -2 * cosw0
	a2 := 1 - alpha

	b.b0 = b0 / a0
	b.b1 = b1 / a0
	b.b2 = b2 / a0
	b.a1 = a1 / a0
	b.a2 = a2 / a0

	b.Reset()
}

func (b *FilterBandPass) Reset() {
	b.z1, b.z2 = 0, 0
}

func (b *FilterBandPass) ProcessSample(x float64) float64 {
	// DF2T
	y := b.b0*x + b.z1
	b.z1 = b.b1*x - b.a1*y + b.z2
	b.z2 = b.b2*x - b.a2*y
	return y
}

func (b *FilterBandPass) ProcessBufferMono32(buf []float32) {
	for i := range buf {
		y := b.ProcessSample(float64(buf[i]))
		buf[i] = float32(y)
	}
}

func (b *FilterBandPass) ProcessBufferStereo32(left, right []float32) {
	n := len(left)
	if len(right) < n {
		n = len(right)
	}
	for i := 0; i < n; i++ {
		yL := b.ProcessSample(float64(left[i]))
		yR := b.ProcessSample(float64(right[i]))
		left[i] = float32(yL)
		right[i] = float32(yR)
	}
}
