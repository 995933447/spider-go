package parser

import (
	"fmt"
	"regexp"
	"spider-go/process/model"
	julypornConfig "spider-go/process/app/julyporn/config"
	"spider-go/config"
	"spider-go/parseutil"
	"strconv"
	"util/timer"
)

var (
	videoNameRe = regexp.MustCompile(`<h4 class="container-title py-3 mb-0">([^<]+)</h4>`)
	videoCollectNumRe = regexp.MustCompile(`<span class="favoriteBtn cursor-p mx-2"[^>]*>[^<]*<i class="fas fa-heart "></i>[^<]*<span>([0-9]+)</span>`)
	videoSeeNumRe = regexp.MustCompile(`<span>免费线路</span>[^<]*</div>[^<]*</div>[^<]*<div>[^0-9]([0-9]+) 次访问[^<]*</div>`)
	videoImgRe = regexp.MustCompile(`<meta property="og:image" content="([^"]+)">`)
	videoM3u8 = regexp.MustCompile(`<video id="video-play" class="video-js vjs-default-skin vjs-big-play-centered" data-poster="[^"]+"[^a-z]*data-src="([^"]+)" data-type="application/x-mpegURL">`)
)

func Video(content []byte, long string, categoryName string, url string) config.ParseResult {
	var result config.ParseResult

	name := parseutil.ExtractString(videoNameRe, content)
	collectNum := parseutil.ExtractString(videoCollectNumRe, content)
	seeNum := parseutil.ExtractString(videoSeeNumRe, content)
	mainImg := parseutil.ExtractString(videoImgRe, content)
	m3u8 := parseutil.ExtractString(videoM3u8, content)

	var video model.Videos

	video.CrawlerUrl = url
	video.Name = name
	video.CollectNum, _ = strconv.Atoi(collectNum)
	video.SeeNum, _ = strconv.Atoi(seeNum)
	video.Mainimg = mainImg
	video.OriginalM3u8 = fmt.Sprintf("https:%s", m3u8)
	video.Long,  _ = timer.DurationToSecond(long)

	if video.Long > 60 * 60 {
		video.LongType = model.VideoLongType
	} else if video.Long > 25 * 60 && video.Long < 60 * 60 {
		video.LongType = model.VideoMiddleType
	} else {
		video.LongType = model.VideoShortType
	}

	video.Status = model.VideoValidStatus

	result.Items = append(result.Items, julypornConfig.VideoItem{
		categoryName,
		video,
	})

	return result
}
