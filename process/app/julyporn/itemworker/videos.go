package itemworker

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"path"
	"regexp"
	"spider-go/process/model"
	julypornConfig "spider-go/process/app/julyporn/config"
	julypornPersist "spider-go/process/app/julyporn/persist"
	"spider-go/process/repository"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/logger"
	"spider-go/parseutil"
	"util/filesystem"
	"util/url"
)

var tsRe = regexp.MustCompile(`.+ts`)

func Videos() chan []config.Item {
	itemsChan := make(chan []config.Item)
	saveChan := julypornPersist.Video()
	go func() {
		items := <- itemsChan
		doVideoWork(&items, &saveChan)
	}()
	return itemsChan
}

func doVideoWork(items *[]config.Item, saveChan *chan model.Videos)  {
	for _, item := range *items {
		go func(item config.Item) {
			video := item.(julypornConfig.VideoItem).Model

			dirPrefix, _ := uuid.NewV1()
			dir := fmt.Sprintf("%s/%s", julypornConfig.BaseFileDir, dirPrefix)
			mainImg, err := downloadVideoMainImg(&video.Mainimg, &dir)
			if err != nil {
				logger.DefaultLogger.Error(err, nil)
				return
			}
			video.Mainimg = mainImg

			videoFile, err := downloadVideo(&video.OriginalM3u8, &dir)
			if err != nil {
				logger.DefaultLogger.Error(err, nil)
				return
			}
			video.M3u8 = videoFile

			category ,err := repository.GetCategory(item.(julypornConfig.VideoItem).Category)
			if err != nil {
				logger.DefaultLogger.Error(err, nil)
				return
			}

			video.CategoryId = int(category.Id)

			*saveChan <- video
		}(item)
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

	dir := fmt.Sprintf("%s/original", *baseDir)

	newM3u8, err := filesystem.UploadFile(content, path.Base(*m3u8File), dir)
	if err != nil {
		return "", err
	}

	tsMatches := tsRe.FindAllSubmatch(content, -1)
	tses := parseutil.ParseMatchesToString(tsMatches)
	uriPrefix := url.GetUriPath(*m3u8File)

	for _, ts := range tses {
		url := uriPrefix + ts
		all, err := fetcher.File(url)
		if err != nil {
			return "", err
		}

		_, err = filesystem.UploadFile(all, ts, dir)
		if err != nil {
			return "", err
		}
	}

	return newM3u8, nil
}