package distinctor

import (
	"github.com/go-xorm/xorm"
	"spider-go/process/model"
	"spider-go/logger/mysql"
	mysqlSaver "spider-go/persist/mysql"
	"sync"
)

type Video struct {
}

var (
	mysqlClient *xorm.Engine
	fetched = sync.Map{}
)

func init()  {
	var err error
	mysqlClient, err = mysqlSaver.NewClient()
	if err != nil {
		panic(err)
	}

	if err = mysql.ListenSql(mysqlClient); err != nil {
		panic(err)
	}

	var videos  []model.Videos
	err = mysqlClient.Cols(model.VideoCrawlerUrlField).Find(&videos)
	if err != nil {
		panic(err)
	}

	for _, video := range videos {
		fetched.Store(video.CrawlerUrl, true)
	}
}

func (*Video) RecordFetched(url string) {
	fetched.Store(url, true)
}

func (*Video) CheckIsFetched(url string) (bool, error) {
	_, has := fetched.Load(url)
	return has, nil
}