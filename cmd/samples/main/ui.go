package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/o98k-ok/voice/internal/ui"
)

func main() {
	program := tea.NewProgram(ui.NewFramework(
		[]ui.Element{
			ui.NewMenuElem([]string{"搜索", "歌单", "当前", "其他"}),
			ui.NewInputElem([]string{"ID", "标题", "时长", "BVID"}, []int{6, 30, 6, 11}),
		},
	))

	fmt.Println(program.Run())
}
