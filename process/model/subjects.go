package model

import "time"

const (
	SubjectValidStatus  = 1
	SubjectNameFiled  = "name"
	SubjectCategoryIdField  = "category_id"
)

type Subjects struct {
	Id int64
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	Name string
	Description []byte
	Remark []byte
	Sort int
	Status int
	CategoryId int
}

