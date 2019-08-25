package persist

import (
	"fmt"
	"github.com/go-xorm/xorm"
	"spider-go/config"
	"spider-go/logger"
	avhd101Config "spider-go/process/app/avhd101/config"
	"spider-go/process/app/avhd101/persist/saverclient"
	"spider-go/process/model"
)

func CategoryWithSubjects() chan []config.Item {
	itemsChan := make(chan []config.Item)
	client := saverclient.GetMysqlClient()
	go func() {
		items := <- itemsChan
		saveCategoryWithSubject(&items, client)
	}()
	return itemsChan
}

func saveCategoryWithSubject(items *[]config.Item, client *xorm.Engine) {
	for _, item := range *items {
		item := item.(avhd101Config.CategoryWithSubjectItem)
		category := item.Category
		err := saveCategory(client, &category)
		if err != nil {
			logger.DefaultLogger.Error(err, nil)
			continue
		}

		if err := saveSubjects(client,&item.Subjects, &category); err != nil {
			logger.DefaultLogger.Error(err, nil)
		}

		logger.DefaultLogger.Debug(fmt.Sprintf("Save item %+v", item), nil)
	}
}

func saveCategory(client *xorm.Engine, category *model.Categories)  error  {
	var oldCategory model.Categories

	_, err := client.Where(model.CategoryNameField + " = ?", category.Name).Get(&oldCategory)
	if err != nil {
		return err
	} else if oldCategory.Id == 0 {
		if _, err = client.Insert(category); err != nil {
			return err
		}
	} else {
		*category = oldCategory
	}

	return nil
}

func saveSubjects(client *xorm.Engine, subjects *[]model.Subjects, category *model.Categories) error {
	for _, subject := range *subjects {
		var oldSubject model.Subjects
		if exists, err := client.Where(model.SubjectNameFiled + " = ? AND " + model.SubjectCategoryIdField + " = ?", subject.Name, category.Id).Exist(&oldSubject); err != nil || exists {
			if err != nil {
				return err
			}
			continue
		}

		subject.CategoryId = int(category.Id)
		if _, err := client.Insert(&subject); err != nil {
			return  err
		}
	}

	return nil
}

