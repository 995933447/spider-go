package parser

import (
	"regexp"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/parseutil"
	tangccConfig "spider-go/process/app/tangcc/config"
	"spider-go/process/app/tangcc/persist"
	"spider-go/process/model"
)

var (
	topicUriRe = regexp.MustCompile(`<li[^>]*><a href="([^"]+)">[^<]+</a></li>`)
	topicNameRe = regexp.MustCompile(`<li[^>]*><a href="[^"]+">([^<]+)</a></li>`)
)

func Topic(content []byte) config.ParseResult {
	var result config.ParseResult

	topicUrisMatches := topicUriRe.FindAllSubmatch(content, -1)
	topicUris := parseutil.ParseMatchesToString(topicUrisMatches)
	topicNameMatches := topicNameRe.FindAllSubmatch(content, -1)
	topicNames := parseutil.ParseMatchesToString(topicNameMatches)

	for index, topicUri := range topicUris {
		name := topicNames[index]
		result.Requests = append(result.Requests, config.Request{
			Url: tangccConfig.Host + topicUri,
			Fetch: fetcher.Html,
			Parse: func(content []byte) config.ParseResult {
				return Categories(content)
			},
			ItemChan: persist.Categories(),
			Distinctor: config.NilDistinctor{},
		})
		result.Items = append(result.Items, model.Subjects{
			Name: name,
		})
	}
	return result
}
