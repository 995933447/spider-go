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
		Name: "网友推荐",
		Status: model.CategoryValidStatus,
	})
	urls := []string{
		"https://www.xvideos.com/video35378047/_",
		"https://www.xvideos.com/video26224731/bat_coc_hiep_dam",
		"https://www.xvideos.com/video31482161/_2_.",
		"https://www.xvideos.com/video44713519/sexy_teen_emo_with_big_eyes_trying_deepthroat_anal",
		"https://www.xvideos.com/video50119849/_",
		"https://www.xvideos.com/video48460677/_a",
		"https://www.xvideos.com/video43564455/_",
		"https://www.xvideos.com/video45636539/_",
		"https://www.xvideos.com/video44689829/_",
		"https://www.xvideos.com/video45228705/runa_momose_amazing_scenes_of_threesome_sex_-_more_at_javhd.net",
		"https://www.xvideos.com/video34622765/fucked_hard_gangbang._who_is_she",
		"https://www.xvideos.com/video15607209/salacious_and_wild_oriental_bang",
		"https://www.xvideos.com/video37644341/friend_s_mothers_2018_720p_hd_korea_18_",
		"https://www.xvideos.com/video48930573/amazing_korean_camgirl_dancing_and_masturbating_part2",
		"https://www.xvideos.com/video38307017/korean_sex",
		"https://www.xvideos.com/video48029249/_",
		"https://www.xvideos.com/video48288378/_",
		"https://www.xvideos.com/video35921705/chich_em_erika_momtani_dam_ang",
		"https://www.xvideos.com/video12642297/japanese_cougar_teacher_fucked_by_her_student",
		"https://www.xvideos.com/video45996077/_39107893_",
		"https://www.xvideos.com/video28735349/korea",
		"https://www.xvideos.com/video29109687/720p_1500k_115773481",
		"https://www.xvideos.com/video46364665/_",
		"https://www.xvideos.com/video12476931/_18_",
		"https://www.xvideos.com/video39860959/sexy_office_bitch_yui_hatano_pussy_pounded_by_co_workers",
		"https://www.xvideos.com/video35979997/bao_dam_vo_yeu",
		"https://www.xvideos.com/video10892360/bukkake",
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
