package persist

import (
	"spider-go/config"
	"spider-go/process/app/xvideo/aichange/persist/saverclient"
	"spider-go/process/model"
)

var categories []model.Categories

func init()  {
	client := saverclient.GetMysqlClient()
	err := client.Where(model.CategoryStatusField + " = ?", model.CategoryValidStatus).Find(&categories)
	if err != nil {
		panic(err)
	}
}

func Categories() chan []config.Item {
	return config.NilItemWork()
}