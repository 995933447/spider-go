package parser

import (
	"regexp"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/parseutil"
	tangccConfig "spider-go/process/app/tangcc/config"
	"spider-go/process/app/tangcc/itemworker"
	"spider-go/process/distinctor"
	"strconv"
	"util/rand"
)

var (
	VideoUriRe = regexp.MustCompile(`<li class="video ">[^<]+<a href="([^"]+)" target="_blank" title="[^"]+">`)
	VideoNameRe = regexp.MustCompile(`<li class="video ">[^<]+<a href="[^"]+" target="_blank" title="([^"]+)">`)
	VideoImgRe = regexp.MustCompile(`<div class="thumb"><img src="([^"]+)" alt="[^"]+" /></div>`)
	VideoSeeNumRe = regexp.MustCompile(`<span class="f-r fa fa-eye"> 已被观看：([0-9]+) 次</span>`)
	VideoLikeNumRe = regexp.MustCompile(`<span class="f-l fa fa-thumbs-o-up" style="margin-left: 15px;">&nbsp;([0-9]+)</span>`)
)

func Videos(content []byte, categoryName string) config.ParseResult {
	var result config.ParseResult

	uriMatches := VideoUriRe.FindAllSubmatch(content, -1)
	uris := parseutil.ParseMatchesToString(uriMatches)

	nameMatches := VideoNameRe.FindAllSubmatch(content, -1)
	names := parseutil.ParseMatchesToString(nameMatches)

	imgMatches := VideoImgRe.FindAllSubmatch(content, -1)
	imgs := parseutil.ParseMatchesToString(imgMatches)

	seeNumMatches := VideoSeeNumRe.FindAllSubmatch(content, -1)
	seeNums := parseutil.ParseMatchesToString(seeNumMatches)

	likeNumMatches := VideoLikeNumRe.FindAllSubmatch(content, -1)
	likeNums := parseutil.ParseMatchesToString(likeNumMatches)


	for index, uri := range uris {
		var item tangccConfig.VideoItem
		item.Name = names[index]
		item.MainImg = tangccConfig.Host + imgs[index]
		if len(seeNums) < (index - 1) || len(seeNums) == 0 {
			item.SeeNum = rand.RandIntN(10000)
		} else {
			seeNum := seeNums[index]
			item.SeeNum, _ = strconv.Atoi(seeNum)
		}
		if len(likeNums) < (index - 1) || len(seeNums) == 0 {
			item.LikeNum = rand.RandIntN(1000)
		} else {
			likeNum := likeNums[index]
			item.LikeNum, _ = strconv.Atoi(likeNum)
		}
		item.CrawlerUrl = tangccConfig.Host + uri
		item.CategoryName = categoryName

		result.Requests = append(result.Requests, config.Request{
			Url: item.CrawlerUrl,
			Fetch: fetcher.Html,
			Parse: func(content []byte) config.ParseResult {
				return Video(content, item)
			},
			ItemChan: itemworker.Videos(),
			Distinctor: &distinctor.Video{},
		})
	}

	return result
}
