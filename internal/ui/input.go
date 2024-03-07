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
	"github.com/o98k-ok/voice/internal/bilibili"
	"github.com/o98k-ok/voice/internal/convertor"
	"github.com/o98k-ok/voice/internal/music"
)

type InputElem struct {
	textInput  textinput.Model
	logo       *LogoElem
	result     *ListElem
	fetcher    bilibili.Fetcher
	mconvertor convertor.Convertor

	playChannel chan music.Music
	fetcherIdx  int
	active      bool
}

func NewInputElem(headers []string, widths []int) *InputElem {
	elem := textinput.New()
	elem.Focus()
	elem.Prompt = "> "

	return &InputElem{
		active:     false,
		textInput:  elem,
		logo:       &LogoElem{},
		result:     NewListElem(headers, widths, nil),
		fetcher:    bilibili.NewBlibliFetcher(netutil.NewHttpClient()),
		mconvertor: convertor.NewAfconvertConvertor("./data"),
		fetcherIdx: 1,
	}
}

func (ie *InputElem) Active() bool          { return ie.active }
func (ie *InputElem) SetActive(active bool) { ie.active = active }

func (ie *InputElem) RegisterPlayer() chan music.Music {
	ie.playChannel = make(chan music.Music, 10)
	return ie.playChannel
}

func (ie *InputElem) Init() tea.Cmd {
	return textinput.Blink
}

func (ie *InputElem) View() string {
	border := lipgloss.RoundedBorder()
	box := lipgloss.NewStyle().
		BorderStyle(border).
		Width(MaxWindowSize * 0.50).
		BorderForeground(highlightColor).
		Render(ie.textInput.View())

	left1 := "\n" + lipgloss.JoinHorizontal(lipgloss.Center, "🥳  ", box)
	left2 := ie.result.View()
	left := lipgloss.JoinVertical(lipgloss.Right, left1, "  ", left2)

	right1 := ie.logo.View()

	var right2 string
	if ie.result.table.SelectedRow() != nil {
		right2 = ie.result.table.SelectedRow()[1]
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
		pack = append(pack, []string{strconv.Itoa(i + (idx-1)*10), m.Name, m.Duration, m.URL})
	}
	return pack
}

func (ie *InputElem) MsgKeyBindings() map[string]map[string]func(v interface{}) tea.Cmd {
	return map[string]map[string]func(v interface{}) tea.Cmd{
		"tea.KeyMsg": {
			ALLMsgKey: func(v interface{}) tea.Cmd {
				var cmd tea.Cmd
				switch {
				case ie.textInput.Focused():
					ie.textInput, cmd = ie.textInput.Update(v)
					return cmd
				case ie.result.table.Focused():
					val, ok := v.(tea.KeyMsg)
					if !ok {
						return cmd
					}

					switch {
					case val.String() == "s":
						ie.result.table.Blur()
						ie.textInput.Focus()
					default:
						ie.result.table, cmd = ie.result.table.Update(v)
					}
					return cmd
				}
				return cmd
			},
			"right": func(v interface{}) tea.Cmd {
				switch {
				case ie.result.table.Focused():
					ie.fetcherIdx += 1
					pack := ie.fetch(ie.fetcherIdx)
					ie.result.ResetList(pack)
				}
				return nil
			},
			"left": func(v interface{}) tea.Cmd {
				switch {
				case ie.result.table.Focused():
					ie.fetcherIdx = mathutil.Max(1, ie.fetcherIdx-1)
					pack := ie.fetch(ie.fetcherIdx)
					ie.result.ResetList(pack)
				}
				return nil
			},
			"enter": func(v interface{}) tea.Cmd {
				switch {
				case ie.textInput.Focused():
					ie.fetcherIdx = 1
					pack := ie.fetch(ie.fetcherIdx)
					ie.result.ResetList(pack)
					ie.textInput.Blur()
					ie.result.table.Focus()
				case ie.result.table.Focused():
					msic := ie.result.table.SelectedRow()
					bvid := msic[3]
					for i, u := range ie.fetcher.GetAudioURL(bvid) {
						go func(bvID string, url string, idx int) {
							namein := fmt.Sprintf("%s/%s_%d.mp4", ROOT, bvID, idx)
							nameout := fmt.Sprintf("%s/%s_%d.wav", ROOT, bvID, idx)
							fin, _ := os.Create(namein)
							fout, _ := os.Create(nameout)
							ie.fetcher.Download(url, fin)
							fin.Close()

							fin, _ = os.Open(namein)

							ie.mconvertor.ConvertM4AToWav(fin, fout)
							fin.Close()
							fout.Close()
							os.Remove(namein)

							if ie.playChannel != nil {
								ie.playChannel <- music.Music{
									Name:      msic[0],
									URL:       url,
									LocalPath: nameout,
								}
							}
						}(bvid, u, i)
					}
				}
				return nil
			},
		},
	}
}
