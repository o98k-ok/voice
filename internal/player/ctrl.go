package player

import (
	"github.com/faiface/beep"
)

type Controller interface {
	Stream(samples [][2]float64) (n int, ok bool)
	Err() error
	Send()
}

type NextAble struct {
	streamer beep.Streamer
	ch       chan bool
}

func (nc *NextAble) Stream(samples [][2]float64) (n int, ok bool) {
	select {
	case <-nc.ch:
		return 0, false
	default:
		return nc.streamer.Stream(samples)
	}
}

func (nc *NextAble) Err() error {
	return nil
}

func (nc *NextAble) Send() {
	nc.ch <- true
}

func NewNextAbleStreamer(streamer beep.Streamer) *NextAble {
	return &NextAble{
		streamer: streamer,
		ch:       make(chan bool),
	}
}

type PauseCtrl struct {
	*beep.Ctrl
}

func (pc *PauseCtrl) Send() {
	pc.Ctrl.Paused = !pc.Ctrl.Paused
}

func NewPauseCtrl(streamer beep.Streamer, defa bool) *PauseCtrl {
	return &PauseCtrl{
		Ctrl: &beep.Ctrl{Streamer: streamer, Paused: defa},
	}
}
