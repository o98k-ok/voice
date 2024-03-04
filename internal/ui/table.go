package ui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type tableModel struct {
	table table.Model
}

type TableModel struct {
	program *tea.Program
}

func NewTable(values [][]string) *TableModel {
	columns := []table.Column{
		{Title: "标题", Width: 40},
		{Title: "简介", Width: 100},
		{Title: "时长", Width: 100},
		{Title: "BVID", Width: 10},
	}

	var rows []table.Row
	for _, r := range values {
		rows = append(rows, table.Row{r[0], r[1], r[2], r[3]})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return &TableModel{
		program: tea.NewProgram(tableModel{t}),
	}
}

func (t *TableModel) Run() string {
	res, err := t.program.Run()
	if err != nil {
		return ""
	}

	v, ok := res.(tableModel)
	if !ok {
		return ""
	}
	// get bvid
	return v.table.SelectedRow()[3]
}

func (m tableModel) Init() tea.Cmd { return nil }

func (m tableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.table.Focused() {
				m.table.Blur()
			} else {
				m.table.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			return m, tea.Quit
		}
	}
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m tableModel) View() string {
	return baseStyle.Render(m.table.View()) + "\n"
}
