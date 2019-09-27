package persist

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"spider-go/config"
	"spider-go/process/model"
	"spider-go/logger"
	"spider-go/process/app/javdove/persist/saverclient"
)

func Categories() chan []config.Item {
	itemsChan := make(chan []config.Item)
	client := saverclient.GetMysqlClient()
	go func() {
		items := <- itemsChan
		saveCategories(client, &items)
	}()
	return itemsChan
}

func saveCategories(client *xorm.Engine, items *[]config.Item) {
	for _, item := range *items {
		item := item.(model.Categories)
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