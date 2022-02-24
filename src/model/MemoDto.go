package model

import (
	"time"
)

type TestTable struct {
	Id       uint      `json:"id"`
	TestName string    `json:"TestName"`
	TestDate time.Time `json:"TestDate"`
}

type Memo struct {
	Id         uint      `json:"id" gorm:"AUTO_INCREMENT"`
	Title      string    `json:"title"`
	DescShow   string    `json:"descShow"`
	NoticeTime time.Time `json:"noticeTime"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}
