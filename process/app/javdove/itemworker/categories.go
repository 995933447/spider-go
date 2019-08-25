package itemworker

import (
	"path"
	"spider-go/process/model"
	javdoveConfig "spider-go/process/app/javdove/config"
	javdovePersist "spider-go/process/app/javdove/persist"
	"spider-go/config"
	"spider-go/fetcher"
	"util/filesystem"
)

func Categories() chan []config.Item {
	itemsChan := make(chan []config.Item)
	saverChan := javdovePersist.Categories()
	var saverItems []model.Categories
	go func() {
		items := <- itemsChan

		for _, item := range items {
			item := item.(javdoveConfig.CategoryItem)
			filename := path.Base(item.Icon)
			if exists, err := filesystem.PathExists(javdoveConfig.CategoryImgDir + "/" + filename); err != nil || !exists  {
				content, err := fetcher.File(item.Icon)
				if err != nil {
					item.Icon = ""
				}

				newImg, err := filesystem.UploadFile(content, filename, javdoveConfig.CategoryImgDir)
				if err != nil {
					item.Icon = ""
				} else {
					item.Icon = newImg


				}
			} else {
				item.Icon = javdoveConfig.CategoryImgDir + "/" + filename
			}
			saverItems = append(saverItems, model.Categories{
				Name: item.Name,
				Icon: item.Icon,
				VideoNum: item.VideoNum,
			})
		}

		saverChan <- saverItems
	}()
	return itemsChan
}
