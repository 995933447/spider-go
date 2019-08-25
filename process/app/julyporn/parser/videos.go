package parser

import (
	"regexp"
	"spider-go/process/distinctor"
	julypornConfig "spider-go/process/app/julyporn/config"
	"spider-go/process/app/julyporn/itemworker"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/parseutil"
	"strings"
)

var (
	uriRe = regexp.MustCompile(`<a target="_blank" class="title text-sub-title mt-2 mb-3" href="([^"]+)">[^<]+</a>`)
	longRe = regexp.MustCompile(`<small class="layer">([^<]+)</small>`)
)

func Videos(content []byte, categoryName string) config.ParseResult {
	var result config.ParseResult
	uriMatches := uriRe.FindAllSubmatch(content, -1)
	uris := parseutil.ParseMatchesToString(uriMatches)
	longMatches := longRe.FindAllSubmatch(content, -1)
	longs := parseutil.ParseMatchesToString(longMatches)

	for index := range uris {
		long := strings.Trim(strings.Trim(longs[index], "\n"), " ")
		url := julypornConfig.Host + uris[index]
		result.Requests = append(result.Requests, config.Request{
			Url: url,
			Fetch: fetcher.Html,
			Parse: func(content []byte) config.ParseResult {
				return Video(content, long, categoryName, url)
			},
			ItemChan:   itemworker.Videos(),
			Distinctor: &distinctor.Video{},
		})
	}
	return result
}
