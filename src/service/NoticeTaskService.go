package service

import (
	"MemoProjects/src/config"
	"MemoProjects/src/logger"
	"MemoProjects/src/model"
	"strconv"
	"time"
)

func AddNoticeTask(qo model.NoticeQo, user model.User) model.MemoTask {
	conn := config.GetConn()

	task := model.MemoTask{}

	noticeType := model.MemoTaskParseNoticeType(qo.NoticeType)
	if noticeType == model.NoticeTypeOnce {
		onceTime := model.MemoTaskParseNoticeOnceTime(qo.NoticeOnceCal, qo.NoticeOnceTime)
		task = model.MemoTask{
			Title:                       qo.Title,
			Desc:                        qo.Desc,
			UserId:                      user.Id,
			NoticeType:                  noticeType,
			NoticeOnceTime:              onceTime,
			NoticePeriodFirstTimeMinute: qo.NoticePeriodFirstTime,
			NoticePeriodTimes:           qo.NoticePeriodTimes,
			NoticePeriodInterval:        qo.NoticePeriodInterval,
			NoticePeriodFirstTime:       model.MemoTaskParseNoticePeriodFirstTime(onceTime, qo.NoticePeriodFirstTime),
			CreateTime:                  time.Now(),
			NoticeTaskStatus:            model.NoticeTaskStatusCreate,
		}
		conn.Table(config.TableMemoTask).Create(&task)
	}

	return task
}

func AddMemoFromTask(task model.MemoTask, user model.User) {
	if task == (model.MemoTask{}) {
		logger.Logger.Info("AddMemoFromTask task is empty")
		return
	}
	conn := config.GetConn()
	if task.NoticeType == model.NoticeTypeOnce {
		//取最后一次提醒的时间，这样 提醒了之后 查看列表能看到
		addMinutes := task.NoticePeriodInterval * task.NoticePeriodTimes
		duration, _ := time.ParseDuration(strconv.Itoa(addMinutes) + "m")
		lastNoticeTime := task.NoticeOnceTime.Add(duration)
		if lastNoticeTime.Before(task.NoticeOnceTime) {
			lastNoticeTime = task.NoticeOnceTime
		}
		memo := model.Memo{
			Title:      task.Title,
			DescShow:   task.Desc,
			TaskId:     task.Id,
			UserId:     user.Id,
			NoticeTime: lastNoticeTime,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		conn.Table(config.TableMemo).
			Create(&memo)
	}
}

func AddMemoNotice(task model.MemoTask) {
	if task == (model.MemoTask{}) {
		logger.Logger.Info("AddMemoNotice task is empty")
		return
	}

	conn := config.GetConn()

	if task.NoticeType == model.NoticeTypeOnce {
		//一次性任务
		//根据 预计第一次提醒时间，间隔，次数 生成一组 提醒，放入数据库中
		//数据库 再记录 是否 已经发过消息
		memoNotices := make([]*model.MemoNotice, task.NoticePeriodTimes)
		var noticeTime time.Time
		for i := 0; i < task.NoticePeriodTimes; i++ {
			addMinutes := task.NoticePeriodInterval * i
			duration, _ := time.ParseDuration(strconv.Itoa(addMinutes) + "m")
			noticeTime = task.NoticePeriodFirstTime.Add(duration)
			notice := model.MemoNotice{
				NoticeTime:   noticeTime,
				Title:        task.Title,
				DescShow:     task.Desc,
				TaskId:       task.Id,
				CreateTime:   time.Now(),
				NoticeStatus: model.MemoNoticeStatusRecord,
			}
			memoNotices[i] = &notice
		}

		conn.Table(config.TableMemoNotice).CreateInBatches(memoNotices, 100)
	}
}

func UpdateMemoTaskDone(memoId int) {
	conn := config.GetConn()

	var memo model.Memo
	conn.Table(config.TableMemo).Where("id = ?", memoId).First(&memo)

	if memo == (model.Memo{}) {
		logger.Logger.Info("memo id null " + strconv.Itoa(memoId))
		return
	}

	conn.Table(config.TableMemoTask).
		Where("id = ?", memo.TaskId).
		Updates(model.MemoTask{NoticeTaskStatus: model.NoticeTaskStatusDone})
	conn.Table(config.TableMemo).Delete(&memo)
}

func MemoTaskInfo(memoId int) model.MemoTask {
	conn := config.GetConn()
	var memo model.Memo
	conn.Table(config.TableMemo).Where("id = ?", memoId).First(&memo)

	if memo == (model.Memo{}) {
		logger.Logger.Info("memo id null " + strconv.Itoa(memoId))
		return model.MemoTask{}
	}

	var memoTask model.MemoTask
	conn.Table(config.TableMemoTask).
		Where("id = ?", memo.TaskId).
		First(&memoTask)

	return memoTask
}
