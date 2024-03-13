package ui

import (
	"container/list"
	"fmt"
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/o98k-ok/voice/internal/music"
	"github.com/o98k-ok/voice/internal/pkg"
	"github.com/o98k-ok/voice/internal/player"
	"github.com/o98k-ok/voice/internal/storage"
)

type HistoryList struct {
	list    *ListElem
	player  *player.VoicePlayer
	storage storage.Storage
	active  bool
	current int

	page  int
	limit int
}

func NewHistoryList(player *player.VoicePlayer, index storage.Storage, headers []string, widths []int) *HistoryList {
	return &HistoryList{
		list:    NewListElem(headers, widths, nil),
		player:  player,
		page:    1,
		limit:   10,
		storage: index,
	}
}

func (hl *HistoryList) Init() tea.Cmd {
	return nil
}

func (hl *HistoryList) View() string {
	help := "enter play • ↓/j move down • ↑/k move up\nbackspace del music • ← page left • → page right"

	size := math.Ceil(float64(hl.player.PlayList.Len()) / float64(hl.limit))
	pageInfo := fmt.Sprintf("%d • %d/%d页", hl.list.table.Cursor()+1, hl.page, int(size))

	v := lipgloss.JoinVertical(lipgloss.Right, hl.list.View(), " ", pageInfo)
	return lipgloss.JoinVertical(lipgloss.Center, v, "\n", help)
}

func (hl *HistoryList) fechByBvID(bvID string, limit int) [][]string {
	var page int = 1
	var values [][]string

	p := hl.player.PlayList.Front()
	for {
		values = [][]string{}
		var got bool
		for i := 0; i < limit; i++ {
			if p == nil {
				break
			}
			m := p.Value.(*music.Music)
			if m.BvID == bvID {
				got = true
				hl.current = i
				hl.page = page
			}
			values = append(values, []string{m.Name, m.Desc, m.Duration, m.BvID, m.LocalPath})
			p = p.Next()
		}
		if got || p == nil {
			break
		}
		page += 1
	}
	return values
}

func (hl *HistoryList) fechList(off, limit int) [][]string {
	var page int = 1
	var values [][]string

	p := hl.player.PlayList.Front()
	for {
		values = [][]string{}
		for i := 0; i < limit; i++ {
			if p == nil {
				break
			}
			m := p.Value.(*music.Music)
			values = append(values, []string{m.Name, m.Desc, m.Duration, m.BvID, m.LocalPath})
			p = p.Next()
		}
		if page >= off || p == nil {
			break
		}
		page++
	}
	hl.current = len(values) / 2
	return values
}

func (hl *HistoryList) MsgKeyBindings() map[string]map[string]func(interface{}) tea.Cmd {
	return map[string]map[string]func(interface{}) tea.Cmd{
		"tea.KeyMsg": {
			"tab": func(i interface{}) tea.Cmd {
				if !hl.active {
					return nil
				}

				hl.list.table.Focus()
				if hl.player.CurrentElem != nil {
					hl.list.ResetList(hl.fechByBvID(hl.player.CurrentElem.Value.(*music.Music).BvID, 10))
				}
				hl.list.table.SetCursor(hl.current)
				return nil
			},
			"backspace": func(i interface{}) tea.Cmd {
				if !hl.active {
					return nil
				}

				m := hl.list.table.SelectedRow()
				hl.storage.DelMusic(music.MusicKey{
					Name:      m[0],
					Desc:      m[1],
					BVID:      m[3],
					Duration:  m[2],
					LocalPath: m[4],
				})

				p := hl.player.PlayList.Front()
				for ; p != nil; p = p.Next() {
					if p.Value.(*music.Music).BvID == m[3] {
						break
					}
				}
				hl.player.PlayList.Remove(p)

				// 刷新列表数据
				size := math.Ceil(float64(hl.player.PlayList.Len()) / float64(hl.limit))
				if hl.page > int(size) {
					hl.page -= 1
				}
				hl.list.ResetList(hl.fechList(hl.page, hl.limit))

				// 如果正在播放，就下一首
				if p == hl.player.CurrentElem {
					hl.player.Next()
				}
				return nil
			},
			"enter": func(i interface{}) tea.Cmd {
				if !hl.active || len(hl.list.table.Rows()) == 0 {
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
			" ": func(i interface{}) tea.Cmd {
				if !hl.active {
					return nil
				}
				hl.player.Pause()
				return nil
			},
			"right": func(i interface{}) tea.Cmd {
				if !hl.active {
					return nil
				}
				size := math.Ceil(float64(hl.player.PlayList.Len()) / float64(hl.limit))

				if hl.page < int(size) {
					hl.page += 1
				}

				hl.list.ResetList(hl.fechList(hl.page, hl.limit))
				hl.list.table.SetCursor(hl.current)
				var cmd tea.Cmd
				hl.list.table, cmd = hl.list.table.Update(i)
				return cmd
			},
			"left": func(i interface{}) tea.Cmd {
				if !hl.active {
					return nil
				}
				if hl.page > 1 {
					hl.page -= 1
				}
				hl.list.ResetList(hl.fechList(hl.page, hl.limit))
				hl.list.table.SetCursor(hl.current)
				var cmd tea.Cmd
				hl.list.table, cmd = hl.list.table.Update(i)
				return cmd
			},
			"up": func(i interface{}) tea.Cmd {
				var cmd tea.Cmd
				hl.list.table, cmd = hl.list.table.Update(i)
				return cmd
			},
			"down": func(i interface{}) tea.Cmd {
				var cmd tea.Cmd
				hl.list.table, cmd = hl.list.table.Update(i)
				return cmd
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
