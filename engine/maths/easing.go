package math

import "github.com/fogleman/ease"

// PingPong t current time, d duration
func PingPong(t, d float64, in, out ease.Function) float64 {
	t = t / d
	if t < .5 {
		return in(t)
	}
	return out(t - 1)
}
