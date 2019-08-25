package itemworker

import (
	"bytes"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"spider-go/process/model"
	javdoveConfig "spider-go/process/app/javdove/config"
	javdovePersist "spider-go/process/app/javdove/persist"
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

type processVideo struct {
	m3u8File   string
	storageDir string
	model      model.Videos
}

func Videos(fetcherProcessNum int, videoProcessNum int) chan []config.Item {
 	itemChan := make(chan []config.Item)
 	saveChan := javdovePersist.Videos()
	videoProcessChan := make(chan processVideo)

	for i := 0; i < videoProcessNum; i++ {
		createVideoProcess(&videoProcessChan, &saveChan)
	}

	for i := 0; i < fetcherProcessNum; i++ {
		createFetcherProcess(&itemChan, &saveChan)
	}

 	return itemChan
}

func createFetcherProcess(itemChan *chan []config.Item, saveChan *chan model.Videos)  {
	go func() {
		for {
			items := <- *itemChan
			for _, item := range items {
				video := item.(javdoveConfig.VideoItem).Model

				dirPrefix, _ := uuid.NewV1()
				dir := fmt.Sprintf("%s/%s", javdoveConfig.BaseFileDir, dirPrefix)
				mainImg, err := downloadVideoMainImg(&video.Mainimg, &dir)
				if err != nil {
					logger.DefaultLogger.Error(err, nil)
					continue
				}
				video.Mainimg = strings.TrimLeft(mainImg, "y:/")

				m3u8File, err := downloadVideo(video.OriginalM3u8, dir)
				if err != nil {
					logger.DefaultLogger.Error(err, nil)
					continue
				}
				video.M3u8 = strings.TrimLeft(m3u8File, "y:/")

				category ,err := repository.GetCategory(item.(javdoveConfig.VideoItem).Category)
				if err != nil {
					logger.DefaultLogger.Error(err, nil)
					continue
				}
				video.CategoryId = int(category.Id)

				//processVideoChan <- processVideo{
				//	m3u8File: m3u8File,
				//	storageDir: dir,
				//	model: video,
				//}
				*saveChan <- video
			}
		}
	}()
}

func downloadVideoMainImg(mainImg *string, baseDir *string) (string, error) {
	content, err := fetcher.File(*mainImg)
	if err != nil {
		return "", err
	}

	return filesystem.UploadFile(content, path.Base(*mainImg), fmt.Sprintf("%s/images", *baseDir))
}

func downloadVideo(m3u8File string, baseDir string) (string, error) {
	content, err := fetcher.File(m3u8File)
	if err != nil {
		return "", err
	}

	dir := fmt.Sprintf("%s/original", baseDir)

	tsMatches := tsRe.FindAllSubmatch(content, -1)
	tses := parseutil.ParseMatchesToString(tsMatches)

	removeUri := url.GetUriPath(tses[0])
	content = bytes.ReplaceAll(content, []byte(removeUri), []byte(""))
	newM3u8, err := filesystem.UploadFile(content, javdoveConfig.NewM3u8Name, dir)
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

func createVideoProcess(processChan *chan processVideo, saveChan *chan model.Videos) {
	go func() {
		for {
			processVideo := <- *processChan

			dir := fmt.Sprintf("%s/new",  processVideo.storageDir)
			exist, err := filesystem.PathExists(dir)

			if err != nil || !exist {
				os.MkdirAll(dir, os.ModePerm)
			}

			outputMergeVideo, err := mergeVideo(&processVideo.m3u8File, &dir)
			if err != nil {
				logger.DefaultLogger.Error(err, nil)
				continue
			}

			videoFile, err := deleteVideoLogo(&outputMergeVideo, &dir)
			if err != nil {
				logger.DefaultLogger.Error(err, nil)
				continue
			}

			if err = os.RemoveAll(outputMergeVideo); err != nil {
				logger.DefaultLogger.Error(err, nil)
			}

			m3u8, err := sliceVideo(&videoFile, &dir)
			if err != nil {
				logger.DefaultLogger.Error(err, nil)
				continue
			}

			if err := os.RemoveAll(videoFile); err != nil {
				logger.DefaultLogger.Error(err, nil)
			}
			oldSliceDir := fmt.Sprintf("%s/original", processVideo.storageDir)
			removeOldSlice(&oldSliceDir)

			processVideo.model.M3u8 = strings.TrimLeft(m3u8, "y:/")

			*saveChan <- processVideo.model
		}
	}()
}

func mergeVideo(m3u8File *string, dir *string) (string, error) {
	content, err := ioutil.ReadFile(*m3u8File)
	if err != nil {
		return "", err
	}

	tsDir := filesystem.Dir(*m3u8File)
	tsMatches := tsRe.FindAllSubmatch(content, -1)
	tses := parseutil.ParseMatchesToString(tsMatches)

	var inputs []string
	for _, ts := range tses {
		input := tsDir + "/" + path.Base(ts)
		inputs = append(inputs, input)
	}

	output := fmt.Sprintf("%s/%s", strings.TrimRight(*dir, "/"), javdoveConfig.MergeVideoName)
	err = filesystem.MergeVideoByM3u8(inputs, output)
	if err != nil {
		return "", err
	}

	return output, nil
}

func deleteVideoLogo(videoFile *string, dir *string) (string, error) {
	exist, err := filesystem.PathExists(*dir)
	if !exist || err != nil {
		os.MkdirAll(*dir, os.ModePerm)
	}
	newVideoFile := strings.TrimRight(*dir, "/") + "/" + javdoveConfig.DeleteLogoVideoName
	err = filesystem.DeleteVideoLogo(*videoFile, newVideoFile, javdoveConfig.VideoLogoX, javdoveConfig.VideoLogoY, javdoveConfig.VideoLogoWidth, javdoveConfig.VideoLogoHeight)
	if err != nil {
		return "", err
	}
	return newVideoFile, nil
}

func sliceVideo(videoFile *string,dir *string) (string, error) {
	m3u8File := fmt.Sprintf("%s/%s", strings.TrimRight(*dir, "/"), "output.m3u8")
	err := filesystem.SliceVideo(*videoFile, m3u8File, 5)
	if err != nil {
		return "", err
	}
	return m3u8File, nil
}

func removeOldSlice(dir *string)  {
	files, _ := ioutil.ReadDir(*dir)
	for _, file := range files {
		if strings.Contains(file.Name(), ".ts")	{
			os.RemoveAll(file.Name())
		}
	}
}