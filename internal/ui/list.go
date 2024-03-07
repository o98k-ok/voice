package ui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ListElem struct {
	table  table.Model
	active bool
}

func NewListElem(headers []string, widths []int, values [][]string) *ListElem {
	var columns []table.Column
	for i := 0; i < len(headers); i++ {
		columns = append(columns, table.Column{
			Title: headers[i],
			Width: widths[i],
		})
	}

	var rows []table.Row
	for _, r := range values {
		rows = append(rows, table.Row(r))
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(highlightColor).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)
	return &ListElem{
		table: t,
	}
}

func (le *ListElem) Active() bool          { return le.active }
func (le *ListElem) SetActive(active bool) { le.active = active }

func (le *ListElem) Init() tea.Cmd { return nil }
func (le *ListElem) View() string {
	return le.table.View()
}
func (le *ListElem) MsgKeyBindings() map[string]map[string]func(interface{}) tea.Cmd { return nil }

func (le *ListElem) ResetList(values [][]string) {
	var rows []table.Row
	for _, r := range values {
		rows = append(rows, table.Row(r))
	}
	le.table.SetRows(rows)
}
