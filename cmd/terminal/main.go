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

	inputElem := ui.NewInputElem([]string{"ID", "标题", "时长", "BVID", "描述"}, []int{6, 30, 6, 12, 0}, localIndex)
	processBar := ui.NewProcessLineElem(player)
	historyList := ui.NewHistoryList([]string{"标题", "描述", "时长", "BVID"}, []int{40, 60, 12, 0}, player)

	elems := []ui.Element{inputElem, historyList, processBar}
	menu := ui.NewMenuElem([]string{"搜索", "列表", "当前"}, elems)
	framework := ui.NewFramework(menu, elems)
	go func() {
		channel := inputElem.RegisterPlayer()
		for {
			msic := <-channel
			player.DryPlay(&msic)
		}
	}()
	program := tea.NewProgram(framework)
	fmt.Println(program.Run())
}
