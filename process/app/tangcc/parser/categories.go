package parser

import (
	"regexp"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/parseutil"
	"spider-go/process/model"
	tangccConfig "spider-go/process/app/tangcc/config"
)

var(
	categoryNameRe = regexp.MustCompile(`<a href="/index.php\?[^"]*">([^<]+)</a>`)
	categoryUriRe = regexp.MustCompile(`<a href="(/index.php\?[^"]*)">[^<]+</a>`)
)

func Categories(content []byte) config.ParseResult {
	var result config.ParseResult

	nameMatches := categoryNameRe.FindAllSubmatch(content, -1)
	names := parseutil.ParseMatchesToString(nameMatches)

	uriMatches := categoryUriRe.FindAllSubmatch(content, -1)
	uris := parseutil.ParseMatchesToString(uriMatches)

	for index, name := range names {
		result.Items = append(result.Items, model.Categories{
			Name: name,
		})

		url := tangccConfig.Host + uris[index]
		result.Requests = append(result.Requests, config.Request{
			Url: url,
			Fetch: fetcher.Html,
			Parse: func(content []byte) config.ParseResult {
				return VideoPages(content, url, name)
			},
			ItemChan: config.NilItemWork(),
			Distinctor: config.NilDistinctor{},
		})
	}

	return result
}