package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuElem struct {
	Tabs         []string
	ActiveTabIdx int
	menus        []Element
}

func NewMenuElem(tabs []string, menus []Element) *MenuElem {
	menus[0].SetActive(true)
	return &MenuElem{
		Tabs:         tabs,
		ActiveTabIdx: 0,
		menus:        menus,
	}
}

func (m *MenuElem) Active() bool          { return true }
func (m *MenuElem) SetActive(active bool) {}

func (m *MenuElem) Init() tea.Cmd { return nil }
func (m *MenuElem) View() string {
	if len(m.Tabs) == 0 {
		return ""
	}
	everySize := MaxWindowSize / len(m.Tabs)
	style := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(highlightColor).Width(everySize).Padding(0, 5, 0, 5)

	res := make([]string, 0, len(m.Tabs))
	for i, r := range m.Tabs {
		if i == m.ActiveTabIdx {
			activeStyle := style.Copy().Foreground(lipgloss.Color("#66CC00")).
				Border(lipgloss.RoundedBorder(), true, true, false, true)
			res = append(res, activeStyle.Render(r))
			continue
		}
		res = append(res, style.Render(r))
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, res...)
}
func (m *MenuElem) MsgKeyBindings() map[string]map[string]func(v interface{}) tea.Cmd {
	return map[string]map[string]func(v interface{}) tea.Cmd{
		"tea.KeyMsg": {
			"tab": func(v interface{}) tea.Cmd {
				m.menus[m.ActiveTabIdx].SetActive(false)
				m.ActiveTabIdx++
				m.ActiveTabIdx = m.ActiveTabIdx % len(m.Tabs)
				m.menus[m.ActiveTabIdx].SetActive(true)
				// fmt.Print("\033[H\033[2J")
				return nil
			},
			"shift+tab": func(v interface{}) tea.Cmd {
				m.menus[m.ActiveTabIdx].SetActive(false)
				m.ActiveTabIdx += 3
				m.ActiveTabIdx = m.ActiveTabIdx % len(m.Tabs)
				m.menus[m.ActiveTabIdx].SetActive(true)
				// fmt.Print("\033[H\033[2J")
				return nil
			}},
	}
}
