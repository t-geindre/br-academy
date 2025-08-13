package game

import "time"

type Settings struct {
	WinSize   [2]int        `json:"w_size"`
	Threshold float64       `json:"a_threshold"`
	Release   time.Duration `json:"a_release"` // in nanoseconds
	Band      [2]float64    `json:"a_band"`    // in Hz
}
