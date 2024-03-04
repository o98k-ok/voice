package main

import (
	"fmt"
	"os"
	"path"

	"github.com/duke-git/lancet/v2/netutil"
	"github.com/o98k-ok/voice/internal/bilibili"
	"github.com/o98k-ok/voice/internal/convertor"
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

	// search page
	input := ui.NewInputModelProcess()
	keyword := input.Run()

	fetcher := bilibili.NewBlibliFetcher(netutil.NewHttpClient())
	musics, err := fetcher.Search(keyword, 1, 10)
	if err != nil {
		return
	}

	var pack [][]string
	for _, m := range musics {
		pack = append(pack, []string{m.Name, m.Desc, "3m30s", m.URL})
	}
	bvid := ui.NewTable(pack).Run()

	mconvertor := convertor.NewAfconvertConvertor("./data")

	for i, u := range fetcher.GetAudioURL(bvid) {
		go func(bvID string, url string, idx int) {
			namein := fmt.Sprintf("%s/%s_%d.mp4", ROOT, bvID, idx)
			nameout := fmt.Sprintf("%s/%s_%d.wav", ROOT, bvID, idx)
			fin, _ := os.Create(namein)
			fout, _ := os.Create(nameout)
			fetcher.Download(url, fin)
			fin.Close()

			fin, _ = os.Open(namein)

			mconvertor.ConvertM4AToWav(fin, fout)
			fin.Close()
			fout.Close()
			os.Remove(namein)

			err = player.AddInQueue(&music.Music{
				Name:      bvID,
				Desc:      bvID,
				LocalPath: path.Join(nameout),
			})
		}(bvid, u, i)
	}

	for {
		fmt.Print("Press [ENTER] to next. ")
		fmt.Scanln()
		player.Next()
	}
}
