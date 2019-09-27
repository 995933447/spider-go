package distinctor

import (
	"fmt"
	"spider-go/process/model"
	"testing"
)

func TestVideo_CheckIsFetched(t *testing.T) {
	var videos  []model.Videos
	err := mysqlClient.Cols(model.VideoCrawlerUrlField).Find(&videos)
	if err != nil {
		panic(err)
	}

	for _, video := range videos {
		fetched.Store(video.CrawlerUrl, true)
	}

	url := "http://www.tangchaoav88.com/show-11-4188-1.html"
	_, has := fetched.Load(url)
	fmt.Println(has)
}
