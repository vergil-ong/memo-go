package service

import (
	"MemoProjects/src/config"
	"MemoProjects/src/logger"
	"MemoProjects/src/model"
	"time"
)

func AddNoticeTask(qo model.NoticeQo) model.MemoTask {
	conn := config.GetConn()

	task := model.MemoTask{}

	noticeType := model.MemoTaskParseNoticeType(qo.NoticeType)
	if noticeType == model.NOTICE_TYPE_ONCE {
		onceTime := model.MemoTaskParseNoticeOnceTime(qo.NoticeOnceCal, qo.NoticeOnceTime)
		task = model.MemoTask{
			Title:                       qo.Title,
			Desc:                        qo.Desc,
			NoticeType:                  noticeType,
			NoticeOnceTime:              onceTime,
			NoticePeriodFirstTimeMinute: qo.NoticePeriodFirstTime,
			NoticePeriodTimes:           qo.NoticePeriodTimes,
			NoticePeriodInterval:        qo.NoticePeriodInterval,
			NoticePeriodFirstTime:       model.MemoTaskParseNoticePeriodFirstTime(onceTime, qo.NoticePeriodFirstTime),
			CreateTime:                  time.Now(),
		}
		conn.Table(config.TableMemoTask).Create(&task)
	}

	return task
}

func AddMemoFromTask(task model.MemoTask) {
	if task == (model.MemoTask{}) {
		logger.Logger.Info("task is empty")
		return
	}
	conn := config.GetConn()
	if task.NoticeType == model.NOTICE_TYPE_ONCE {
		memo := model.Memo{
			Title:      task.Title,
			DescShow:   task.Desc,
			TaskId:     task.Id,
			NoticeTime: task.NoticePeriodFirstTime,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		conn.Table(config.TableMemo).
			Create(&memo)
	}
}
