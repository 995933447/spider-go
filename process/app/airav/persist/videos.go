package persist

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"spider-go/logger"
	"spider-go/process/app/airav/persist/saverclient"
	"spider-go/process/model"
)

func Videos() chan model.Videos {
	itemChan := make(chan model.Videos)
	client := saverclient.GetMysqlClient()
	go func() {
		for {
			item := <- itemChan
			saveVideos(client, &item)
		}
	}()
	return itemChan
}

func saveVideos(client *xorm.Engine, item *model.Videos) {
	go func(client *xorm.Engine, item *model.Videos) {
		//var oldVideo model.Videos
		//if exist, err := client.Where("crawler_url = ?", item.CrawlerUrl).Get(&oldVideo); err != nil {
		//		logger.DefaultLogger.Error(err, nil)
		//} else if !exist {
		var tag model.Tags
		if _, err := client.Where(model.TagNameField + " = ?", "HD").Get(&tag); err != nil {
			item.TagIds = fmt.Sprintf("%d", tag.Id)
		}

		_, err := client.Insert(item)
		if err != nil {
			logger.DefaultLogger.Error(err, nil)
		}

		var category model.Categories
		if _, err := client.Id(item.CategoryId).Get(&category); err == nil {
			category.VideoNum++
			_, err = client.Id(category.Id).Update(&category)
			if err != nil {
				logger.DefaultLogger.Error(err, nil)
			}
		} else {
			logger.DefaultLogger.Error(err, nil)
		}

		logger.DefaultLogger.Debug(fmt.Sprintf("save item %+v", *item), nil)
		//}
	}(client, item)
}