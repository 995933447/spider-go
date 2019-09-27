package parser

import (
	"fmt"
	"regexp"
	"spider-go/config"
	"spider-go/logger"
	"spider-go/parseutil"
	javhdConfig "spider-go/process/app/xvideo/javhd/config"
)

var (
	titleRe = regexp.MustCompile(`<h2 class="page-title">([^<]+)<span class="duration">[0-9]+ min</span>`)
	longRe = regexp.MustCompile(`<h2 class="page-title">[^<]+<span class="duration">([^<]+)</span>`)
	imgRe = regexp.MustCompile(`html5player.setThumbUrl169\('([^']+)'\);`)
	//m3u8Re = regexp.MustCompile(`html5player.setVideoHLS\('([^']+)'\)`)
	mp4Re = regexp.MustCompile(`html5player.setVideoUrlHigh\('([^']+)'\);`)
)

func Video(content []byte, url string) config.ParseResult {
	var result config.ParseResult
	title := parseutil.ExtractString(titleRe, content)
	long := parseutil.ExtractString(longRe, content)
	img := parseutil.ExtractString(imgRe, content)
	//m3u8 := parseutil.ExtractString(m3u8Re, content)
	mp4 := parseutil.ExtractString(mp4Re, content)

	result.Items = append(result.Items, javhdConfig.Video{
		Name: title,
		Long: long,
		MainImg: img,
		Mp4: mp4,
		CrawlerUrl: javhdConfig.Host + url,
	})
	logger.DefaultLogger.Info(fmt.Sprintf("%+v\n", result.Items), nil)
	return result
}
