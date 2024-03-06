package bilibili

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/netutil"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/o98k-ok/voice/internal/music"
	"github.com/o98k-ok/voice/internal/pkg"
)

type Fetcher interface {
	// https://api.bilibili.com/x/web-interface/search/type?__refresh__=true&_extra=&context=&page=1&page_size=42&platform=pc&highlight=1&single_column=0&keyword=%E5%91%A8%E6%9D%B0%E4%BC%A6&category_id=&search_type=video&dynamic_offset=0&preload=true&com2co=true
	Search(keyword string, page, pageSize int) ([]*music.Music, error)
	GetAudioURL(bvid string) []string
	Download(url string, writer io.Writer) error
}

type BilibiliFetcher struct {
	cli *netutil.HttpClient
}

func NewBlibliFetcher(cli *netutil.HttpClient) *BilibiliFetcher {
	return &BilibiliFetcher{
		cli: cli,
	}
}

type BiliResult struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Result []struct {
			Type             string        `json:"type"`
			ID               int           `json:"id"`
			Author           string        `json:"author"`
			Mid              int           `json:"mid"`
			Typeid           string        `json:"typeid"`
			Typename         string        `json:"typename"`
			Arcurl           string        `json:"arcurl"`
			Aid              int           `json:"aid"`
			Bvid             string        `json:"bvid"`
			Title            string        `json:"title"`
			Description      string        `json:"description"`
			Arcrank          string        `json:"arcrank"`
			Pic              string        `json:"pic"`
			Play             int           `json:"play"`
			VideoReview      int           `json:"video_review"`
			Favorites        int           `json:"favorites"`
			Tag              string        `json:"tag"`
			Review           int           `json:"review"`
			Pubdate          int           `json:"pubdate"`
			Senddate         int           `json:"senddate"`
			Duration         string        `json:"duration"`
			Badgepay         bool          `json:"badgepay"`
			HitColumns       []string      `json:"hit_columns"`
			ViewType         string        `json:"view_type"`
			IsPay            int           `json:"is_pay"`
			IsUnionVideo     int           `json:"is_union_video"`
			RecTags          interface{}   `json:"rec_tags"`
			NewRecTags       []interface{} `json:"new_rec_tags"`
			RankScore        int           `json:"rank_score"`
			Like             int           `json:"like"`
			Upic             string        `json:"upic"`
			Corner           string        `json:"corner"`
			Cover            string        `json:"cover"`
			Desc             string        `json:"desc"`
			URL              string        `json:"url"`
			RecReason        string        `json:"rec_reason"`
			Danmaku          int           `json:"danmaku"`
			BizData          interface{}   `json:"biz_data"`
			IsChargeVideo    int           `json:"is_charge_video"`
			Vt               int           `json:"vt"`
			EnableVt         int           `json:"enable_vt"`
			VtDisplay        string        `json:"vt_display"`
			Subtitle         string        `json:"subtitle"`
			EpisodeCountText string        `json:"episode_count_text"`
			ReleaseStatus    int           `json:"release_status"`
			IsIntervene      int           `json:"is_intervene"`
		} `json:"result"`
	} `json:"data"`
}

func (bf *BilibiliFetcher) Search(keyword string, page, pageSize int) ([]*music.Music, error) {
	req := netutil.HttpRequest{
		RawURL:      "https://api.bilibili.com/x/web-interface/search/type",
		Method:      http.MethodGet,
		Headers:     make(http.Header),
		QueryParams: make(url.Values),
	}
	bf.fillHeader(&req)

	req.QueryParams.Add("__refresh__", "true")
	req.QueryParams.Add("platform", "pc")
	req.QueryParams.Add("highlight", "1")
	req.QueryParams.Add("signgle_column", "0")
	req.QueryParams.Add("search_type", "video")
	req.QueryParams.Add("dynamic_offset", "0")
	req.QueryParams.Add("preload", "true")
	req.QueryParams.Add("com2co", "true")

	req.QueryParams.Add("page", convertor.ToString(page))
	req.QueryParams.Add("page_size", convertor.ToString(pageSize))
	req.QueryParams.Add("keyword", convertor.ToString(keyword))

	result, err := pkg.Request(bf.cli, &req, func(result *BiliResult) bool { return result.Code == 0 || len(result.Data.Result) != 0 })
	if err != nil {
		return nil, err
	}

	musics := make([]*music.Music, 0, len(result.Data.Result))
	for _, item := range result.Data.Result {
		musics = append(musics, &music.Music{
			Name: func() string {
				extra := fmt.Sprintf("<em class=\"keyword\">%s</em>", keyword)
				str := strings.ReplaceAll(item.Title, extra, keyword)
				return strutil.RemoveNonPrintable(str)
			}(),
			Desc:     strutil.RemoveNonPrintable(item.Description),
			URL:      item.Bvid,
			Duration: strings.ReplaceAll(item.Duration, ":", "m") + "s",
		})
	}
	return musics, nil
}

