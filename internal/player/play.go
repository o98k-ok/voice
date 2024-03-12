package player

import (
	"time"

	"container/list"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/o98k-ok/voice/internal/music"
	"github.com/o98k-ok/voice/internal/pkg"
	"github.com/o98k-ok/voice/internal/storage"
)

type Player interface {
	SetMode(mode string) error
	GetMode() string
	AddInQueue(song *music.Music) error
	PlayIt(song *music.Music) error
	Play() error
	PlayList() ([]*music.Music, error)
	Exit() error
	Pause() error
	Next() (*music.Music, error)
	Prev() (*music.Music, error)
}

type VoicePlayer struct {
	VERSION      string
	Status       int64
	SampleRate   int64
	PlayList     *list.List
	PlayingQueue *StreamerQueue

	CurrentElem *list.Element
}

func NewVoicePlayer(sampleRate int64) *VoicePlayer {
	rate := beep.SampleRate(sampleRate)
	speaker.Init(rate, rate.N(time.Second/10))

	return &VoicePlayer{
		VERSION:      "1.0.0",
		Status:       0,
		SampleRate:   sampleRate,
		PlayList:     list.New(),
		PlayingQueue: NewStreamerQueue(int(sampleRate)),
	}
}

func (vp *VoicePlayer) InitPlayList(storage storage.Storage) {
	for m := range storage.HistoryMusics() {
		vp.AddInQueue(&music.Music{
			Name:      m.Name,
			Desc:      m.Desc,
			BvID:      m.BVID,
			LocalPath: m.LocalPath,
			Duration:  m.Duration,
		})
	}
}

func (vp *VoicePlayer) Pause() error {
	ctrl := vp.Current()
	if ctrl != nil && ctrl.PauseTrigger != nil {
		ctrl.PauseTrigger()
	}
	return nil
}

func (vp *VoicePlayer) FastForward() error {
	ctrl := vp.Current()
	if ctrl != nil && ctrl.SeekTrigger != nil && ctrl.PositionCallback != nil && ctrl.DurationCallback != nil {
		p := ctrl.PositionCallback() + ctrl.SampleRate.N(time.Second*5)
		if p > ctrl.DurationCallback() {
			p = ctrl.DurationCallback()
		}

		speaker.Lock()
		ctrl.SeekTrigger(p)
		speaker.Unlock()

		return nil
	}
	return nil
}

func (vp *VoicePlayer) FastBackward() error {
	ctrl := vp.Current()
	if ctrl != nil && ctrl.SeekTrigger != nil && ctrl.PositionCallback != nil {
		p := ctrl.PositionCallback() - ctrl.SampleRate.N(5*time.Second)
		if p < 0 {
			p = 0
		}
		return ctrl.SeekTrigger(p)
	}
	return nil
}

func (vp *VoicePlayer) Next() (*music.Music, error) {
	ctrl := vp.Current()
	if ctrl != nil && ctrl.NextTrigger != nil {
		ctrl.NextTrigger()
	}
	return nil, nil
}

func (vp *VoicePlayer) NextP(p *list.Element) (*music.Music, error) {
	ctrl := vp.Current()
	vp.CurrentElem = p
	if ctrl != nil && ctrl.NextTrigger != nil {
		ctrl.NextTrigger()
	}
	return nil, nil
}

func (vp *VoicePlayer) Current() *music.Music {
	if vp.CurrentElem == nil {
		return nil
	}
	return vp.CurrentElem.Value.(*music.Music)
}

func (vp *VoicePlayer) Run() error {
	speaker.Play(beep.Seq(vp.PlayingQueue, beep.Callback(func() {})))

	go func() {
		for {
			if vp.CurrentElem == nil {
				vp.CurrentElem = vp.PlayList.Front()
			}

			if vp.PlayingQueue.Size() == 0 && vp.CurrentElem != nil {
				// avoid stream conflict
				// time.Sleep(time.Millisecond * 100)
				vp.CurrentElem = pkg.NextForward(vp.PlayList, vp.CurrentElem)
				if vp.CurrentElem != nil {
					vp.PlayingQueue.Add(vp.CurrentElem)
				}
			}
			time.Sleep(time.Microsecond * 100)
		}
	}()
	return nil
}

func (vp *VoicePlayer) AddInQueue(song *music.Music) error {
	vp.PlayList.PushBack(song)
	return nil
}

func (vp *VoicePlayer) DryPlay(song *music.Music) error {
	if vp.CurrentElem == nil {
		vp.AddInQueue(song)
		return nil
	}
	vp.PlayList.InsertAfter(song, vp.CurrentElem)
	vp.Next()
	return nil
}
