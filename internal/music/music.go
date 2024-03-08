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
	BvID       string
	LocalPath  string
	SampleRate beep.SampleRate
	Duration   string

	NextTrigger      func()
	PauseTrigger     func()
	DurationCallback func() int
	PositionCallback func() int
}

func (m Music) DurationRate() string {
	if m.PositionCallback == nil || m.DurationCallback == nil {
		return ""
	}
	return fmt.Sprintf("%s/%s", m.SampleRate.D(m.PositionCallback()).Round(time.Second).String(),
		m.SampleRate.D(m.DurationCallback()).Round(time.Second).String())
}

type MusicKey struct {
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	BVID      string `json:"bvid"`
	LocalPath string `json:"local_path"`
	Duration  string `json:"duration"`
}
