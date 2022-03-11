package model

type MemoVo struct {
	Id         uint   `json:"id"`
	Title      string `json:"title"`
	DescShow   string `json:"descShow"`
	NoticeTime int64  `json:"noticeTime"`
}
