package itemworker

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"path"
	"regexp"
	avhd101Config "spider-go/process/app/avhd101/config"
	avhd101Persist "spider-go/process/app/avhd101/persist"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/logger"
	"spider-go/parseutil"
	"strings"
	"util/filesystem"
	"util/url"
)

var (
	baseDir  = strings.TrimRight(avhd101Config.BaseFileDir, "/")
	tsRe = regexp.MustCompile(`.+ts`)
)

func Videos() chan []config.Item {
	itemChan := make(chan []config.Item)
	saveItemChan := avhd101Persist.Video()
	go func() {
		item := <- itemChan
		doVideoWork(&item, &saveItemChan)
	}()
	return itemChan
}

func doVideoWork(items *[]config.Item, saveChan *chan avhd101Config.VideoItem)  {
	for _, item := range *items {
		go func() {
			item := item.(avhd101Config.VideoItem)
			err := downVideoMainImg(&item.MainImg, &item)
			if err != nil {
				logger.DefaultLogger.Error(err, nil)
				return
			}
			err = downVideo(&item.M3u8File, &item)
			if err != nil {
				logger.DefaultLogger.Error(err, nil)
				return
			}

			*saveChan <- item
		}()
	}
}

func downVideoMainImg(mainImg *string, item *avhd101Config.VideoItem) error {
	content, err := fetcher.File(*mainImg)

	dirPrefix, _ := uuid.NewV1()
	dir := fmt.Sprintf("%s/%s/images", baseDir, dirPrefix)
	newImg, err := filesystem.UploadFile(content, path.Base(*mainImg), dir)
	if err != nil {
		return err
	}

	item.MainImg = newImg
	item.StoreDirPrefix = fmt.Sprintf("%s/%s", baseDir, dirPrefix)

	return nil
}

func downVideo(m3u8File *string, item *avhd101Config.VideoItem) error {
	content, err := fetcher.File(*m3u8File)

	dir := fmt.Sprintf("%s/original", item.StoreDirPrefix)
	newM3u8, err := filesystem.UploadFile(content, path.Base(item.M3u8File), dir)
	if err != nil {
		return err
	}

	tsMatches := tsRe.FindAllSubmatch(content, -1)
	tses := parseutil.ParseMatchesToString(tsMatches)
	uriPrefix := url.GetUriPath(item.M3u8File)

	for _, ts := range tses {
		url := uriPrefix + ts
		all, err := fetcher.File(url)
		if err != nil {
			return err
		}

		_, err = filesystem.UploadFile(all, ts, dir)
		if err != nil {
			return err
		}
	}

	item.M3u8File = newM3u8
	return nil
}