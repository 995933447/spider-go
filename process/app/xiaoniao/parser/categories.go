package parser

import (
	"regexp"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/parseutil"
	xiaoniaoConfig "spider-go/process/app/xiaoniao/config"
	"spider-go/process/model"
)

var categoryUriRe = regexp.MustCompile(`<li class=""><a class="nav-link" href="([^"]+)">[^<]+</a></li>`)
var categoryNameRe = regexp.MustCompile(`<li class=""><a class="nav-link" href="[^"]+">([^<]+)</a></li>`)

func Categories(content []byte) config.ParseResult {
	var result config.ParseResult
	categoryNameMatches := categoryNameRe.FindAllSubmatch(content, -1)
	categoryNames := parseutil.ParseMatchesToString(categoryNameMatches)
	categoryUriMatches := categoryUriRe.FindAllSubmatch(content, -1)
	categoryUris := parseutil.ParseMatchesToString(categoryUriMatches)
	for index, categoryUri := range categoryUris {
		if index == 0 {
			continue
		}
		categoryName := categoryNames[index]
		result.Requests = append(result.Requests, config.Request{
			Url: xiaoniaoConfig.Host + categoryUri,
			Parse: func(content []byte) config.ParseResult {
				return VideoPages(content, categoryName)
			},
			Fetch: fetcher.Html,
			ItemChan:config.NilItemWork(),
			Distinctor: config.NilDistinctor{},
		})

		result.Items = append(result.Items, model.Categories{
			Name: categoryName,
		})
	}
	return result
}