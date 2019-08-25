package model

import "time"

const(
	VideoValidStatus  = 1
	VideoLongType  = 2
	VideoShortType  = 0
	VideoMiddleType  =  1
	VideoNeedCharge  = 1
	VideoNeedntChange  = 0
	VideoNeedLogin = 1
	VideoNeedntLogin = 0
)

const VideoCrawlerUrlField = "crawler_url"

type Videos struct {
	Id int64
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	Name string
	Icon string
	CategoryId int
	SubjectId  int
	Mainimg string
	OriginalM3u8 string
	CrawlerUrl string
	M3u8 string
	OriginalMp4 string
	Mp4 string
	Long int
	TagIds string
	NeedCharge int
	SeeNum int
	RealSeeNum int
	CollectNum int
	RealCollectNum int
	Status int
	LongType int
	PreviewVideo string
	LikeNum int
	DislikeNum int
	NeedLogin int
}
