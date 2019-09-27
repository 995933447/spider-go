package parser

import (
	"regexp"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/parseutil"
	javhdConfig "spider-go/process/app/xvideo/javhd/config"
	"spider-go/process/app/xvideo/javhd/itemworker"
	"spider-go/process/distinctor"
)

var urlRe = regexp.MustCompile(`<p class="title"><a href="([^"]+)" title="[^"]+">[^<]+</a></p>`)

func Videos(content []byte) config.ParseResult {
	var result config.ParseResult
	urlMatches := urlRe.FindAllSubmatch(content, -1)
	urls := parseutil.ParseMatchesToString(urlMatches)
	for _, url := range urls {
		result.Requests = append(result.Requests, config.Request{
			Url: javhdConfig.Host + url,
			Fetch: fetcher.Html,
			Parse: func(content []byte) config.ParseResult {
				return Video(content, url)
			},
			ItemChan: itemworker.Video(),
			Distinctor: &distinctor.Video{},
		})
	}

	return result
}
