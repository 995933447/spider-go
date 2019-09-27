package parser

import (
	"regexp"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/parseutil"
	xiaoniaoConfig "spider-go/process/app/xiaoniao/config"
	"spider-go/process/app/xiaoniao/itemworker"
	"spider-go/process/distinctor"
)

var uriRe = regexp.MustCompile(`<a href="([^"]+)" title="[^"]+" class="thumbnail" target="_blank">`)
var mainimgRe = regexp.MustCompile(`<img onerror="this.src='[^']+'" src="([^"]+)" alt="[^"]+" width="200" height="135">`)

func Videos(content []byte, categoryName string) config.ParseResult {
	var result config.ParseResult

	uriMatches := uriRe.FindAllSubmatch(content, -1)
	uris := parseutil.ParseMatchesToString(uriMatches)

	mainimgMathes := mainimgRe.FindAllSubmatch(content, -1)
	mainimgs := parseutil.ParseMatchesToString(mainimgMathes)

	for index, uri := range uris {
		url := xiaoniaoConfig.Host + uri
		mainimg := mainimgs[index]
		result.Requests = append(result.Requests, config.Request{
			Url: url,
			Fetch: fetcher.Html,
			Parse: func(content []byte) config.ParseResult {
				return Video(content, categoryName, mainimg, url)
			},
			ItemChan: itemworker.Videos(),
			Distinctor: &distinctor.Video{},
		})
	}
	
	return result
}
