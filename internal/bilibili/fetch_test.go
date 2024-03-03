package bilibili

import (
	"testing"

	"github.com/duke-git/lancet/v2/netutil"
)

func TestFetch(t *testing.T) {
	fetcher := BilibiliFetcher{
		cli: netutil.NewHttpClient(),
	}

	all, err := fetcher.Search("If a were a boy", 1, 10)
	if err != nil {
		t.Error(err)
	}

	for _, a := range all {
		t.Log(fetcher.GetAudioURL(a.URL))
	}
}
