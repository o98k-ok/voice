package pkg

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/duke-git/lancet/v2/mathutil"
)

func RenderWithWidth(value string, width int) string {
	width = mathutil.Min(len(value), width)
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(value)
}
