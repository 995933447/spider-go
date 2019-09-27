package parser

import (
	"regexp"
	"spider-go/config"
	"spider-go/parseutil"
	xiaoniaoConfig "spider-go/process/app/xiaoniao/config"
	"strconv"
	"strings"
)

var videoNameRe = regexp.MustCompile(`<div id="video" class="row">[^<]*<h2>([^<]+)</h2>`)
var m3u8Re = regexp.MustCompile(`,"url":"([^"]+m3u8)"`)
var likeNumRe = regexp.MustCompile(`(\([0-9]+\))`)

func Video(content []byte, categoryName string, mainimg string, url string) config.ParseResult {
	var result config.ParseResult

	name := parseutil.ExtractString(videoNameRe, content)
	m3u8 := strings.ReplaceAll(parseutil.ExtractString(m3u8Re, content), `\/`, `/`)
	like := parseutil.ExtractString(likeNumRe, content)
	
	var item xiaoniaoConfig.VideoItem
	item.M3u8 = m3u8
	item.Like, _ = strconv.Atoi(like)
	item.Name = name
	item.Category = categoryName
	item.CrawlerUrl = url
	item.Mainimg = mainimg

	result.Items = append(result.Items, item)

	return result
}
