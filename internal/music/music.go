package music

import (
	"fmt"
	"time"

	"github.com/faiface/beep"
)

type Music struct {
	Status int64

	Name       string
	Desc       string
	URL        string
	LocalPath  string
	SampleRate beep.SampleRate

	NextTrigger      func()
	PauseTrigger     func()
	DurationCallback func() int
	PositionCallback func() int
}

func (m Music) DurationRate() string {
	return fmt.Sprintf("%s/%s", m.SampleRate.D(m.PositionCallback()).Round(time.Second).String(),
		m.SampleRate.D(m.DurationCallback()).Round(time.Second).String())
}

type MusicRealtime struct {
	Name      string
	Desc      string
	URL       string
	LocalPath string
	Duration  int
	Position  int
}
