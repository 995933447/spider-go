package parser

import (
	"regexp"
	"spider-go/config"
	"spider-go/parseutil"
	tangccConfig "spider-go/process/app/tangcc/config"
)

var (
	m3u8Re = regexp.MustCompile(`a:'(.+m3u8[^']*)'`)
	tagsRe = regexp.MustCompile(`<a href="[^"]+" class="blue">([^<]+)</a>`)
)

func Video(content []byte, item tangccConfig.VideoItem) config.ParseResult {
	var result config.ParseResult

	item.M3u8 = parseutil.ExtractString(m3u8Re, content)
	tagsMatches := tagsRe.FindAllSubmatch(content, -1)
	tags := parseutil.ParseMatchesToString(tagsMatches)
	for _, tag := range tags {
		item.Tags = append(item.Tags, tag)
	}

	result.Items = append(result.Items, item)

	return result
}
