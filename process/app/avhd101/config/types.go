package config

import "spider-go/process/model"

type CategoryWithSubjectItem struct {
	Category model.Categories
	Subjects []model.Subjects
}

type VideoItem struct {
	Url string
	Name string
	Tags []string
	MainImg string
	M3u8File string
	Category string
	Subject string
	StoreDirPrefix string
}