package bilibili

import (
	"net/http"
	"net/url"

	"github.com/duke-git/lancet/v2/netutil"
	"github.com/o98k-ok/voice/internal/pkg"
)

type DouyinHot struct {
	StatusCode int `json:"status_code"`
	MusicList  []struct {
		MusicInfo struct {
			ID     int64  `json:"id"`
			IDStr  string `json:"id_str"`
			Title  string `json:"title"`
			Author string `json:"author"`
			Album  string `json:"album"`
		} `json:"music_info"`
	} `json:"music_list"`
}

type DouyinHelper interface {
	HotKeys() []string
}

type Douyin struct {
	cli *netutil.HttpClient
}

func NewDouyinHelper(cli *netutil.HttpClient) DouyinHelper {
	return &Douyin{cli: cli}
}

func (dy *Douyin) HotKeys() []string {
	req := netutil.HttpRequest{
		RawURL:      "https://aweme.snssdk.com/aweme/v1/chart/music/list/?chart_id=6853972723954146568&count=20",
		Method:      http.MethodGet,
		Headers:     make(http.Header),
		QueryParams: make(url.Values),
	}

	hot, err := pkg.Request(dy.cli, &req, func(result *DouyinHot) bool { return result.StatusCode == 0 && len(result.MusicList) > 0 })
	if err != nil {
		return []string{}
	}

	var keys []string
	for _, music := range hot.MusicList {
		keys = append(keys, music.MusicInfo.Title+" "+music.MusicInfo.Author)
	}
	return keys
}
