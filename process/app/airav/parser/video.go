package parser

import (
	"regexp"
	airavConfig "spider-go/process/app/airav/config"
	"spider-go/config"
	"spider-go/parseutil"
	"strconv"
)

var (
	videoNameRe = regexp.MustCompile(`<span id="ContentPlaceHolder1_Label1" itemprop="name">([^<]+)</span></h1>`)
	videoImgRe = regexp.MustCompile(`<img id="ContentPlaceHolder1_Image_itemscope" itemprop="image" src="([^"]+)" />`)
	videoSeeNumRe = regexp.MustCompile(`<span id="ContentPlaceHolder1_Label4">([0-9]+)</span></span>`)
	m3u8Re = regexp.MustCompile(`file: "([^"]+index.m3u8")`)
)

func Video(content []byte, categoryName string) config.ParseResult {
	var result config.ParseResult
	name := parseutil.ExtractString(videoNameRe, content)
	mainImg := parseutil.ExtractString(videoImgRe, content)
	seeNum, _ := strconv.Atoi(parseutil.ExtractString(videoSeeNumRe, content))
	m3u8 := parseutil.ExtractString(m3u8Re, content)

	result.Items = append(result.Items, airavConfig.VideoItem{
		CategoryName: categoryName,
		Name: name,
		MainImg: mainImg,
		SeeNum: seeNum,
		OriginalM3u8: m3u8,
	})

	return result
}