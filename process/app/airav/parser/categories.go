package parser

import (
	"regexp"
	"spider-go/process/model"
	airavConfig "spider-go/process/app/airav/config"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/parseutil"
)

var (
	categoryUriRe = regexp.MustCompile(`<li class="AirMenu_li[^A-Z]*(Hide_For_JP_AND_EN)*"><a href="(/index.aspx\?status=[0-9])">[^<]*<span id="[^"]+" class="notranslate">[^<]+</span></a></li>`)
	categoryNameRe = regexp.MustCompile(`<li class="AirMenu_li[^A-Z]*(Hide_For_JP_AND_EN)*"><a href="/index.aspx\?status=[0-9]">[^<]*<span id="[^"]+" class="notranslate">([^<]+)</span></a></li>`)
)

func Categories(content []byte) config.ParseResult {
	var result config.ParseResult
	nameMatches := categoryNameRe.FindAllSubmatch(content, -1)
	names := parseutil.ParseMatchesToString(nameMatches)

	uriMatches := categoryUriRe.FindAllSubmatch(content, -1)
	uris := parseutil.ParseMatchesToString(uriMatches)

	num := 1
	for index, uri := range uris {
		if num <= 0 {
			continue
		}
		name := names[index]
		url := airavConfig.Host + uri
		result.Requests = append(result.Requests, config.Request{
			Url: url,
			Fetch: fetcher.Html,
			Parse: func(content []byte) config.ParseResult {
				return VideoPages(content, name, url)
			},
			ItemChan: config.NilItemWork(),
			Distinctor: config.NilDistinctor{},
		})

		result.Items = append(result.Items, model.Categories{
			Name: name,
		})
		num--
	}

	return result
}
