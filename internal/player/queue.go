package player

import (
	"container/list"
	"os"

	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
	"github.com/o98k-ok/voice/internal/music"
)

type StreamerQueue struct {
	streamers  []*Task[*list.Element]
	sampleRate int
}

func NewStreamerQueue(rate int) *StreamerQueue {
	return &StreamerQueue{
		streamers:  make([]*Task[*list.Element], 0, 2),
		sampleRate: rate,
	}
}

func (q *StreamerQueue) Add(elem *list.Element) {
	data, ok := elem.Value.(*music.Music)
	if !ok {
		return
	}

	f, err := os.Open(data.LocalPath)
	if err != nil {
		return
	}
	streamer, format, err := wav.Decode(f)
	if err != nil {
		return
	}

	s := beep.Resample(4, format.SampleRate, beep.SampleRate(q.sampleRate), streamer)
	c1 := NewPauseCtrl(s, false)
	c2 := NewNextAbleStreamer(c1)

	// setting music trigger
	data.NextTrigger = c2.Send
	data.PauseTrigger = c1.Send
	data.DurationCallback = streamer.Len
	data.PositionCallback = streamer.Position
	data.SampleRate = format.SampleRate

	q.streamers = append(q.streamers, NewTask(elem, c2))
}

func (q *StreamerQueue) Stream(samples [][2]float64) (n int, ok bool) {
	// We use the filled variable to track how many samples we've
	// successfully filled already. We loop until all samples are filled.
	filled := 0
	for filled < len(samples) {
		// There are no streamers in the queue, so we stream silence.
		if len(q.streamers) == 0 {
			for i := range samples[filled:] {
				samples[i][0] = 0
				samples[i][1] = 0
			}
			break
		}

		// We stream from the first streamer in the queue.
		n, ok := q.streamers[0].Stream(samples[filled:])
		// If it's drained, we pop it from the queue, thus continuing with
		// the next streamer.
		if !ok {
			q.streamers = q.streamers[1:]
		}
		// We update the number of filled samples.
		filled += n
	}

	return len(samples), true
}

func (q *StreamerQueue) Err() error {
	return nil
}

func (q *StreamerQueue) Next() {
	q.streamers[0].Send()
}

func (q *StreamerQueue) Size() int {
	return len(q.streamers)
}

func (q *StreamerQueue) Clear() {
	q.streamers = q.streamers[1:]
}

func (q *StreamerQueue) Current() *list.Element {
	if len(q.streamers) == 0 || q.streamers[0] == nil {
		return nil
	}
	return q.streamers[0].val
}
