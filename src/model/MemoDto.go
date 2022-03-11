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

type Memo struct {
	Id         uint   `json:"id" gorm:"AUTO_INCREMENT"`
	Title      string `json:"title"`
	DescShow   string `json:"descShow"`
	TaskId     uint
	NoticeTime time.Time `json:"noticeTime"`
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

type MemoTask struct {
	Id                          uint   `json:"id" gorm:"AUTO_INCREMENT"`
	Title                       string `json:"title"`
	Desc                        string `json:"desc"`
	NoticeType                  int
	NoticeOnceTime              time.Time
	NoticePeriodFirstTime       time.Time
	NoticePeriodFirstTimeMinute int
	NoticePeriodTimes           int
	NoticePeriodInterval        int
	CreateTime                  time.Time
}

const NOTICE_TYPE_ONCE = 1
const NOTICE_TYPE_CIRCLE = 2

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
