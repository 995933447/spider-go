package persist

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"spider-go/logger"
	"spider-go/process/app/tangcc/config"
	"spider-go/process/app/tangcc/persist/saverclient"
	"spider-go/process/model"
	"strings"
)

func Video() chan config.VideoItem {
	saveChan := make(chan config.VideoItem)
	client := saverclient.GetMysqlClient()
	go func() {
		for {
			video := <- saveChan
			saveVideo(client, video)
		}
	}()
	return saveChan
}

func saveVideo(client *xorm.Engine, videoItem config.VideoItem) {
	var videoModel model.Videos
	for _, tag := range videoItem.Tags {
		var tagModel model.Tags
		if has, err := client.Where(model.TagNameField + " = ?", tag).Exist(&tagModel); err == nil && !has {
			tagModel.Name = tag
			tagModel.Status = model.TagValidStatus
			if _, err = client.Insert(&tagModel); err != nil {
				logger.DefaultLogger.Error(err, nil)
			}
			videoModel.TagIds = fmt.Sprintf("%s,%d", videoModel.TagIds, tagModel.Id)
		}
	}
	videoModel.TagIds = strings.Trim(videoModel.TagIds, ",")

	var category model.Categories
	if has, err := client.Where(model.CategoryNameField + " = ?", videoItem.CategoryName).Get(&category); err == nil && !has {
		category.Name = videoItem.Name
		category.Status = model.CategoryValidStatus
		if _, err := client.Insert(&category); err != nil {
			logger.DefaultLogger.Error(err, nil)
		}
	}
	videoModel.CategoryId = int(category.Id)

	//var subject model.Subjects
	//if has, err := client.Where(model.SubjectNameFiled + " = ?", videoItem.SubjectName).Get(&subject); err == nil && !has {
	//	subject.Name = videoItem.Name
	//	subject.Status = model.SubjectValidStatus
	//	if _, err := client.Insert(&subject); err != nil {
	//		logger.DefaultLogger.Error(err, nil)
	//	}
	//}
	//videoModel.SubjectId = int(subject.Id)

	videoModel.Name = videoItem.Name
	videoModel.OriginalM3u8 = videoItem.M3u8
	videoModel.M3u8 = videoItem.M3u8
	videoModel.Status = model.VideoValidStatus
	videoModel.Mainimg = videoItem.MainImg
	videoModel.SeeNum = videoItem.SeeNum
	videoModel.LikeNum = videoItem.LikeNum
	videoModel.Long = videoItem.Long
	videoModel.CrawlerUrl = videoItem.CrawlerUrl

	if videoModel.Long > 60 * 60 {
		videoModel.LongType = model.VideoLongType
	} else if videoModel.Long > 25 * 60 && videoModel.Long < 60 * 60 {
		videoModel.LongType = model.VideoMiddleType
	} else {
		videoModel.LongType = model.VideoShortType
	}

	if _, err := client.Insert(&videoModel); err != nil {
		logger.DefaultLogger.Error(err, nil)
	}
	logger.DefaultLogger.Debug(fmt.Sprintf("saved item :%+v", videoModel), nil)
}


