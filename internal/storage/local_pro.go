package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/huntsman-li/go-cache"

	"github.com/o98k-ok/voice/internal/music"
)

type LocalFileStorageV2 struct {
	Cache    cache.Cache
	MetaPath string
	RawPath  string
}

func NewLocalFileStorageV2(root string) *LocalFileStorageV2 {
	if _, err := os.Stat(root); err != nil {
		panic(err)
	}

	var storage = &LocalFileStorageV2{
		MetaPath: root + "/meta",
		RawPath:  root + "/raw",
	}

	if _, err := os.Stat(storage.RawPath); err != nil {
		os.Mkdir(storage.RawPath, os.ModePerm)
	}
	if _, err := os.Stat(storage.MetaPath); err != nil {
		os.Mkdir(storage.MetaPath, os.ModePerm)
	}

	var err error
	options := cache.Options{
		Adapter:       "file",
		AdapterConfig: storage.MetaPath,
		Interval:      2,
	}

	if storage.Cache, err = cache.Cacher(options); err != nil {
		panic(err)
	}
	return storage
}

func (ls *LocalFileStorageV2) GetRootPath() string {
	return ls.RawPath
}

func (ls *LocalFileStorageV2) SaveMusic(music music.MusicKey) error {
	// if len(ls.Root) == 0 {
	// 	return nil
	// }

	var err error
	music.Key = uuid.NewString()
	func() {
		keys, _ := ls.Cache.Get(ls.allKeys()).([]string)
		keys = append(keys, music.Key)
		err = ls.Cache.Put(ls.allKeys(), keys, -1)
	}()
	if err != nil {
		return err
	}

	func() {
		data, _ := json.Marshal(music)
		err = ls.Cache.Put(ls.musicKey(music.Key), data, -1)
	}()
	return err
}

func (ls *LocalFileStorageV2) musicKey(key string) string {
	return fmt.Sprintf("raw_music:%s", key)
}

func (ls *LocalFileStorageV2) allKeys() string {
	return "VOICE_KEYS"
}

func (ls *LocalFileStorageV2) DelMusic(music music.MusicKey) error {
	// if len(ls.Root) == 0 {
	// 	return nil
	// }

	keys, _ := ls.Cache.Get(ls.allKeys()).([]string)
	for i, k := range keys {
		if k == music.Key {
			keys = append(keys[:i], keys[i+1:]...)
			break
		}
	}
	if err := ls.Cache.Put(ls.allKeys(), keys, -1); err != nil {
		return err
	}

	if err := ls.Cache.Delete(ls.musicKey(music.Key)); err != nil {
		return err
	}
	return nil
}

func (ls *LocalFileStorageV2) HistoryMusics() chan music.MusicKey {
	chans := make(chan music.MusicKey, 4)
	keys, _ := ls.Cache.Get(ls.allKeys()).([]string)

	go func() {
		defer close(chans)
		for _, key := range keys {
			musicData := ls.Cache.Get(ls.musicKey(key))
			if musicData == nil {
				continue // 如果音乐数据不存在，跳过
			}

			var music music.MusicKey
			if err := json.Unmarshal(musicData.([]byte), &music); err != nil {
				continue // 如果解析失败，跳过
			}

			chans <- music
		}
	}()

	return chans
}
