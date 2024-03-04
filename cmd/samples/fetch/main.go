package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/duke-git/lancet/v2/netutil"
	"github.com/o98k-ok/voice/internal/bilibili"
	"github.com/olekukonko/tablewriter"
)

func main() {
	fetcher := bilibili.NewBlibliFetcher(netutil.NewHttpClient())
	all, err := fetcher.Search("以父之名", 1, 10)
	if err != nil {
		fmt.Println("error happen", err)
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"标题", "描述", "bvid"})

	// 格式存在问题，不想调了
	for _, a := range all {
		if len(a.Desc) > 100 {
			a.Desc = a.Desc[:100]
		}
		a.Desc = strings.ReplaceAll(a.Desc, "\n", "")
		table.Append([]string{a.Name, a.Desc, a.URL})
	}
	table.Render()
}
