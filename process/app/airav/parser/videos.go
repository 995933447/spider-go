package parser

import (
	"regexp"
	"spider-go/process/distinctor"
	airavConfig "spider-go/process/app/airav/config"
	"spider-go/process/app/airav/itemworker"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/parseutil"
)

var uriRe = regexp.MustCompile(`<a class="ga_click" href='([^']+)' title="[^"]+">`)

func Videos(content []byte, categoryName string) config.ParseResult {
	var result config.ParseResult

	uriMatches := uriRe.FindAllSubmatch(content, -1)
	uris := parseutil.ParseMatchesToString(uriMatches)

	for _, uri := range uris {
		result.Requests = append(result.Requests, config.Request{
			Url:   airavConfig.Host + "/" + uri,
			Fetch: fetcher.Html,
			Parse: func(content []byte) config.ParseResult {
				return Video(content, categoryName)
			},
			ItemChan:   itemworker.Videos(),
			Distinctor: &distinctor.Video{},
		})
	}

	return result
}
