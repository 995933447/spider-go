package parser

import (
	"regexp"
	"spider-go/process/distinctor"
	javdoveConfig "spider-go/process/app/javdove/config"
	"spider-go/process/app/javdove/itemworker"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/parseutil"
	"strings"
	"util/timer"
)

var videoMainImgRe = regexp.MustCompile(`(<div class="thumb-overlay">[^<]*<img src="([^?]+)[^"]*"[^>]+>)|(<img data-original="([^?]+)[^"]*"[^>]*>)`)
var videoUrlRe = regexp.MustCompile(`<div class="well well-sm">[^<]*<a href="([^"]+)">`)
var videoLongRe = regexp.MustCompile(`<div class="duration">([^<]+)</div>`)

var videoItemChan = itemworker.Videos(javdoveConfig.VideoProcessNum, javdoveConfig.FetchingVideoNum)

func Videos(content []byte) config.ParseResult {
	var result config.ParseResult

	urlMatches := videoUrlRe.FindAllSubmatch(content, -1)
	urls := parseutil.ParseMatchesToString(urlMatches)

	longMatches := videoLongRe.FindAllSubmatch(content, -1)
	longs := parseutil.ParseMatchesToString(longMatches)

	mainImgMatches := videoMainImgRe.FindAllSubmatch(content, -1)
	var mainImgs []string
	for _, mainImgMatch := range mainImgMatches {
		if string(mainImgMatch[4]) == "" {
			mainImgs = append(mainImgs, string(mainImgMatch[2]))
		} else {
			mainImgs = append(mainImgs, string(mainImgMatch[4]))
		}
	}

	for index, url := range urls {
		url := javdoveConfig.Host + url
		long, _ := timer.DurationToSecond(strings.Trim(strings.Trim(longs[index], "\n"), " "))
		mainImg := mainImgs[index]
		result.Requests = append(result.Requests, config.Request{
			Url: url,
			Fetch: fetcher.Html,
			Parse: func(content []byte) config.ParseResult {
				return Video(content, mainImg, long, url)
			},
			ItemChan: videoItemChan,
			Distinctor: &distinctor.Video{},
		})
	}

	return result
}
