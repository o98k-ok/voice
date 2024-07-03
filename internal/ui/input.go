package ui

import (
	"fmt"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/duke-git/lancet/v2/mathutil"
	"github.com/duke-git/lancet/v2/netutil"
	"github.com/duke-git/lancet/v2/random"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/o98k-ok/voice/internal/bilibili"
	"github.com/o98k-ok/voice/internal/convertor"
	"github.com/o98k-ok/voice/internal/music"
	"github.com/o98k-ok/voice/internal/pkg"
	"github.com/o98k-ok/voice/internal/player"
	"github.com/o98k-ok/voice/internal/storage"
)

type InputElem struct {
	textInput  textinput.Model
	logo       *LogoElem
	result     *ListElem
	fetcher    bilibili.Fetcher
	dyHelper   bilibili.DouyinHelper
	mconvertor convertor.Convertor

	fetcherIdx  int
	active      bool
	storage     storage.Storage
	player      *player.VoicePlayer
	suggestKeys []string
	suggestKey  string
}

func NewInputElem(player *player.VoicePlayer, storage storage.Storage, headers []string, widths []int) *InputElem {
	elem := textinput.New()
	elem.Focus()
	elem.Prompt = "> "

	return &InputElem{
		active:     false,
		textInput:  elem,
		logo:       &LogoElem{},
		result:     NewListElem(headers, widths, nil),
		fetcher:    bilibili.NewBlibliFetcher(netutil.NewHttpClient()),
		dyHelper:   bilibili.NewDouyinHelper(netutil.NewHttpClient()),
		mconvertor: convertor.NewAfconvertConvertor(storage.GetRootPath()),
		fetcherIdx: 1,
		storage:    storage,
		player:     player,
	}
}

func (ie *InputElem) Active() bool          { return ie.active }
func (ie *InputElem) SetActive(active bool) { ie.active = active }

func (ie *InputElem) Init() tea.Cmd {
	keys := ie.dyHelper.HotKeys()
	ie.suggestKeys = keys
	ie.suggestKey = keys[0]
	return textinput.Blink
}

func (ie *InputElem) View() string {
	ie.textInput.Placeholder = ie.suggestKey
	border := lipgloss.RoundedBorder()
	box := lipgloss.NewStyle().
		BorderStyle(border).
		Width(MaxWindowSize * 0.50).
		BorderForeground(highlightColor).
		Render(ie.textInput.View())

	left1 := "\n" + lipgloss.JoinHorizontal(lipgloss.Center, "🥳  ", box)
	left2 := ie.result.View()
	help := "> switch mode • ↓ move down • ↑ move up\n← page left • → page right • tab next menu"
	left3 := lipgloss.NewStyle().Bold(true).Render(help)
	left := lipgloss.JoinVertical(lipgloss.Right, left1, "  ", left2, "  ", left3)

	right1 := ie.logo.View()

	var right2 string
	if ie.result.table.SelectedRow() != nil {
		right2 = pkg.RenderWithWidth(ie.result.table.SelectedRow()[1], MaxWindowSize*0.4)
	}
	right := lipgloss.JoinVertical(lipgloss.Right, right1, "  ", right2, "🛵🛵🛵🛵 ")

	return lipgloss.JoinHorizontal(lipgloss.Top, left, "     ", right)
}

func (ie *InputElem) fetch(idx int) [][]string {
	musics, err := ie.fetcher.Search(ie.textInput.Value(), idx, 10)
	if err != nil {
		return nil
	}

	var pack [][]string
	for i, m := range musics {
		pack = append(pack, []string{strconv.Itoa(i + (idx-1)*10), m.Name, m.Duration, m.BvID, m.Desc})
	}
	return pack
}

