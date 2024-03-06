package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LogoElem struct{}

func (le *LogoElem) Init() tea.Cmd { return nil }
func (le *LogoElem) View() string {

	v1 := `
__$$$$$___ __$$$$$___ __
__$$__$$__ __$$__$$__ __
__$$$$$___ __$$$$$___ __
__$$___$$_ __$$___$$_ __
__$$___$$_ __$$___$$_ __
__$$$$$$__ __$$$$$$__ __`
	v2 := `
__$$___$$_ ____$$$___ __$$$$_ ____$$$$_ __$$$$$$$_ __
__$$___$$_ ___$$_$$__ ___$$__ ___$$____ __$$______ __
___$$_$$__ __$$___$$_ ___$$__ __$$_____ __$$$$$___ __
___$$_$$__ __$$___$$_ ___$$__ __$$_____ __$$______ __
____$$$___ ___$$_$$__ ___$$__ ___$$____ __$$______ __
_____$____ ____$$$___ __$$$$_ ____$$$$_ __$$$$$$$_ __`
	return lipgloss.JoinVertical(lipgloss.Left, v1, v2)
}
func (le *LogoElem) MsgKeyBindings() map[string]map[string]func(interface{}) tea.Cmd { return nil }
