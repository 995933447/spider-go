package repository

import (
	"spider-go/process/model"
	"spider-go/persist/mysql"
	"strings"
)

func GetCategory(name string) (model.Categories, error) {
	var category model.Categories
	client, err := mysql.NewClient()
	if err != nil {
		return category, err
	}

	_, err = client.Where(model.CategoryNameField + " = ?", strings.Trim(name, " ")).Get(&category)
	if err != nil {
		return category, err
	}
	return category, nil
}