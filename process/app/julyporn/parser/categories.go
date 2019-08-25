package parser

import (
	"regexp"
	"spider-go/process/model"
	julypornConfig "spider-go/process/app/julyporn/config"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/parseutil"
	"strings"
)

var (
	categoryUrlRe = regexp.MustCompile(`<div class="item col-30[^"]*" click="redirect" data-href="([^"]+)">[^<]+</div>`)
	categoryNameRe = regexp.MustCompile(`<div class="item col-30[^"]*" click="redirect" data-href="[^"]+">([^<]+)</div>`)
)

func Categories(content []byte) config.ParseResult {
	var result config.ParseResult

	urlMatches := categoryUrlRe.FindAllSubmatch(content, -1)
	urls := parseutil.ParseMatchesToString(urlMatches)
	urls = urls[1:]
	nameMatches := categoryNameRe.FindAllSubmatch(content, -1)
	names := parseutil.ParseMatchesToString(nameMatches)
	names = names[1:]

	for index := range urls {
		name := strings.Trim(names[index], " ")
		if name != "动漫" {
			continue
		}
		result.Requests = append(result.Requests, config.Request{
			Url:   julypornConfig.Host + urls[index],
			Fetch: fetcher.Html,
			Parse: func(content []byte) config.ParseResult {
				return VideoPages(content, name)
			},
			ItemChan:config.NilItemWork(),
			Distinctor: config.NilDistinctor{},
		})

		result.Items = append(result.Items, model.Categories{
			Name:   name,
			Status: model.CategoryValidStatus,
		})
	}

	return result
}
