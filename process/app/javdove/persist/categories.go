package persist

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"spider-go/process/model"
	"spider-go/logger"
	"spider-go/process/app/javdove/persist/saverclient"
)

func Categories() chan []model.Categories {
	itemsChan := make(chan []model.Categories)
	client := saverclient.GetMysqlClient()
	go func() {
		items := <- itemsChan
		saveCategories(client, &items)
	}()
	return itemsChan
}

func saveCategories(client *xorm.Engine, items *[]model.Categories) {
	for _, item := range *items {
		var oldCategory model.Categories
		exist, err := client.Where(model.CategoryNameField + " = ?", item.Name).Get(&oldCategory)
		if err != nil {
			logger.DefaultLogger.Error(err, nil)
			continue
		}

		if exist {
			logger.DefaultLogger.Info(fmt.Sprintf("Category existed %+v", oldCategory), nil)
			if oldCategory.Icon != "" {
				item.Icon = ""
			}
			item.VideoNum = oldCategory.VideoNum + item.VideoNum
			if _, err := client.Id(oldCategory.Id).Update(&item); err != nil {
				logger.DefaultLogger.Error(err, nil)
			}
			continue
		}

		item.Status = model.CategoryValidStatus
		if _, err := client.Insert(&item); err != nil {
			logger.DefaultLogger.Error(err, nil)
		} else {
			logger.DefaultLogger.Debug(fmt.Sprintf("save category: %+v\n", item), nil)
		}
	}
}