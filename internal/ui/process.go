package ui

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/o98k-ok/voice/internal/pkg"
	"github.com/o98k-ok/voice/internal/player"
)

type tickMsg time.Time

func (t tickMsg) String() string {
	return "tickMsg"
}

type ProceeLineElem struct {
	progress progress.Model
	player   *player.VoicePlayer
	active   bool
}

func NewProcessLineElem(player *player.VoicePlayer) *ProceeLineElem {
	return &ProceeLineElem{
		progress: progress.New(progress.WithDefaultGradient(), progress.WithoutPercentage()),
		player:   player,
	}
}

func (pe *ProceeLineElem) Active() bool          { return pe.active }
func (pe *ProceeLineElem) SetActive(active bool) { pe.active = active }

func (pe *ProceeLineElem) Init() tea.Cmd { return tickCmd() }

func (pe *ProceeLineElem) View() string {
	if pe.player == nil || pe.player.Current() == nil {
		return ""
	}

	music := pe.player.Current()
	return lipgloss.JoinVertical(lipgloss.Center,
		pkg.RenderWithWidth(music.Name, MaxWindowSize*0.6), "\n",
		pkg.RenderWithWidth(music.Desc, MaxWindowSize*0.6), "\n",
		pe.progress.ViewAs(pe.progress.Percent())+"  "+music.DurationRate(), "\n",
		lipgloss.NewStyle().Bold(true).Render("tab next menu • p prev • space pause/play • n next"))
}

func (pe *ProceeLineElem) MsgKeyBindings() map[string]map[string]func(interface{}) tea.Cmd {
	return map[string]map[string]func(interface{}) tea.Cmd{
		"ui.tickMsg": {
			"tickMsg": func(i interface{}) tea.Cmd {
				if pe.progress.Percent() == 1.0 || pe.player == nil || pe.player.Current() == nil {
					return tickCmd()
				}

				if pe.player.Current().PositionCallback == nil || pe.player.Current().DurationCallback == nil {
					return tickCmd()
				}
				cmd := pe.progress.SetPercent(float64(pe.player.Current().PositionCallback()) / float64(pe.player.Current().DurationCallback()))
				return tea.Batch(tickCmd(), cmd)
			},
		},
		"tea.KeyMsg": {
			" ": func(v interface{}) tea.Cmd {
				if pe.active {
					pe.player.Pause()
				}
				return nil
			},
			"n": func(v interface{}) tea.Cmd {
				if pe.active {
					pe.player.Next()
				}
				return nil
			},
			"p": func(interface{}) tea.Cmd {
				if pe.active {
					p := pkg.NextN(pe.player.PlayList, pe.player.CurrentElem, -2)
					pe.player.NextP(p)
				}
				return nil
			},
		},
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