type BiliDetail struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Bvid string `json:"bvid"`
		Cid  int    `json:"cid"`
	} `json:"data"`
}

type BiliPlayURL struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		Dash struct {
			Duration       int     `json:"duration"`
			MinBufferTime  float64 `json:"minBufferTime"`
			MinBufferTime0 float64 `json:"min_buffer_time"`
			Audio          []struct {
				ID            int      `json:"id"`
				BaseURL       string   `json:"baseUrl"`
				BaseURL0      string   `json:"base_url"`
				BackupURL     []string `json:"backupUrl"`
				BackupURL0    []string `json:"backup_url"`
				Bandwidth     int      `json:"bandwidth"`
				MimeType      string   `json:"mimeType"`
				MimeType0     string   `json:"mime_type"`
				Codecs        string   `json:"codecs"`
				Width         int      `json:"width"`
				Height        int      `json:"height"`
				FrameRate     string   `json:"frameRate"`
				FrameRate0    string   `json:"frame_rate"`
				Sar           string   `json:"sar"`
				StartWithSap  int      `json:"startWithSap"`
				StartWithSap0 int      `json:"start_with_sap"`
				SegmentBase   struct {
					Initialization string `json:"Initialization"`
					IndexRange     string `json:"indexRange"`
				} `json:"SegmentBase"`
				SegmentBase0 struct {
					Initialization string `json:"initialization"`
					IndexRange     string `json:"index_range"`
				} `json:"segment_base"`
				Codecid int `json:"codecid"`
			} `json:"audio"`
		} `json:"dash"`
	} `json:"data"`
}

func (b *BilibiliFetcher) fillHeader(req *netutil.HttpRequest) {
	req.Headers.Set("authority", "api.bilibili.com")
	req.Headers.Set("accept", "application/json, text/plain, */*")
	req.Headers.Set("accept-language", "zh-CN")
	req.Headers.Set("cookie", "buvid3=0")
	req.Headers.Set("referer", "https://www.bilibili.com/")
	req.Headers.Set("sec-fetch-dest", "empty")
	req.Headers.Set("sec-fetch-mode", "cors")
	req.Headers.Set("sec-fetch-site", "cross-site")
	req.Headers.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36")
}

func (bf *BilibiliFetcher) GetAudioURL(bvid string) []string {
	req := netutil.HttpRequest{
		RawURL:  fmt.Sprintf("https://api.bilibili.com/x/web-interface/view?bvid=%s", bvid),
		Method:  http.MethodGet,
		Headers: make(http.Header),
	}
	bf.fillHeader(&req)

	response, err := bf.cli.SendRequest(&req)
	if err != nil {
		return nil
	}
	defer response.Body.Close()

	result, err := pkg.Request(bf.cli, &req, func(result *BiliDetail) bool {
		return result.Code == 0 || len(result.Data.Bvid) != 0 || result.Data.Cid != 0
	})
	if err != nil {
		return nil
	}

	req = netutil.HttpRequest{
		RawURL: fmt.Sprintf("https://api.bilibili.com/x/player/playurl?fnval=16&bvid=%s&cid=%d", result.Data.Bvid, result.Data.Cid),
		Method: http.MethodGet,
	}

	play, err := pkg.Request(bf.cli, &req, func(result *BiliPlayURL) bool {
		return result.Code == 0 || len(result.Data.Dash.Audio) != 0
	})
	if err != nil {
		return nil
	}

	var urls []string
	for _, u := range play.Data.Dash.Audio {
		urls = append(urls, u.BaseURL)
	}
	// 同一个bvid下面的音频基本重复，具体情况后面再看
	return urls[:1]
}

func (bf *BilibiliFetcher) Download(url string, writer io.Writer) error {
	req := netutil.HttpRequest{
		RawURL:  url,
		Method:  http.MethodGet,
		Headers: make(http.Header),
	}
	req.Headers.Set("accept", "*/*")
	req.Headers.Set("accept-language", "zh-CN")
	req.Headers.Set("referer", "https://www.bilibili.com/")
	req.Headers.Set("sec-fetch-dest", "audio")
	req.Headers.Set("sec-fetch-mode", "no-cors")
	req.Headers.Set("sec-fetch-site", "cross-site")
	req.Headers.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36")
	req.Headers.Set("range", "bytes=0-")

	response, err := bf.cli.SendRequest(&req)
	if err != nil {
		return err
	}
	if response.StatusCode != http.StatusOK && response.StatusCode != http.StatusPartialContent {
		return pkg.ErrHttpRequest
	}
	defer response.Body.Close()

	if _, err := io.Copy(writer, response.Body); err != nil {
		return err
	}
	return err
}
