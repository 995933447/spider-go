package parser

import (
	"regexp"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/parseutil"
	"strconv"
)

var pageRe = regexp.MustCompile(`<li><a href="[^"]+">([0-9]+)</a></li>`)

func VideoPages(content []byte, categoryUrl string, categoryName string) config.ParseResult {
	var result config.ParseResult
	pageMatches := pageRe.FindAllSubmatch(content, -1)
	var lastPage int
	if len(pageMatches) <= 0 {
		lastPage = 1
	} else {
		pages := parseutil.ParseMatchesToString(pageMatches)
		lastPage, _ = strconv.Atoi(pages[len(pages) - 1])
		for i:= 1; i <= lastPage; i++ {
			result.Requests = append(result.Requests, config.Request{
				Url: categoryUrl + "&page=" + strconv.Itoa(i),
				Fetch: fetcher.Html,
				Parse: func(content []byte) config.ParseResult {
					return Videos(content, categoryName)
				},
				ItemChan: config.NilItemWork(),
				Distinctor: config.NilDistinctor{},
			})
		}
	}

	return result
}
