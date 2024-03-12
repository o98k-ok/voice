package storage

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/o98k-ok/voice/internal/music"
)

type Storage interface {
	SaveMusic(music music.MusicKey) error
	HistoryMusics() chan music.MusicKey
	GetRootPath() string
	DelMusic(music music.MusicKey) error
}

type LocalFileStorage struct {
	Root string
}

func NewLocalFileStorage(root string) *LocalFileStorage {
	return &LocalFileStorage{
		Root: root,
	}
}

func (ls *LocalFileStorage) GetRootPath() string {
	return ls.Root
}

func (ls *LocalFileStorage) SaveMusic(music music.MusicKey) error {
	if len(ls.Root) == 0 {
		return nil
	}

	key, _ := json.Marshal(music)
	return os.WriteFile(fmt.Sprintf("%s/%x.raw", ls.Root, md5.Sum(key)), key, 0644)
}

func (ls *LocalFileStorage) DelMusic(music music.MusicKey) error {
	if len(ls.Root) == 0 {
		return nil
	}

	key, _ := json.Marshal(music)
	path := fmt.Sprintf("%s/%x.raw", ls.Root, md5.Sum(key))
	os.Remove(path)
	os.Remove(music.LocalPath)
	return nil
}

func (ls *LocalFileStorage) HistoryMusics() chan music.MusicKey {
	chans := make(chan music.MusicKey, 4)
	fdir, err := os.ReadDir(ls.Root)
	if err != nil {
		return chans
	}

	go func() {
		defer close(chans)
		for _, f := range fdir {
			if f.IsDir() {
				continue
			}

			if !strings.HasSuffix(f.Name(), ".raw") {
				continue
			}

			d, err := os.ReadFile(ls.Root + "/" + f.Name())
			if err != nil {
				continue
			}
			var key music.MusicKey
			if err := json.Unmarshal(d, &key); err != nil || len(key.BVID) == 0 {
				continue
			}

			chans <- key
		}
	}()
	return chans
}
