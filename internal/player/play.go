package player

import (
	"time"

	"container/list"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/o98k-ok/voice/internal/music"
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

func (vp *VoicePlayer) Pause() error {
	ctrl := vp.current()
	if ctrl != nil && ctrl.PauseTrigger != nil {
		ctrl.PauseTrigger()
	}
	return nil
}

func (vp *VoicePlayer) Next() (*music.Music, error) {
	ctrl := vp.current()
	if ctrl != nil && ctrl.NextTrigger != nil {
		ctrl.NextTrigger()
	}
	return nil, nil
}

func (vp *VoicePlayer) current() *music.Music {
	if vp.PlayingQueue.Current() == nil {
		return nil
	}
	return vp.PlayingQueue.Current().Value.(*music.Music)
}

func (vp *VoicePlayer) Info() *music.MusicRealtime {
	m := vp.current()
	return &music.MusicRealtime{
		Name:     m.Name,
		Desc:     m.Desc,
		URL:      m.URL,
		Duration: m.DurationCallback(),
		Position: m.PositionCallback(),
	}
}

func (vp *VoicePlayer) Run() error {
	speaker.Play(beep.Seq(vp.PlayingQueue, beep.Callback(func() {})))

	go func() {
		for {
			switch vp.PlayingQueue.Size() {
			case 0:
				elem := vp.PlayList.Front()
				if elem != nil {
					vp.PlayingQueue.Add(elem)
				}
			// try cache
			case 1:
				elem := vp.PlayingQueue.Current().Next()
				if elem == nil {
					elem = vp.PlayList.Front()
				}
				vp.PlayingQueue.Add(elem)
			default:
				time.Sleep(time.Millisecond * 100)
			}
		}
	}()

	// process := ui.NewMusicProcess()
	// go func() {
	// 	for {
	// 		if vp.current() == nil {
	// 			time.Sleep(time.Millisecond * 100)
	// 		}
	// 		process.Run(vp.current())
	// 	}
	// }()
	return nil
}

// shit lancet without insertAfter
// func (vp *VoicePlayer) insertSong(curret *datastructure.LinkNode[*music.Music], song *music.Music) *datastructure.LinkNode[*music.Music] {
// 	if vp.Queue.Head == nil || curret == nil {
// 		vp.Queue.InsertAtHead(song)
// 		return vp.Queue.Head
// 	}

// 	next := curret.Next
// 	node := datastructure.NewLinkNode(song)
// 	curret.Next = node
// 	node.Pre = curret
// 	node.Next = next

// 	if next != nil {
// 		next.Pre = node
// 	}
// 	return node
// }

// func (vp *VoicePlayer) PlayIt(song *music.Music) error {
// 	// process old song
// 	if !reflect.ValueOf(vp.NextCtrl).IsNil() {
// 		vp.NextCtrl.Send()
// 	}

// 	// process new song
// 	// vp.Playing = vp.insertSong(vp.Playing, song)
// 	return vp.playIt(song)
// }

func (vp *VoicePlayer) AddInQueue(song *music.Music) error {
	vp.PlayList.PushBack(song)
	return nil
}

func (vp *VoicePlayer) DryPlay(song *music.Music) error {
	p := vp.PlayingQueue.Current()
	if p == nil {
		vp.AddInQueue(song)
		return nil
	}

	vp.PlayList.InsertAfter(song, p.Next())

	vp.Next()
	// BUGS

	time.Sleep(time.Millisecond * 300)
	vp.Next()
	return nil
}

// func (vp *VoicePlayer) Play() error {
// 	if vp.Playing == nil {
// 		if vp.PlayList.Size() == 0 {
// 			// holding
// 			return nil
// 		}

// 		vp.Playing = vp.PlayList.Head
// 	}
// 	return vp.playIt(vp.Playing.Value)
// }
