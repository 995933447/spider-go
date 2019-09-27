package parser

import (
	"regexp"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/parseutil"
	xiaoniaoConfig "spider-go/process/app/xiaoniao/config"
)

var pageUriRe = regexp.MustCompile(`<li class="mbyc"><a class="pagelink_b" href="([^"]+)" title="[^"]+">[^<]+</a></li>`)

func VideoPages(content []byte, categoryName string) config.ParseResult {
	var result config.ParseResult

	pageUriMatches := pageUriRe.FindAllSubmatch(content, -1)
	pageUris := parseutil.ParseMatchesToString(pageUriMatches)

	for _, pageUri := range pageUris {
		result.Requests = append(result.Requests, config.Request{
			Url: xiaoniaoConfig.Host + pageUri,
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
