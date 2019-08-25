package parser

import (
	"regexp"
	"spider-go/process/model"
	avhd101Config "spider-go/process/app/avhd101/config"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/parseutil"
	"strings"
)

var (
	CategoryWithSubjectRe = regexp.MustCompile(`<article class="data-collapse">.*?</article>`)
	SubjectRe = regexp.MustCompile(`<a class="category" href="(/search\?q=[^"]+)" itemprop="url"><span itemprop="name">([^<]+)</span>[^<]+</a>`)
	CategoryRe = regexp.MustCompile(`<div class="title"> <h2><i class="fa fa-angle-double-right" aria-hidden="true"></i>([^<]+)</h2>`)
)

func CategoryWithSubjects(content []byte) config.ParseResult {
	var result config.ParseResult

	matches := CategoryWithSubjectRe.FindAllSubmatch(content, -1)
	for _, match := range matches {
		categoryName := parseutil.ExtractString(CategoryRe, match[0])
		 item := avhd101Config.CategoryWithSubjectItem{
			model.Categories{
				Name: strings.Trim(categoryName, " "),
			},
			[]model.Subjects{},
		}

		subjectMatches := SubjectRe.FindAllSubmatch(match[0], -1)
		for _, subjectMatch := range subjectMatches {
			subjectName := string(subjectMatch[2])
			item.Subjects = append(item.Subjects, model.Subjects{
				Name: strings.Trim(subjectName, " "),
			})

			result.Items = append(result.Items, item)
			result.Requests = append(result.Requests, config.Request{
				Url:        avhd101Config.Host + string(subjectMatch[1]),
				Fetch:      fetcher.Html,
				Parse:      MakeVideosParse(categoryName, subjectName),
				ItemChan:   config.NilItemWork(),
				Distinctor: config.NilDistinctor{},
			})
		}
	}

	return result
}


