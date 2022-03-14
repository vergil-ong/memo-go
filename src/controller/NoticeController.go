package controller

import (
	"MemoProjects/src/config"
	"MemoProjects/src/logger"
	"MemoProjects/src/model"
	"MemoProjects/src/service"
	"MemoProjects/src/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func NoticeList(context *gin.Context) {
	conn := config.GetConn()

	_, user := service.GetAuthInfo(context)
	if user == (model.User{}) {
		model.ReturnSuccess(make([]*model.MemoVo, 0), context)
		return
	}

	var memos []model.Memo
	conn.
		Table(config.TableMemo).
		Where("user_id = ? ", user.Id).
		Order("notice_time asc").
		Limit(10).
		Find(&memos)

	memoVos := make([]*model.MemoVo, len(memos))

	for index, memo := range memos {
		mVo := new(model.MemoVo)
		mVo.Id = memo.Id
		mVo.Title = memo.Title
		mVo.DescShow = memo.DescShow
		mVo.NoticeTime = memo.NoticeTime.UnixMilli()

		diffTodayEnd := utils.GetDiffDays(utils.GetTodayEndUnixMilli(), memo.NoticeTime.UnixMilli())
		diffTodayStart := utils.GetDiffDays(memo.NoticeTime.UnixMilli(), utils.GetTodayStartUnixMilli())
		//days := utils.GetDiffDays(utils.GetTodayEndUnixMilli(), memo.NoticeTime.UnixMilli())
		noticeTimeShow := ""
		if diffTodayEnd > 0 {
			noticeTimeShow = strconv.Itoa(diffTodayEnd) + "天后" + memo.NoticeTime.Format("15:04:05")
		} else if diffTodayStart < 0 {
			noticeTimeShow = strconv.Itoa(diffTodayEnd) + "天前" + memo.NoticeTime.Format("15:04:05")
		} else {
			noticeTimeShow = "当天" + memo.NoticeTime.Format("15:04:05")
		}
		/*if days > 0 {
			noticeTimeShow = strconv.FormatInt(days, 10) + "天后" + memo.NoticeTime.Format("15:04:05")
		} else if days == 0 {
			noticeTimeShow = "当天" + memo.NoticeTime.Format("15:04:05")
		} else {
			noticeTimeShow = strconv.FormatInt(days, 10) + "天前" + memo.NoticeTime.Format("15:04:05")
		}*/
		mVo.NoticeTimeShow = noticeTimeShow

		memoVos[index] = mVo
	}

	model.ReturnSuccess(memoVos, context)
}

func NoticeAdd(context *gin.Context) {
	var memoQo model.NoticeQo

	context.BindJSON(&memoQo)
	logger.Logger.Info("memoQo is ", zap.String("memoQo", logger.GetJson(memoQo)))

	conn := config.GetConn()
	memo := model.Memo{
		Title:      memoQo.Title,
		DescShow:   memoQo.Desc,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	conn.Table(config.TableMemo).
		Select("Title", "DescShow", "CreateTime", "UpdateTime").
		Create(&memo)

	success := model.Success(memo)
	context.JSON(model.HttpSuccess, success)
}

func NoticeTaskAdd(context *gin.Context) {
	var memoQo model.NoticeQo

	var task model.MemoTask

	_, user := service.GetAuthInfo(context)
	if user == (model.User{}) {
		model.ReturnSuccess(task, context)
		return
	}

	context.BindJSON(&memoQo)
	logger.Logger.Info("memoQo is ", zap.String("memoQo", logger.GetJson(memoQo)))

	//生成task
	task = service.AddNoticeTask(memoQo, user)
	if task == (model.MemoTask{}) {
		model.ReturnSuccess(task, context)
		return
	}

	//根据task 生成 memo
	service.AddMemoFromTask(task, user)

	service.AddMemoNotice(task)
}

func NoticeMemoDone(context *gin.Context) {
	var memoQo model.NoticeQo
	err := context.BindJSON(&memoQo)

	if err != nil {
		logger.Logger.Info("NoticeMemoDone bind memoQo error")
		return
	}

	logger.Logger.Info("memoQo is ", zap.String("memoQo", logger.GetJson(memoQo)))

	memoId := memoQo.NoticeMemoId
	service.UpdateMemoTaskDone(memoId)
}

func NoticeMemoTask(context *gin.Context) {
	memoId := context.Param("memoId")
	memoIdInt, err := strconv.Atoi(memoId)
	if err != nil {
		logger.Logger.Error("NoticeMemoTask parse int error " + memoId)
		return
	}
	info := service.MemoTaskInfo(memoIdInt)

	model.ReturnSuccess(info, context)
}
