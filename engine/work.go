package engine

import (
	"spider-go/config"
)

func Work(request config.Request) (config.ParseResult, error) {
	if hasDuplicateWork(request) {
		return config.ParseResult{}, nil
	}

	content, err := request.Fetch(request.Url)
	if err != nil {
		return config.ParseResult{}, err
	}

	parseResult := request.Parse(content)
	request.ItemChan <- parseResult.Items

	workDone(request)

	return parseResult, nil
}

func hasDuplicateWork(request config.Request) bool {
	if exists, err := request.Distinctor.CheckIsFetched(request.Url); err != nil {
		return true
	} else if exists {
		return true
	}
	return false
}

func workDone(request config.Request) {
	request.Distinctor.RecordFetched(request.Url)
}