package persist

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"spider-go/config"
	"spider-go/logger"
	"spider-go/process/model"
	"spider-go/process/app/tangcc/persist/saverclient"
)

func Categories() chan []config.Item {
	saveChan := make(chan []config.Item)
	client := saverclient.GetMysqlClient()
	go func() {
		for {
			items := <- saveChan
			saveCategories(client, &items)
		}
	}()
	return saveChan
}

func saveCategories(client *xorm.Engine,categories *[]config.Item) {
	for _, category := range *categories {
		category := category.(model.Categories)
		has, err := client.Where(model.CategoryNameField + " = ?", category.Name).Get(&category)
		if err != nil {
			logger.DefaultLogger.Error(err, nil)
			continue
		}
		if !has {
			category.Status = model.CategoryValidStatus
			_, err := client.Insert(category)
			if err != nil {
				logger.DefaultLogger.Error(err, nil)
			}
			logger.DefaultLogger.Debug(fmt.Sprintf("saved category %+v", category), nil)
		}
	}
}
