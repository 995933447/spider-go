package parser

import (
	"spider-go/config"
	"spider-go/fetcher"
	javhdConfig "spider-go/process/app/xvideo/javhd/config"
	"strconv"
)

func VideoPages(content []byte) config.ParseResult {
	var result config.ParseResult

	// Total pages 149
	for i := 1; i < 150; i ++ {
		result.Requests = append(result.Requests, config.Request{
			Url: javhdConfig.Host + "/?k=javhd&p=" + strconv.Itoa(i),
			Fetch: fetcher.Html,
			Parse: Videos,
			ItemChan: config.NilItemWork(),
			Distinctor: config.NilDistinctor{},
		})
	}

	return result
}
