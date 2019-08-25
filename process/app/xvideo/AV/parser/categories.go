package parser

import (
	"spider-go/process/distinctor"
	"spider-go/process/model"
	"spider-go/process/app/xvideo/aichange/itemworker"
	"spider-go/config"
	"spider-go/fetcher"
)

func Categories(content []byte) config.ParseResult {
	var result config.ParseResult
	result.Items = append(result.Items, model.Categories{
		Name: "女优",
		Status: model.CategoryValidStatus,
	})
	urls := []string{
		"https://www.xvideos.com/video39860959/sexy_office_bitch_yui_hatano_pussy_pounded_by_co_workers",
		"https://www.xvideos.com/video40462833/handsome_office_gal_yui_hatano_pussy_drilled_by_her_colleagues",
		"https://www.xvideos.com/video16730161/slippery_and_moist_oriental_group_sex",
		"https://www.xvideos.com/video35979997/bao_dam_vo_yeu",
		"https://www.xvideos.com/video23426304/akiho_yoshizawa_hotest_1st_hard_pony_ride_cock_hard",
	}

	for index := range urls {
		url := urls[index]
		result.Requests = append(result.Requests, config.Request{
			Url: url,
			Fetch: fetcher.Html,
			Parse: func(content []byte) config.ParseResult {
				return Video(content, url)
			},
			ItemChan: itemworker.Video(),
			Distinctor: &distinctor.Video{},
		})
	}
	return result
}