func (ie *InputElem) MsgKeyBindings() map[string]map[string]func(v interface{}) tea.Cmd {
	return map[string]map[string]func(v interface{}) tea.Cmd{
		"tea.KeyMsg": {
			ALLMsgKey: func(v interface{}) tea.Cmd {
				if !ie.active || strutil.ContainsAny(v.(tea.KeyMsg).String(), []string{"up", "down", ">", " "}) {
					return nil
				}
				var cmd tea.Cmd
				if ie.textInput.Focused() {
					ie.textInput, cmd = ie.textInput.Update(v)
				}
				return cmd
			},
			"down": func(v interface{}) tea.Cmd {
				if !ie.active {
					return nil
				}
				if ie.textInput.Focused() {
					ie.textInput.Blur()
					ie.result.table.Focus()
					return nil
				}
				var cmd tea.Cmd
				ie.result.table, cmd = ie.result.table.Update(v)
				return cmd
			},
			"up": func(v interface{}) tea.Cmd {
				if !ie.active {
					return nil
				}
				var cmd tea.Cmd
				if ie.result.table.Cursor() <= 0 {
					ie.result.table.Blur()
					ie.textInput.Focus()
					return cmd
				}
				ie.result.table, cmd = ie.result.table.Update(v)
				return cmd
			},
			" ": func(v interface{}) tea.Cmd {
				if !ie.active {
					return nil
				}
				var cmd tea.Cmd
				switch {
				case ie.textInput.Focused():
					ie.textInput, cmd = ie.textInput.Update(v)
					return cmd
				case ie.result.table.Focused():
					ie.player.Pause()
					return cmd
				}
				return cmd
			},
			">": func(v interface{}) tea.Cmd {
				if !ie.active {
					return nil
				}
				switch {
				case ie.textInput.Focused():
					ie.textInput.Blur()
					ie.result.table.Focus()
				case ie.result.table.Focused():
					ie.result.table.Blur()
					ie.textInput.Focus()
				default:
				}
				return nil
			},
			"right": func(v interface{}) tea.Cmd {
				if !ie.active {
					return nil
				}
				switch {
				case ie.result.table.Focused():
					ie.fetcherIdx += 1
					pack := ie.fetch(ie.fetcherIdx)
					ie.result.ResetList(pack)
				case ie.textInput.Focused():
					ie.textInput.SetValue(ie.suggestKey)
					ie.suggestKey = ie.suggestKeys[random.RandInt(0, len(ie.suggestKeys))]
				}
				return nil
			},
			"left": func(v interface{}) tea.Cmd {
				if !ie.active {
					return nil
				}
				switch {
				case ie.result.table.Focused():
					ie.fetcherIdx = mathutil.Max(1, ie.fetcherIdx-1)
					pack := ie.fetch(ie.fetcherIdx)
					ie.result.ResetList(pack)
				}
				return nil
			},
			"enter": func(v interface{}) tea.Cmd {
				if !ie.active {
					return nil
				}
				switch {
				case ie.textInput.Focused():
					if len(ie.textInput.Value()) == 0 {
						return nil
					}
					ie.fetcherIdx = 1
					pack := ie.fetch(ie.fetcherIdx)
					ie.result.ResetList(pack)
					ie.textInput.Blur()
					ie.result.table.Focus()
				case ie.result.table.Focused():
					if len(ie.result.table.Rows()) == 0 {
						return nil
					}
					msic := ie.result.table.SelectedRow()
					bvid := msic[3]
					for i, u := range ie.fetcher.GetAudioURL(bvid) {
						go func(bvID string, url string, idx int) {
							root := ie.storage.GetRootPath()
							namein := fmt.Sprintf("%s/%s_%d.mp4", root, bvID, idx)
							nameout := fmt.Sprintf("%s/%s_%d.wav", root, bvID, idx)

							fin, _ := os.Create(namein)
							fout, _ := os.Create(nameout)
							ie.fetcher.Download(url, fin)
							fin.Close()

							fin, _ = os.Open(namein)

							ie.mconvertor.ConvertM4AToWav(fin, fout)
							fin.Close()
							fout.Close()
							os.Remove(namein)

							ie.storage.SaveMusic(music.MusicKey{
								Name:      msic[1],
								Desc:      msic[4],
								BVID:      bvID,
								LocalPath: nameout,
								Duration:  msic[2],
							})

							mm := music.Music{
								Name:      msic[1],
								Desc:      msic[4],
								URL:       url,
								BvID:      bvID,
								LocalPath: nameout,
								Duration:  msic[2],
							}
							ie.player.DryPlay(&mm)
						}(bvid, u, i)
					}
				}
				return nil
			},
		},
	}
}
