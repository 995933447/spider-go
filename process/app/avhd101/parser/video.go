package parser

import (
	"regexp"
	avhd101Config "spider-go/process/app/avhd101/config"
	"spider-go/config"
	"spider-go/parseutil"
	"strings"
)

var baseDir  = strings.TrimRight(avhd101Config.BaseFileDir, "/")

var (
	nameRe = regexp.MustCompile(`<h1 class="container video-name">[^<]*?(<span[^>]+>[^<]+</span>.)*([^<]+)</h1>`)
	tagsHtmlRe = regexp.MustCompile(`<h1 class="container video-name">[^<]*?((<span[^>]+>([^<]+)</span>.)*)[^<]+</h1>`)
	tagRe = regexp.MustCompile(`<span[^>]+>([^<]+)</span>`)
	mainImgRe = regexp.MustCompile(`<video id="my-player" class="video-js vjs-big-play-centered vjs-16-9" controls poster="([^"]+)">`)
	m3u8FileReg = regexp.MustCompile(`<source src="([^"]+)" type="application/x-mpegURL" />`)
)

func Video(content []byte, categoryName, subjectName string, url string) config.ParseResult {
	var result config.ParseResult

	name := parseName(content)
	tags := parseTags(content)
	mainImg := strings.Trim(parseutil.ExtractString(mainImgRe, content), " ")
	m3u8File := strings.Trim(parseutil.ExtractString(m3u8FileReg, content), " ")

	item := avhd101Config.VideoItem{
		Url: url,
		Name: name,
		Tags: tags,
		MainImg: mainImg,
		M3u8File: m3u8File,
		Category: categoryName,
		Subject: subjectName,
	}

	result.Items = append(result.Items, item)

	return result
}

func parseName(content []byte) string {
	matches := nameRe.FindAllSubmatch(content, -1)
	return parseutil.ParseMatchesToString(matches)[0]
}

func parseTags(content []byte) []string {
	tagsHtmlMatches := tagsHtmlRe.FindAllSubmatch(content, -1)
	tagsHtml := parseutil.ParseMatchesToHtml(tagsHtmlMatches)
	if len(tagsHtml) > 0 {
		tagMatches := tagRe.FindAllSubmatch(tagsHtml[0], -1)
		return  parseutil.ParseMatchesToString(tagMatches)
	}
	return []string{}
}