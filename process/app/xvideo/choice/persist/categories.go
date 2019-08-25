package persist

import (
	"spider-go/config"
	"spider-go/process/app/xvideo/aichange/persist/saverclient"
	"spider-go/process/model"
)

var category model.Categories

const categoryName = "网友推荐"

func init()  {
	client := saverclient.GetMysqlClient()
	has, err := client.Where(model.CategoryNameField + " = ?", categoryName).Get(&category)
	if err != nil {
		panic(err)
	}
	if !has {
		category.Name = categoryName
		category.Status = model.CategoryValidStatus
		_, err := client.Insert(&category)
		if err != nil {
			panic(err)
		}
	}
}

func Categories() chan []config.Item {
	return config.NilItemWork()
}