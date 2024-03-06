package main

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder())

type model struct {
	table table.Model
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m model) View() string {
	return lipgloss.NewStyle().Border(lipgloss.DoubleBorder()).Render(m.table.View())
}

type TableModel struct {
	program *tea.Program
}

func NewTable(values [][]string) *TableModel {
	columns := []table.Column{
		{Title: "ID", Width: 20},
		// {Title: "标题", Width: 0},
		// {Title: "简介", Width: 0},
		// {Title: "时长", Width: 0},
		// {Title: "BVID", Width: 0},
	}

	var rows []table.Row
	for _, r := range values {
		rows = append(rows, table.Row{r[0]})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(1),
	)

	// s := table.DefaultStyles()
	// s.Header = s.Header.
	// 	BorderStyle(lipgloss.NormalBorder()).
	// 	BorderForeground(lipgloss.Color("240")).
	// 	BorderBottom(true).
	// 	Bold(false)
	// s.Selected = s.Selected.
	// 	Foreground(lipgloss.Color("229")).
	// 	Background(lipgloss.Color("57")).
	// 	Bold(false)
	// t.SetStyles(s)

	return &TableModel{
		program: tea.NewProgram(model{t}),
	}
}
func (t *TableModel) Run() string {
	res, err := t.program.Run()
	if err != nil {
		return ""
	}

	v, ok := res.(model)
	if !ok {
		return ""
	}
	// get id
	return v.table.SelectedRow()[0]
}

func main() {
	values := [][]string{
		{"1", "2", "3", "4", "5"},
	}
	NewTable(values).Run()
}
