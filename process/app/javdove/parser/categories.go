package parser

import (
	"regexp"
	javdoveConfig "spider-go/process/app/javdove/config"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/parseutil"
	"strconv"
)

var categoryIconRe = regexp.MustCompile(`<div class="thumb-overlay">[^<]*<img src="([^?]+.jpg)[^"]*"[^>]*>[^<]*</div>`)
var categoryNameRe = regexp.MustCompile(`<div class="classificagion-title m-t-5">[^<]*<button herf="">([^<]+)</button>[^<]*<span>[0-9]+</span>[^<]*</div>[^<]*</div>`)
var categoryVideoNumRe = regexp.MustCompile(`<div class="classificagion-title m-t-5">[^<]*<button herf="">[^<]+</button>[^<]*<span>([0-9]+)</span>[^<]*</div>[^<]*</div>`)
var categoryUrlRe = regexp.MustCompile(`<div class="col-xs-6 col-sm-6 col-md-4 col-lg-4 m-b-15 px-3">[^<]*<a href="([^"]+)">[^<]*<div class="thumb-overlay">`)

func Categories(content []byte) config.ParseResult {
	var result config.ParseResult

	iconMatches := categoryIconRe.FindAllSubmatch(content, -1)
	icons :=  parseutil.ParseMatchesToString(iconMatches)

	nameMatches := categoryNameRe.FindAllSubmatch(content, -1)
	names := parseutil.ParseMatchesToString(nameMatches)

	videoNumMatches := categoryVideoNumRe.FindAllSubmatch(content, -1)
	videoNums := parseutil.ParseMatchesToString(videoNumMatches)

	urlMatches := categoryUrlRe.FindAllSubmatch(content, -1)
	urls := parseutil.ParseMatchesToString(urlMatches)


	for index := range urls {
		videoNum, _ :=strconv.Atoi(videoNums[index])

		pages := videoNum / 16 * 3
		if (videoNum % 16 * 3) > 0 {
			pages++
		}

		url := urls[index]
		for i := 0; i <= pages; i++ {
			pageNum := strconv.Itoa(i)
			result.Requests = append(result.Requests, config.Request{
				Url:        javdoveConfig.Host + url + "?page=" + pageNum,
				Fetch:      fetcher.Html,
				Parse:      Videos,
				ItemChan:   config.NilItemWork(),
				Distinctor: config.NilDistinctor{},
			})
		}

		categoryName := names[index]
		icon := icons[index]

		result.Items = append(result.Items, javdoveConfig.CategoryItem{
			Name:     categoryName,
			Icon:     javdoveConfig.Host + icon,
			VideoNum: videoNum,
		})
	}

	return result
}