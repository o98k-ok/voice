package ui

import (
	"reflect"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Msg interface {
	String() string
}

type Element interface {
	Init() tea.Cmd
	View() string
	MsgKeyBindings() map[string]map[string]func(interface{}) tea.Cmd
}

type DefaultElem struct{}

func (de *DefaultElem) Init() tea.Cmd {
	return nil
}
func (de *DefaultElem) View() string {
	return ""
}
func (de *DefaultElem) MsgKeyBindings() map[string]map[string]func(interface{}) tea.Cmd {
	return map[string]map[string]func(interface{}) tea.Cmd{
		"tea.KeyMsg": {"ctrl+c": func(s interface{}) tea.Cmd { return tea.Quit }},
	}
}

// imp tea.Model
type Framework struct {
	Elems []Element
}

func NewFramework(elems []Element) *Framework {
	work := &Framework{
		Elems: []Element{&DefaultElem{}},
	}

	work.Elems = append(work.Elems, elems...)
	return work
}

func (f *Framework) Init() tea.Cmd {
	var cmds []tea.Cmd
	for _, e := range f.Elems {
		if cmd := e.Init(); cmd != nil {
			cmds = append(cmds, cmd)
		}
	}
	return tea.Batch(cmds...)
}

func (f *Framework) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m, ok := msg.(Msg)
	if !ok {
		// without String(), ignore anything
		return f, nil
	}
	return f.update(m)
}

func (f *Framework) update(msg Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	tpe := reflect.TypeOf(msg).String()
	for _, e := range f.Elems {
		bindings := e.MsgKeyBindings()[tpe]
		for key, fn := range bindings {
			if key == msg.String() || key == ALLMsgKey {
				if cmd := fn(msg); cmd != nil {
					cmds = append(cmds, cmd)
				}
			}
		}
	}
	return f, tea.Batch(cmds...)
}

// View mainly refactor here
func (f *Framework) View() string {
	var blocks []string
	for _, e := range f.Elems {
		blocks = append(blocks, e.View())
	}
	return lipgloss.JoinVertical(lipgloss.Left, blocks...)
}
