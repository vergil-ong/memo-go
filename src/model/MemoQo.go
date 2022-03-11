package model

type NoticeQo struct {
	Title                 string `json:"title" form:"title"`
	Desc                  string `json:"desc" form:"desc"`
	NoticeType            string `json:"noticeType" form:"noticeType"`
	NoticeOnceCal         string `json:"noticeOnceCal" form:"noticeOnceCal"`
	NoticeOnceTime        string `json:"noticeOnceTime" form:"noticeOnceTime"`
	NoticePeriodFirstTime int    `json:"noticePeriodFirstTime" form:"noticePeriodFirstTime"`
	NoticePeriodTimes     int    `json:"noticePeriodTimes" form:"noticePeriodTimes"`
	NoticePeriodInterval  int    `json:"noticePeriodInterval" form:"noticePeriodInterval"`
}
