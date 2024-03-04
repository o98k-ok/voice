package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type (
	errMsg error
)

type InputModel struct {
	textInput textinput.Model
	err       error
}

type Input struct {
	process *tea.Program
	model   *InputModel
}

func NewInputModelProcess() *Input {
	model := initialModel()
	return &Input{
		model:   &model,
		process: tea.NewProgram(&model),
	}
}

func (i *Input) Run() string {
	result, _ := i.process.Run()
	v, ok := result.(InputModel)
	if !ok {
		return ""
	}
	return v.textInput.Value()
}

func initialModel() InputModel {
	ti := textinput.New()
	ti.Placeholder = "以父之名"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return InputModel{
		textInput: ti,
		err:       nil,
	}
}

func (m InputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m InputModel) View() string {
	return fmt.Sprintf(
		"请输入搜索内容\n\n%s\n\n%s",
		m.textInput.View(),
		"回车搜索",
	) + "\n"
}
