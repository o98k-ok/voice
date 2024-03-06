package main

import (
	"fmt"
	"path"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/o98k-ok/voice/internal/music"
	"github.com/o98k-ok/voice/internal/player"
	"github.com/o98k-ok/voice/internal/ui"
)

const (
	ROOT = "./data"
)

func main() {
	player := player.NewVoicePlayer(48000)

	for _, v := range []string{"xc", "wwas"} {
		player.AddInQueue(&music.Music{
			Name:      v,
			Desc:      "开心快乐每一天",
			LocalPath: path.Join(ROOT, v+".wav"),
		})
	}
	player.Run()

	inputElem := ui.NewInputElem([]string{"ID", "标题", "时长", "BVID"}, []int{6, 30, 6, 11})
	framework := ui.NewFramework(
		[]ui.Element{
			ui.NewMenuElem([]string{"搜索", "歌单", "当前", "其他"}),
			inputElem,
		},
	)

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
