package itemworker

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"path"
	"spider-go/config"
	"spider-go/fetcher"
	"spider-go/logger"
	xvideoConfig "spider-go/process/app/xvideo/aichange/config"
	"spider-go/process/app/xvideo/aichange/persist"
	"spider-go/process/model"
	"strconv"
	"strings"
	"util/filesystem"
)

//var (
	//childM3u8Re = regexp.MustCompile(`(hls-360p\.m3u8.*)`)
	//tsRe = regexp.MustCompile(`(.+ts.*)`)
//)

var saveChan = persist.Video()

var sliceWithSaveVideoChan = SliceWithSaveVideoPipeLine()

func Video() chan []config.Item {
	itemChan := make(chan []config.Item)
	go func() {
		for {
			items := <- itemChan
			doVideoWork(&items, saveChan)
		}
	}()
	return itemChan
}

func doVideoWork(items *[]config.Item, saveChan chan model.Videos)  {
	for _, item := range *items {
		//var video model.Videos
		//
		item := item.(xvideoConfig.Video)
		//
		//video.Name = item.Name
		//video.Status = model.VideoValidStatus

		dirPrefix, _ := uuid.NewV1()
		dir := fmt.Sprintf("%s/%s", xvideoConfig.BaseFileDir, dirPrefix)
		mainImg, err := downloadVideoMainImg(&item.MainImg, &dir)
		if err != nil {
			logger.DefaultLogger.Error(err, nil)
			return
		}
		item.MainImg = mainImg

		//video.Mainimg = mainImg

		//newM3u8, err := downloadVideo(&item.M3u8, &dir)

		//video.OriginalM3u8 = newM3u8
		//video.M3u8 = newM3u8

		newMp4, err := downloadVideo(&item.Mp4, &dir)
		if err != nil {
			logger.DefaultLogger.Error(err, nil)
			return
		}
		item.Mp4 = newMp4
		item.M3u8 = filesystem.Dir(newMp4) + "/" + xvideoConfig.M3u8File

		sliceWithSaveVideoChan <- item

		//saveChan <- video
	}
}

func downloadVideoMainImg(mainImg *string, baseDir *string) (string, error) {
	content, err := fetcher.File(*mainImg)
	if err != nil {
		return "", err
	}

	return filesystem.UploadFile(content, path.Base(*mainImg), fmt.Sprintf("%s/images", *baseDir))
}

func downloadVideo(mp4 *string, baseDir *string) (string, error) {
	content, err := fetcher.File(*mp4)
	if err != nil {
		return "", err
	}

	dir := fmt.Sprintf("%s/original", *baseDir)
	return filesystem.UploadFile(content, path.Base(strings.Split(*mp4, "?")[0]), dir)
}

func SliceWithSaveVideoPipeLine() chan xvideoConfig.Video {
	sliceVideoChan := make(chan xvideoConfig.Video)
	go func() {
		for {
			item := <- sliceVideoChan
			err := filesystem.SliceVideo(item.Mp4, item.M3u8, 5)
			if err != nil {
				logger.DefaultLogger.Error(err, nil)
			}

			var video model.Videos
			video.Name = item.Name
			video.Status = model.VideoValidStatus
			video.Mainimg = strings.TrimLeft(item.MainImg, "y:/")
			video.M3u8 = strings.TrimLeft(item.M3u8, "y:/")
			video.Mp4 = strings.TrimLeft(item.Mp4, "y:/")
			video.OriginalMp4 = video.Mp4
			long, _ := strconv.Atoi(strings.Replace(item.Long, " min", "", 1))
			video.Long = long * 60

			if video.Long > 60 * 60 {
				video.LongType = model.VideoLongType
			} else if video.Long > 25 * 60 && video.Long < 60 * 60 {
				video.LongType = model.VideoMiddleType
			} else {
				video.LongType = model.VideoShortType
			}

			video.CrawlerUrl = item.CrawlerUrl
			video.NeedLogin = model.VideoNeedLogin
			saveChan <- video
		}
	}()
	return sliceVideoChan
}

//func downloadVideo(m3u8File *string, baseDir *string) (string, error) {
//	content, err := fetcher.File(*m3u8File)
//	if err != nil {
//		return "", err
//	}
//
//
//	uriPrefix := url.GetUriPath(*m3u8File)
//	childM3u8 := uriPrefix + parseutil.ExtractString(childM3u8Re, content)
//	content, err = fetcher.File(childM3u8)
//	if err != nil {
//		return "", err
//	}
//
//	dir := fmt.Sprintf("%s/original", *baseDir)
//	newM3u8, err := filesystem.UploadFile(content, path.Base(strings.Split(childM3u8, "?")[0]), dir)
//	if err != nil {
//		return "", err
//	}
//
//	tsMatches := tsRe.FindAllSubmatch(content, -1)
//	tses := parseutil.ParseMatchesToString(tsMatches)
//
//	for _, ts := range tses {
//		url := uriPrefix + ts
//		logger.DefaultLogger.Info(url, nil)
//		all, err := fetcher.File(url)
//		if err != nil {
//			return "", err
//		}
//
//		_, err = filesystem.UploadFile(all, path.Base(strings.Split(ts, "?")[0]), dir)
//		if err != nil {
//			return "", err
//		}
//	}
//
//	return newM3u8, nil
//}