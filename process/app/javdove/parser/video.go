package parser

import (
	"regexp"
	"spider-go/process/model"
	javdoveConfig "spider-go/process/app/javdove/config"
	"spider-go/config"
	"spider-go/parseutil"
	"strconv"
	"strings"
)

var videoNameRe = regexp.MustCompile(`<div class="video-info-wrap">[^<]*<h4>([^<]+)</h4>`)
var videoSeeNumRe = regexp.MustCompile(`<div class="video-info-section hidden-xs">[^<]*<p>观看次数：<span>([0-9]+)</span></p>`)
var videoCollectNumRe = regexp.MustCompile(`<p id="video_likes" style="display: inline-block;text-decoration: none;">([0-9]+)</p>`)
var videoMp4Re = regexp.MustCompile(`<source src="([^"]+)" type='video/mp4'>`)
var videoCategoryRe = regexp.MustCompile(`<span property="name">全部免費A片</span>[^<]*</a>[^<]*<meta[^>]*>[^<]*</span>[^<]*<span>[^<]+</span>[^<]*<span property="itemListElement"[^>]*>[^<]*<a property="item"[^>]*>[^<]*<span property="name">([^<]+)</span>`)
var videoM3u8Re = regexp.MustCompile(`<source src="([^?]+\.m3u8\?e=[^"]+)" type='application/x-mpegURL'>`)

func Video(content []byte, mainImg string, long int, url string) config.ParseResult {
	var result config.ParseResult

	video := model.Videos{}
	video.Name = strings.ReplaceAll(parseutil.ExtractString(videoNameRe, content), "番号鸽", "")
	video.SeeNum, _ = strconv.Atoi(parseutil.ExtractString(videoSeeNumRe, content))
	video.CollectNum, _ = strconv.Atoi(parseutil.ExtractString(videoCollectNumRe, content))
	video.OriginalMp4 = parseutil.ExtractString(videoMp4Re, content)
	video.OriginalM3u8 = parseutil.ExtractString(videoM3u8Re, content)
	video.Long = long
	if long > 60 * 60 {
		video.LongType = model.VideoLongType
	} else if long > 25 * 60 && long < 60 * 60 {
		video.LongType = model.VideoMiddleType
	} else {
		video.LongType = model.VideoShortType
	}
	video.Mainimg = mainImg
	video.CrawlerUrl = url

	item := javdoveConfig.VideoItem{
		Category: parseutil.ExtractString(videoCategoryRe, content),
		Model:    video,
	}

	result.Items = append(result.Items, item)
	return result
}
