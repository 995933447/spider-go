package itemworker

import (
	"bytes"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"path"
	"regexp"
	"spider-go/process/model"
	airavConfig "spider-go/process/app/airav/config"
	"spider-go/process/app/airav/persist"
	"spider-go/process/repository"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/logger"
	"spider-go/parseutil"
	"strings"
	"util/filesystem"
	"util/url"
)

var tsRe = regexp.MustCompile(`.+ts`)

func Videos() chan []config.Item {
	itemChan := make(chan []config.Item)
	saveChan := persist.Videos()

	go func() {
		items := <- itemChan
		doVideoWork(&items, &saveChan)
	}()

	return itemChan
}

func doVideoWork(items *[]config.Item, saveChan *chan model.Videos)  {
	for _, item := range *items {
		go func(originalItem *config.Item, saveChan *chan model.Videos) {
			item := (*originalItem).(airavConfig.VideoItem)
			var video model.Videos

			dirPrefix, _ := uuid.NewV1()
			dir := fmt.Sprintf("%s/%s", airavConfig.BaseFileDir, dirPrefix)
			mainImg, err := downloadVideoMainImg(&video.Mainimg, &dir)
			if err != nil {
				logger.DefaultLogger.Error(err, nil)
				return
			}
			video.Mainimg = strings.TrimLeft(mainImg, "d:/")

			m3u8File, err := downloadVideo(&video.OriginalM3u8, &dir)
			if err != nil {
				logger.DefaultLogger.Error(err, nil)
				return
			}
			video.M3u8 = strings.TrimLeft(m3u8File, "d:/")

			category ,err := repository.GetCategory(item.CategoryName)
			if err != nil {
				logger.DefaultLogger.Error(err, nil)
				return
			}
			video.CategoryId = int(category.Id)

			*saveChan <- video
		}(&item, saveChan)
	}
}

func downloadVideoMainImg(mainImg *string, baseDir *string) (string, error) {
	content, err := fetcher.File(*mainImg)
	if err != nil {
		return "", err
	}

	return filesystem.UploadFile(content, path.Base(*mainImg), fmt.Sprintf("%s/images", *baseDir))
}

func downloadVideo(m3u8File *string, baseDir *string) (string, error) {
	content, err := fetcher.File(*m3u8File)
	if err != nil {
		return "", err
	}

	dir := fmt.Sprintf("%s/original", *baseDir)

	tsMatches := tsRe.FindAllSubmatch(content, -1)
	tses := parseutil.ParseMatchesToString(tsMatches)

	removeUri := url.GetUriPath(tses[0])
	content = bytes.ReplaceAll(content, []byte(removeUri), []byte(""))
	newM3u8, err := filesystem.UploadFile(content, airavConfig.NewM3u8Name, dir)
	if err != nil {
		return "", err
	}

	for _, ts := range tses {
		all, err := fetcher.File(ts)
		if err != nil {
			logger.DefaultLogger.Error(err, nil)
			return "", err
		}

		_, err = filesystem.UploadFile(all, path.Base(ts), dir)
		if err != nil {
			return "", err
		}
	}

	return newM3u8, nil
}

