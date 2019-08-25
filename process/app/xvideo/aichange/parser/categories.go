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
		Name: "AI换脸",
		Status: model.CategoryValidStatus,
	})
	urls := []string{
		"https://www.xvideos.com/video49562097/ai_",
		"https://www.xvideos.com/video49869765/_yuxinyalove",
		"https://www.xvideos.com/video50100269/sex_korea_-_trao_oi_vo_cho_nhau_chich_p3",
		"https://www.xvideos.com/video9774463/_phimse.net_k-pop_sex_scandal_korean_celebrities_prostituting_vol_17",
		"https://www.xvideos.com/video46101621/_706998358",
		"https://www.xvideos.com/video49548561/_sex_video_part_1_ai_hotel66",
		"https://www.xvideos.com/video49171653/_hotel66",
		"https://www.xvideos.com/video42769395/_---_",
		"https://www.xvideos.com/video48130903/_sex_video_ai_-_deepfakesporn.com",
		"https://www.xvideos.com/video47737807/china_chinese_yang_mi_hot_deepfakes_riding_sex_in_japan_-_deepfakesporn.com",
		"https://www.xvideos.com/video12546871/_",
		//"https://www.xvideos.com/video18486803/_-_",
		//"https://www.xvideos.com/video25570075/_1_-_-_thisav.com-_",
		//"https://www.xvideos.com/video48130919/_sex_video_part_2_ai_-_deepfakesporn.com",
		//"https://www.xvideos.com/video42229811/_",
		//"https://www.xvideos.com/video11630121/g.e.m_",
		//"https://www.xvideos.com/video8632055/hk_movie_sex_scene_akotube.com",
		//"https://www.xvideos.com/video38012961/_",
		//"https://www.xvideos.com/video18955147/_",
		//"https://www.xvideos.com/video46035213/_706998358",
		//"https://www.xvideos.com/video44748959/_",
		//"https://www.xvideos.com/video42967897/_",
		//"https://www.xvideos.com/video46101749/_706998358",
		//"https://www.xvideos.com/video49580039/_",
		//"https://www.xvideos.com/video30336405/_partb",
		//"https://www.xvideos.com/video47864567/kpop_red_velvet_wendy_sex_-_kpopdeepfakes.com_for_more_videos",
		//"https://www.xvideos.com/video31973571/fucking_a_moaning_big-tits_japanese_girl",
		//"https://www.xvideos.com/video12759527/_-6",
		//"https://www.xvideos.com/video42767247/_",
		//"https://www.xvideos.com/video28929745/xvideos.com_15e4093628d58e207eb57e38c2b323d5-2.mp4",
		//"https://www.xvideos.com/video46094191/_ai_angelababy_",
		//"https://www.xvideos.com/video14210035/_..._1_",
		//"https://www.xvideos.com/video7918805/beautiful_chinese_girl",
		//"https://www.xvideos.com/video48654993/nancy_korea_film_uncenco_link_full_1shortlink.com_link_yclh0yhz7_pass_123_",
		//"https://www.xvideos.com/video22045291/chinese_forced_sex_part_3_",
		//"https://www.xvideos.com/video45873631/_iu_",
		//"https://www.xvideos.com/video46317131/chinese_beautiful_model_compilation",
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
