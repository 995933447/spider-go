package persist

import (
	"github.com/go-xorm/xorm"
	"spider-go/logger"
	"spider-go/process/app/avhd101/config"
	"spider-go/process/app/avhd101/persist/saverclient"
	"spider-go/process/model"
	"strconv"
	"strings"
	"util/rand"
)

func Video() chan config.VideoItem {
	itemsChan := make(chan config.VideoItem)
	client := saverclient.GetMysqlClient()
	go func() {
		item := <- itemsChan
		saveVideo(client, &item)
	}()
	return itemsChan
}

func saveVideo(client *xorm.Engine, item *config.VideoItem)  {
		var video model.Videos
		exist, err := client.Where("crawler_url = ?", item.Url).Exist(&video)
		if err != nil ||exist {
			return
		}

		video.TagIds = getTagIds(client, item)

		if categoryId, err := getCategoryId(client,item); err == nil {
			video.CategoryId = categoryId
		} else {
			logger.DefaultLogger.Error(err, nil)
		}

		if subjectId, err := getSubjectId(client, item); err == nil {
			video.SubjectId = subjectId
		} else {
			logger.DefaultLogger.Error(err, nil)
		}

		video.Name = item.Name
		video.Status = model.VideoValidStatus
		video.Mainimg = item.MainImg
		video.M3u8 = item.M3u8File
		video.OriginalM3u8 = item.M3u8File
		video.LongType = model.VideoShortType
		video.CrawlerUrl = item.Url
		video.NeedCharge = model.VideoNeedntChange
		video.CollectNum = rand.RandIntN(100)
		video.SeeNum = rand.RandIntN(1000)
		if _, err := client.Insert(&video); err != nil {
			logger.DefaultLogger.Error(err, nil)
		}
}

func getTagIds(client *xorm.Engine,item *config.VideoItem) string {
	var tags []model.Tags
	tagIds := ""
	for _, tag := range item.Tags {
		tagModel := model.Tags{
			Name:   tag,
			Status: model.TagValidStatus,
		}
		tags = append(tags, tagModel)
	}

	var distinctTags []model.Tags
	for _, tag := range tags {
		exist, err := client.Where(model.TagNameField + " = ?", tag.Name).Exist(&tag)
		if err == nil && !exist {
			distinctTags = append(distinctTags, tag)
		}
	}

	if _, err := client.Insert(&distinctTags); err != nil {
		var tagIdArr []string
		for _, tagModel := range distinctTags {
			tagIdArr = append(tagIdArr, strconv.Itoa(int(tagModel.Id)))
			tagIds = strings.Join(tagIdArr, ",")
		}
	}
	return tagIds
}

func getCategoryId(client *xorm.Engine, item *config.VideoItem)(int, error) {
	var category model.Categories
	_, err := client.Where(model.CategoryNameField + " = ?", strings.Trim(item.Category, " ")).Get(&category)
	if err != nil {
		logger.DefaultLogger.Error(err, nil)
		return -1, err
	}
	category.VideoNum++
	_, err = client.Id(category.Id).Update(&category)
	if err != nil {
		logger.DefaultLogger.Error(err, nil)
		return -1, err
	}
	return int(category.Id), nil
}

func getSubjectId(client *xorm.Engine, item *config.VideoItem) (int, error) {
	var subject model.Subjects
	_, err := client.Where(model.SubjectNameFiled + " = ?", strings.Trim(item.Subject, " ")).Get(&subject)
	if err != nil {
		return -1, err
	}
	return int(subject.Id), nil
}