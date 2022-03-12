package service

import (
	"MemoProjects/src/config"
	"MemoProjects/src/constant"
	"MemoProjects/src/logger"
	"MemoProjects/src/model"
	"MemoProjects/src/utils"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func DoNoticeTimeService(ticker time.Ticker) {
	for t := range ticker.C {
		fmt.Println("Tick at", t)
		SendMemoNotice()
	}
}

func SendMemoNotice() {
	conn := config.GetConn()
	//查询 前后1分钟的数据 进行发送
	oneMinuteBefore, _ := time.ParseDuration("-1m")
	oneMinuteAfter, _ := time.ParseDuration("30s")
	startTime := time.Now().Add(oneMinuteBefore)
	endTime := time.Now().Add(oneMinuteAfter)

	var memoNotices []model.MemoNotice
	conn.
		Table(config.TableMemoNotice).
		Where("notice_status = 1 AND notice_time BETWEEN ? AND ?", startTime, endTime).
		Find(&memoNotices)

	taskOpenMap := getTaskIdOpenIdMap(memoNotices, conn)

	if len(taskOpenMap) == 0 {
		logger.Logger.Info("taskOpenMap is Empty")
		return
	}

	memoNoticeIds := make([]uint, len(memoNotices))
	for index, notice := range memoNotices {
		openId := taskOpenMap[notice.TaskId]
		if openId == "" {
			logger.Logger.Info("openId is empty id is " + strconv.Itoa(int(notice.Id)) + "task id is " + strconv.Itoa(int(notice.TaskId)))
			continue
		}
		noticeTimeStr := notice.NoticeTime.Format("2006-01-02 15:04:05")
		SendSubscribeMsg(openId, constant.TemplateId1, noticeTimeStr, notice.Title, notice.DescShow)
		memoNoticeIds[index] = notice.Id
	}

	conn.Table(config.TableMemoNotice).
		Where("id in ?", memoNoticeIds).
		Updates(model.MemoNotice{NoticeStatus: model.MemoNoticeStatusSend})
}

func getTaskIdOpenIdMap(memoNotices []model.MemoNotice, conn *gorm.DB) map[uint]string {
	var taskIdOpenIdMap map[uint]string
	if len(memoNotices) == 0 {
		logger.Logger.Info("memoNotices is empty")
		return taskIdOpenIdMap
	}

	taskIds := getTaskIds(memoNotices)
	if len(taskIds) == 0 {
		logger.Logger.Info("taskIds is empty")
		return taskIdOpenIdMap
	}

	taskIdUserIdMap := getTaskIdUserIdMap(taskIds, conn)
	if len(taskIdUserIdMap) == 0 {
		logger.Logger.Info("taskIdUserIdMap is empty")
		return taskIdOpenIdMap
	}

	taskIdOpenIdMap = getTaskIdOpenIdMapByTaskUser(taskIdUserIdMap, conn)
	return taskIdOpenIdMap
}

func getTaskIds(memoNotices []model.MemoNotice) []uint {
	var taskIds []uint

	taskMap := make(map[uint]uint)
	for _, memoNotice := range memoNotices {
		taskMap[memoNotice.TaskId] = memoNotice.TaskId
	}

	if len(taskMap) == 0 {
		logger.Logger.Info("taskMap is empty")
		return taskIds
	}
	taskIds = utils.MapKey2Slice(taskMap)

	return taskIds
}

func getTaskIdUserIdMap(taskIds []uint, conn *gorm.DB) map[uint]uint {
	var taskIdUserIdMap = make(map[uint]uint)

	var memoTasks []model.MemoTask
	conn.Table(config.TableMemoTask).
		Where("id in ?", taskIds).
		Find(&memoTasks)

	if len(memoTasks) == 0 {
		logger.Logger.Info("memoTasks is empty")
		return taskIdUserIdMap
	}

	for _, memoTask := range memoTasks {
		taskIdUserIdMap[memoTask.Id] = memoTask.UserId
	}
	return taskIdUserIdMap
}

func getTaskIdOpenIdMapByTaskUser(taskIdUserIdMap map[uint]uint, conn *gorm.DB) map[uint]string {
	var taskIdOpenIdMap = make(map[uint]string)

	userIds := utils.MapVal2Slice(taskIdUserIdMap)
	var memoUsers []model.User
	conn.Table(config.TableUser).
		Where("id in ?", userIds).
		Find(&memoUsers)
	if len(memoUsers) == 0 {
		logger.Logger.Info("memoUsers is empty")
		return taskIdOpenIdMap
	}

	var userIdMap = make(map[uint]model.User)
	for _, user := range memoUsers {
		userIdMap[user.Id] = user
	}

	for taskId, userId := range taskIdUserIdMap {
		user := userIdMap[userId]
		if user == (model.User{}) {
			continue
		}
		taskIdOpenIdMap[taskId] = user.OpenId
	}

	return taskIdOpenIdMap
}
