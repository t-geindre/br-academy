package grid

import "github.com/hajimehoshi/ebiten/v2"

type Stats struct {
	Level int
	Score int
	Lines int
}

func NewStats() *Stats {
	return &Stats{}
}

func (s *Stats) Reset() {
	s.Level = 30
	s.Score = 0
	s.Lines = 0
}

func (s *Stats) AddLines(lines int) {
	if lines < 1 {
		return
	}

	s.Lines += lines
	s.Level = s.Lines / 10
	s.Score += []int{0, 100, 300, 500, 800}[lines] * (s.Level + 1)
}

func (s *Stats) GetTickRate() int {
	for _, r := range []struct {
		l int
		m float64
	}{
		{0, 48.0 / 60.0},
		{1, 43.0 / 60.0},
		{10, 28.0 / 60.0},
		{20, 18.0 / 60.0},
		{30, 8.0 / 60.0},
	} {
		if s.Level <= r.l {
			return int(r.m * float64(ebiten.TPS()))
		}
	}

	return int(1.0 / 60.0 * float64(ebiten.TPS()))
}
