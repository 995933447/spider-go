package model

import "time"

const (
	TagValidStatus  = 1
)

const (
	TagIdField = "id"
	TagNameField  = "name"
)

type Tags struct {
	Id int64
	Name string
	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	Status int
	Sort int
}
