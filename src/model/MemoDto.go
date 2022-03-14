package model

import (
	"strconv"
	"time"
)

type TestTable struct {
	Id       uint      `json:"id"`
	TestName string    `json:"TestName"`
	TestDate time.Time `json:"TestDate"`
}

type User struct {
	Id     uint
	OpenId string
}

type Memo struct {
	Id         uint   `json:"id" gorm:"AUTO_INCREMENT"`
	Title      string `json:"title"`
	DescShow   string `json:"descShow"`
	TaskId     uint
	UserId     uint
	NoticeTime time.Time `json:"noticeTime"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

type MemoTask struct {
	Id                          uint   `json:"id" gorm:"AUTO_INCREMENT"`
	Title                       string `json:"title"`
	Desc                        string `json:"desc"`
	UserId                      uint
	NoticeType                  int
	NoticeOnceTime              time.Time
	NoticePeriodFirstTime       time.Time
	NoticePeriodFirstTimeMinute int
	NoticePeriodTimes           int
	NoticePeriodInterval        int
	CreateTime                  time.Time
	NoticeTaskStatus            int
}

const NoticeTypeOnce = 1
const NOTICE_TYPE_CIRCLE = 2

const NoticeTaskStatusCreate = 1
const NoticeTaskStatusDone = 2

func MemoTaskParseNoticeType(noticeTypeStr string) int {
	if "一次性提醒" == noticeTypeStr {
		return 1
	} else if "周期性提醒" == noticeTypeStr {
		return 2
	}
	return 0
}

func MemoTaskParseNoticeOnceTime(noticeOnceCal string, noticeOnceTime string) time.Time {
	timeStr := noticeOnceCal + " " + noticeOnceTime + ":00"
	time, _ := time.ParseInLocation("2006/1/02 15:04:05", timeStr, time.Local)
	return time
}

func MemoTaskParseNoticePeriodFirstTime(noticeOnceTime time.Time, beforeMinutes int) time.Time {
	duration, _ := time.ParseDuration("-" + strconv.Itoa(beforeMinutes) + "m")
	return noticeOnceTime.Add(duration)
}

const MemoNoticeStatusRecord = 1
const MemoNoticeStatusSend = 2

type MemoNotice struct {
	Id           uint
	NoticeTime   time.Time
	Title        string
	DescShow     string
	TaskId       uint
	CreateTime   time.Time
	NoticeStatus int
}
