package itemworker

import (
	"fmt"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"os"
	"path"
	"regexp"
	"spider-go/config"
	"spider-go/logger"
	"spider-go/parseutil"
	tangccConfig "spider-go/process/app/tangcc/config"
	"spider-go/process/app/tangcc/fetcher"
	"spider-go/process/app/tangcc/persist"
	"strconv"
	"strings"
	"time"
	"util/filesystem"
	"util/url"
)

var (
	childM3u8Re = regexp.MustCompile(`(.+m3u8)`)
	tsRe = regexp.MustCompile(`(.+ts.*)`)
	longRe = regexp.MustCompile(`#EXTINF:([^,]+),`)
	keyRe = regexp.MustCompile(`#EXT-X-KEY:METHOD=AES-128,URI="([^"]+)"`)
)

var saveChan = persist.Video()

func Videos() chan []config.Item {
	itemChan := make(chan []config.Item)
	go func() {
		items := <- itemChan
		for _, item := range items {
			item := item.(tangccConfig.VideoItem)
			go doVideoWork(item)
		}
	}()
	return itemChan
}

func doVideoWork(item tangccConfig.VideoItem) {
	dirPrefix, _ := uuid.NewV1()
	dir := fmt.Sprintf("%s/%s", tangccConfig.BaseFileDir, dirPrefix)
	if err := downloadVideo(&item, &dir); err != nil {
		logger.DefaultLogger.Error(err, nil)
		return
	}
	if err := downloadVideoMainImg(&item, &dir); err != nil {
		logger.DefaultLogger.Error(err, nil)
		return
	}

	item.M3u8 = strings.TrimLeft(item.M3u8, "y:/")
	item.MainImg = strings.TrimLeft(item.MainImg, "y:/")

	saveChan <- item
}

func downloadVideo(item *tangccConfig.VideoItem, baseDir *string) error {
	// 抓取M3U8内容
	m3u8Content, err := fetcher.File(item.M3u8)
	if err != nil {
		return err
	}

	// 抓取子M3U8
	uriPrefix := url.GetUriPath(item.M3u8)
	childM3u8 := parseutil.ExtractString(childM3u8Re, m3u8Content)
	if childM3u8 != "" {
		host, _ := url.GetHost(item.M3u8)
		scheme, _ := url.GetScheme(item.M3u8)
		uriPrefix = fmt.Sprintf("%s://%s/", scheme, strings.TrimLeft(host, "/"))
		item.M3u8 = uriPrefix + strings.TrimLeft(childM3u8, "/")
		if m3u8Content, err = fetcher.File(item.M3u8); err != nil {
			return err
		}
	}

	// 获取视频时长
	segmentLongMathces := longRe.FindAllSubmatch(m3u8Content, -1)
	segmentLongs := parseutil.ParseMatchesToString(segmentLongMathces)
	var long float64
	for _, segmentLong := range segmentLongs {
		segLong, _ := strconv.ParseFloat(segmentLong, 64)
		long += segLong
	}
	item.Long = int(long)

	// 获取ts路径
	tsMatches := tsRe.FindAllSubmatch(m3u8Content, -1)
	tses := parseutil.ParseMatchesToString(tsMatches)
	if len(tses) <= 0 {
		return errors.New("Unknow error")
	}
	//创建ts和M3U8存放目录
	m3u8Dir := fmt.Sprintf("%s/original", *baseDir)
	tsDir := fmt.Sprintf("%s/%s", m3u8Dir, strings.TrimRight(path.Dir(tses[0]), "."))
	if has, err := filesystem.PathExists(tsDir); err != nil || !has {
		if err := os.MkdirAll(tsDir, os.ModePerm); err != nil {
			logger.DefaultLogger.Error(err, nil)
		}
	}

	// 获取并下载密钥
	key := parseutil.ExtractString(keyRe, m3u8Content)
	if key != "" {
		keyContent, err := fetcher.File(fmt.Sprintf("%s%s", uriPrefix, strings.TrimLeft(key, "/")))
		if err != nil {
			return err
		}
		keyDir := fmt.Sprintf("%s/%s", m3u8Dir, strings.TrimRight(path.Dir(key), "."))
		if has, err := filesystem.PathExists(keyDir); err != nil || !has {
			if err := os.MkdirAll(keyDir, os.ModePerm); err != nil {
				logger.DefaultLogger.Error(err, nil)
			}
		}

		if _, err =filesystem.UploadFile(keyContent, path.Base(key), keyDir); err != nil {
			return err
		}
	}

	// 下载M3U8
	if item.M3u8, err = filesystem.UploadFile(m3u8Content, strings.Split(path.Base(item.M3u8), "?")[0], m3u8Dir); err != nil {
		return err
	}

	//下载TS
	var all []byte
	for _, ts := range tses {
		all, err = fetcher.File(fmt.Sprintf("%s%s", uriPrefix, strings.TrimLeft(ts, "/")))
		try := 0
		for try <= 5 {
			if all, err = fetcher.File(fmt.Sprintf("%s%s", uriPrefix, strings.TrimLeft(ts, "/"))); err != nil {
				time.Sleep(time.Second * 3)
				try++
			} else {
				break
			}
		}
		if err != nil {
			return err
		}

		//if err != nil {
		//	time.Sleep(time.Second * time.Duration(rand.RandIntN(5)))
		//	all, err = fetcher.File(fmt.Sprintf("%s%s", uriPrefix, strings.TrimLeft(ts, "/")))
		//	time.Sleep(time.Second * time.Duration(rand.RandIntN(5)))
		//	if err != nil {
		//		all, err = fetcher.File(fmt.Sprintf("%s%s", uriPrefix, strings.TrimLeft(ts, "/")))
		//		time.Sleep(time.Second * time.Duration(rand.RandIntN(5)))
		//		if err != nil {
		//			return err
		//		}
		//	}
		//
		//}

		_, err = filesystem.UploadFile(all, path.Base(ts), tsDir)
		if err != nil {
			return fmt.Errorf("fetching ts err:%s", err.Error())
		}
	}

	return nil
}

func downloadVideoMainImg(item *tangccConfig.VideoItem, baseDir *string) error {
	content, err := fetcher.File(item.MainImg)
	if err != nil {
		return err
	}

	if mainImg, err := filesystem.UploadFile(content, path.Base(item.MainImg), fmt.Sprintf("%s/images", *baseDir)); err != nil {
		return err
	} else {
		item.MainImg = mainImg
	}

	return  nil
}
