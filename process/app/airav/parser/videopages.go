package parser

import (
	"regexp"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/parseutil"
	"strconv"
)

var pagesRe = regexp.MustCompile(`<span id="ContentPlaceHolder1_Label2">([0-9]+)</span></h1>`)

func VideoPages(content []byte, categoryName string, url string) config.ParseResult {
	var result config.ParseResult

	pages, _ := strconv.Atoi(parseutil.ExtractString(pagesRe, content))

	if pages > 1 {
		pages = 1
	}

	for i := 1; i <= pages; i++ {
		result.Requests = append(result.Requests, config.Request{
			Url: url + "&idx=" + strconv.Itoa(i),
			Fetch: fetcher.Html,
			Parse: func(content []byte) config.ParseResult {
				return Videos(content, categoryName)
			},
			ItemChan: config.NilItemWork(),
			Distinctor: config.NilDistinctor{},
		})
	}

	return result
}