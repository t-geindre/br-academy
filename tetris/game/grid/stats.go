package grid

import "github.com/hajimehoshi/ebiten/v2"

type Stats struct {
	Level    int
	Score    int
	Lines    int
	TickRate int
}

func NewStats() *Stats {
	return &Stats{}
}

func (s *Stats) Reset() {
	s.Level = 0
	s.Score = 0
	s.Lines = 0
	s.TickRate = int(48.0 / 60.0 * float64(ebiten.TPS()))
}

func (s *Stats) AddLines(lines int) {
	if lines < 1 {
		return
	}

	s.Lines += lines
	s.Score += []int{0, 100, 300, 500, 800}[lines] * (s.Level + 1)
	s.Level = s.Lines / 10

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
			s.TickRate = int(r.m * float64(ebiten.TPS()))
			return
		}
	}

	s.TickRate = int(1.0 / 60.0 * float64(ebiten.TPS()))
}

func (s *Stats) GetTickRate() int {
	return s.TickRate
}
