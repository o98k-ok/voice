package main

import (
	"path"

	"github.com/o98k-ok/voice/internal/music"
	"github.com/o98k-ok/voice/internal/player"
	"github.com/o98k-ok/voice/internal/ui"
)

const (
	ROOT = "./data"
)

func main() {
	player := player.NewVoicePlayer(48000)

	for _, v := range []string{"pd", "wmdr", "xc", "wwas"} {
		player.AddInQueue(&music.Music{
			Name:      v,
			Desc:      "开心快乐每一天",
			LocalPath: path.Join(ROOT, v+".wav"),
		})
	}
	player.Run()

	<-ui.GlobalExit
}
