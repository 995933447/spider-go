package persist

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"spider-go/process/model"
	"spider-go/logger"
	"spider-go/process/app/xvideo/aichange/persist/saverclient"
)

func Video() chan model.Videos {
	saveChan := make(chan model.Videos)
	client := saverclient.GetMysqlClient()
	go func() {
		for {
			video := <- saveChan
			saveVideo(client, &video)
		}
	}()
	return saveChan
}

func saveVideo(client *xorm.Engine,video *model.Videos) {
	video.CategoryId = int(category.Id)
	_, err := client.Insert(video)
	if err != nil {
		logger.DefaultLogger.Error(err, nil)
	}
	category.VideoNum++
	_, err = client.Id(category.Id).Update(&category)
	if err != nil {
		logger.DefaultLogger.Error(err, nil)
	}
	logger.DefaultLogger.Debug(fmt.Sprintf("save item:%+v", video), nil)
}