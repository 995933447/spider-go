package config

import "spider-go/process/model"

type(
	CategoryItem struct {
		Icon string
		Name string
		VideoNum int
	}

	VideoItem struct {
		Category string
		Model    model.Videos
	}
)