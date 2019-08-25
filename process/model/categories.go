package model

import (
	"time"
)

const (
	CategoryValidStatus  = 1
	CategoryNameField  = "name"
	CategoryIdField = "id"
)

type Categories struct {
	Id int64
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	Name string
	Icon string
	Sort int
	Status int
	Remark []byte
	Description []byte
	VideoNum int
	PlayNum int
}

func (Categories) TableName() string {
	return "categories"
}