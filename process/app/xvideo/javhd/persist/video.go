package persist

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"spider-go/logger"
	"spider-go/process/app/xvideo/aichange/persist/saverclient"
	"spider-go/process/model"
	"util/rand"
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
	video.CategoryId = int(categories[rand.RandIntN(len(categories) - 1)].Id)
	_, err := client.Insert(video)
	if err != nil {
		logger.DefaultLogger.Error(err, nil)
	} else {
		logger.DefaultLogger.Debug(fmt.Sprintf("save item:%+v", video), nil)
	}

}