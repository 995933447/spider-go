package parser

import (
	"regexp"
	julypornConfig "spider-go/process/app/julyporn/config"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/parseutil"
	"strconv"
)

var (
	videosUriRe = regexp.MustCompile(`<a class="page-link" href="//julyporn.com/videos([^?]+)\?page=[0-9]+">[0-9]+</a>`)
	videosPageRe = regexp.MustCompile(`<a class="page-link" href="//julyporn.com/videos[^?]+\?page=([0-9]+)">[0-9]+</a>`)
)

func VideoPages(content []byte, categoryName string) config.ParseResult {
	var result config.ParseResult
	uri := parseutil.ExtractString(videosUriRe, content)
	pageMatches := videosPageRe.FindAllSubmatch(content, -1)
	pages := parseutil.ParseMatchesToString(pageMatches)
	var pageNum int
	if len(pages) > 1 {
		pageNum, _ = strconv.Atoi(pages[len(pages) - 1])
	} else {
		pageNum = 1
	}


	for i := 1; i < pageNum || i == pageNum; i++ {
		result.Requests = append(result.Requests, config.Request{
			Url:   julypornConfig.Host + "/videos" + uri + "?page=" + strconv.Itoa(i),
			Fetch: fetcher.Html,
			Parse: func(content []byte) config.ParseResult {
				return Videos(content, categoryName)
			},
			ItemChan:   config.NilItemWork(),
			Distinctor: config.NilDistinctor{},
		})
	}

	return result
}
