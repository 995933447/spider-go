package parser

import (
	"regexp"
	"spider-go/process/distinctor"
	avhd101Config "spider-go/process/app/avhd101/config"
	avhd101Fetcher "spider-go/process/app/avhd101/fetcher"
	"spider-go/process/app/avhd101/itemworker"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/parseutil"
	"strconv"
)

var (
	currentPageRe = regexp.MustCompile(`<link rel="alternate" hreflang="zh-Hans" href="https://avhd101.com/cn/search\?q.*?&ft=&p=([0-9]+)" />`)
	pageRe = regexp.MustCompile(`<a href="(/search\?q=.*?&ft=&p=[0-9]+)">([0-9]+)</a>`)
	videoRe = regexp.MustCompile(`<a href="(/watch\?v=[^"]+)">`)
)

func Videos(content []byte, categoryName string, subjectName string) config.ParseResult {
	var result config.ParseResult

	findMoreVideoPage(&content, &result, &categoryName, &subjectName)

	videoMatches := videoRe.FindAllSubmatch(content, -1)
	for _, videoMatch := range videoMatches {
		url := avhd101Config.Host + string(videoMatch[1])

		result.Requests = append(result.Requests, config.Request{
			Url:   url,
			Fetch: avhd101Fetcher.Video,
			Parse: func(content []byte) config.ParseResult {
				return Video(content, categoryName, subjectName, url)
			},
			ItemChan:   itemworker.Videos(),
			Distinctor: &distinctor.Video{},
		})
	}

	return result
}

func findMoreVideoPage(content *[]byte, result *config.ParseResult, categoryName, subjectName *string) {
	currentPageMatch := parseutil.ExtractString(currentPageRe, *content)
	currentPage := 1
	if currentPageMatch != "" {
		currentPage, _ = strconv.Atoi(currentPageMatch)
	}

	pageMatches := pageRe.FindAllSubmatch(*content, -1)
	for _, pageMatch := range pageMatches {
		page, _ := strconv.Atoi(string(pageMatch[2]))
		if page > currentPage {
			result.Requests = append(result.Requests, config.Request{
				Url:        avhd101Config.Host + string(pageMatch[1]),
				Fetch:      fetcher.Html,
				Parse:      MakeVideosParse(*categoryName, *subjectName),
				ItemChan:   config.NilItemWork(),
				Distinctor: config.NilDistinctor{},
			})
		}
	}
}

func MakeVideosParse(categoryName, subjectName string) func([]byte) config.ParseResult {
	return func (content []byte) config.ParseResult {
		return Videos(content , categoryName, subjectName)
	}
}