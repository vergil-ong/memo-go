package model

type NoticeQo struct {
	Title string `json:"title" form:"title"`
	Desc  string `json:"desc" form:"desc"`
}
