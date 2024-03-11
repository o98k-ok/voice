package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/o98k-ok/voice/internal/player"
	"github.com/o98k-ok/voice/internal/storage"
	"github.com/o98k-ok/voice/internal/ui"
)

const (
	ROOT = "./data"
)

func main() {
	localIndex := storage.NewLocalFileStorage(ROOT)
	player := player.NewVoicePlayer(48000)
	player.InitPlayList(localIndex)
	player.Run()

	inputElem := ui.NewInputElem(player, localIndex,
		[]string{"ID", "标题", "时长", "BVID", "描述"},
		[]int{6, 30, 6, 12, 0},
	)
	processBar := ui.NewProcessLineElem(player)
	historyList := ui.NewHistoryList(player,
		[]string{"标题", "描述", "时长", "BVID"},
		[]int{40, 60, 12, 0},
	)

	elems := []ui.Element{processBar, inputElem, historyList}
	menu := ui.NewMenuElem([]string{"🤓  当前", "😂  搜索", "😳  列表"}, elems)
	framework := ui.NewFramework(menu, elems)

	program := tea.NewProgram(framework)
	fmt.Println(program.Run())
}
