package ui

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/o98k-ok/voice/internal/music"
)

const (
	padding  = 2
	maxWidth = 14
)

var (
	GlobalExit = make(chan interface{})
)

type tickMsg time.Time

type MusicProcess struct {
	progress progress.Model
	music    *music.Music
}

func NewMusicProcess() MusicProcess {
	return MusicProcess{
		progress: progress.New(progress.WithDefaultGradient(), progress.WithoutPercentage()),
	}
}

func (m MusicProcess) Run(music *music.Music) {
	m.music = music
	tea.NewProgram(m).Run()
}

func (m MusicProcess) Init() tea.Cmd {
	return tickCmd()
}

func (m MusicProcess) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == tea.KeyCtrlC.String() {
			GlobalExit <- struct{}{}
			return m, tea.Quit
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}

		cmd := m.progress.SetPercent(float64(m.music.PositionCallback()) / float64(m.music.DurationCallback()))
		return m, tea.Batch(tickCmd(), cmd)

	// FrameMsg is sent when the progress bar wants to animate itself
	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (m MusicProcess) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + m.music.Name + pad + m.music.DurationRate() + "\n" +
		pad + m.music.Desc + "\n" +
		pad + m.progress.View() + "\n\n"
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
