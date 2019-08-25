package persist

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"spider-go/config"
	"spider-go/logger"
	"spider-go/process/model"
	"spider-go/process/app/tangcc/persist/saverclient"
)

func Subjects() chan []config.Item {
	saveChan := make(chan []config.Item)
	client := saverclient.GetMysqlClient()
	go func() {
		for {
			items := <- saveChan
			saveSubjects(client, &items)
		}
	}()
	return saveChan
}

func saveSubjects(client *xorm.Engine, subjects *[]config.Item) {
	for _, category := range *subjects {
		subject := category.(model.Subjects)
		has, err := client.Where(model.SubjectNameFiled + " = ?", subject.Name).Get(&subject)
		if err != nil {
			logger.DefaultLogger.Error(err, nil)
			continue
		}
		if !has {
			subject.Status = model.SubjectValidStatus
			_, err := client.Insert(subject)
			if err != nil {
				logger.DefaultLogger.Error(err, nil)
			}
			logger.DefaultLogger.Debug(fmt.Sprintf("saved subject %+v", subject), nil)
		}
	}
}