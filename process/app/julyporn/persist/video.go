package persist

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"spider-go/logger"
	"spider-go/process/app/julyporn/persist/saverclient"
	"spider-go/process/model"
)

func Video() chan model.Videos {
	saveChan := make(chan model.Videos)
	client := saverclient.GetMysqlClient()
	go func() {
		saveItem := <- saveChan
		saveVideo(client, &saveItem)
	}()
	return saveChan
}

func saveVideo(client *xorm.Engine,item *model.Videos)  {
	//var oldVideo model.Videos
	//if exist, err := client.Where("crawler_url = ?", item.CrawlerUrl).Get(&oldVideo); err != nil {
	//	logger.DefaultLogger.Error(err, nil)
	//} else if !exist {
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
}
