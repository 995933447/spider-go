package persist

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"spider-go/logger"
	"spider-go/process/app/javdove/persist/saverclient"
	xiaoniaoConfig "spider-go/process/app/xiaoniao/config"
	"spider-go/process/model"
)

func Videos() chan xiaoniaoConfig.VideoItem {
	itemChan := make(chan xiaoniaoConfig.VideoItem)
	client := saverclient.GetMysqlClient()
	go func() {
		for {
			item := <- itemChan
			saveVideos(client, &item)
		}
	}()
	return itemChan
}

func saveVideos(client *xorm.Engine, item *xiaoniaoConfig.VideoItem) {
	go func(item *xiaoniaoConfig.VideoItem) {
		var videoModel model.Videos
		videoModel.Name = item.Name
		videoModel.M3u8 = item.M3u8
		videoModel.Mainimg = item.Mainimg
		videoModel.CrawlerUrl = item.CrawlerUrl
		videoModel.OriginalM3u8 = item.M3u8
		videoModel.LikeNum = item.Like
		videoModel.Status = model.VideoValidStatus
		videoModel.Long = item.Long

		var categoryModel model.Categories
		_, err := client.Where(model.CategoryNameField + " = ?", item.Category).Get(&categoryModel)
		if err != nil {
			logger.DefaultLogger.Error(err, nil)
			return
		}

		videoModel.CategoryId = int(categoryModel.Id)
		videoModel.Status = model.VideoValidStatus
		_, err = client.Insert(&videoModel)
		if err != nil {
			logger.DefaultLogger.Error(err, nil)
		}
		logger.DefaultLogger.Debug(fmt.Sprintf("save item %+v", *item), nil)
	}(item)
}