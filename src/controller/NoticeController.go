package controller

import (
	"MemoProjects/src/config"
	"MemoProjects/src/logger"
	"MemoProjects/src/model"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func NoticeList(context *gin.Context) {
	conn := config.GetConn()
	var memos []model.Memo
	conn.Table(config.TableMemo).Limit(10).Find(&memos)

	memoVos := make([]*model.MemoVo, len(memos))

	for index, memo := range memos {
		mVo := new(model.MemoVo)
		mVo.Id = memo.Id
		mVo.Title = memo.Title
		mVo.DescShow = memo.DescShow
		mVo.NoticeTime = memo.NoticeTime.UnixMilli()

		memoVos[index] = mVo
	}

	success := model.Success(memoVos)
	context.JSON(model.HttpSuccess, success)
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
