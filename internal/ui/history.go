package ui

import (
	"container/list"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/o98k-ok/voice/internal/music"
	"github.com/o98k-ok/voice/internal/pkg"
	"github.com/o98k-ok/voice/internal/player"
)

type HistoryList struct {
	list    *ListElem
	player  *player.VoicePlayer
	active  bool
	current int
}

func NewHistoryList(headers []string, widths []int, player *player.VoicePlayer) *HistoryList {
	return &HistoryList{
		list:   NewListElem(headers, widths, nil),
		player: player,
	}
}

func (hl *HistoryList) Init() tea.Cmd {
	return nil
}

func (hl *HistoryList) View() string {
	return hl.list.View()
}

func (hl *HistoryList) fechList() [][]string {
	if hl.player.CurrentElem == nil {
		return nil
	}

	var values [][]string
	p := hl.player.CurrentElem
	hl.current = 0
	for i := 0; i < 3; i++ {
		if p.Prev() == nil {
			break
		}
		p = p.Prev()
		hl.current += 1
	}

	for i := 0; i < 10; i++ {
		if p == nil {
			break
		}

		m := p.Value.(*music.Music)
		if m != nil {
			values = append(values, []string{m.Name, m.Desc, m.Duration, m.BvID})
		}
		p = p.Next()
	}
	return values
}

func (hl *HistoryList) MsgKeyBindings() map[string]map[string]func(interface{}) tea.Cmd {
	return map[string]map[string]func(interface{}) tea.Cmd{
		"tea.KeyMsg": {
			ALLMsgKey: func(v interface{}) tea.Cmd {
				if !hl.active {
					return nil
				}
				hl.list.table.Focus()

				hl.list.ResetList(hl.fechList())
				var cmd tea.Cmd
				hl.list.table, cmd = hl.list.table.Update(v)
				return cmd
			},
			"tab": func(i interface{}) tea.Cmd {
				if !hl.active {
					return nil
				}

				hl.list.table.SetCursor(hl.current)
				return nil
			},
			"enter": func(i interface{}) tea.Cmd {
				if !hl.active {
					return nil
				}

				bvID := hl.list.table.SelectedRow()[3]
				var newElem *list.Element
				for p := hl.player.PlayList.Front(); p != nil; p = p.Next() {
					if p.Value.(*music.Music).BvID == bvID {
						newElem = p
						break
					}
				}
				p := pkg.NextBackward(hl.player.PlayList, newElem)
				hl.player.NextP(p)
				return nil
			},
		},
	}
}

func (hl *HistoryList) Active() bool {
	return hl.active
}

func (hl *HistoryList) SetActive(active bool) {
	hl.active = active
}
