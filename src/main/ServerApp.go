package main

import (
	"MemoProjects/src/config"
	"MemoProjects/src/controller"
	"MemoProjects/src/logger"
	"MemoProjects/src/middleware"
	"MemoProjects/src/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	// 配置文件读写
	config.InitViperConfig()
	// 日志框架初始化 控制台输出,日志文件,文件切割
	logger.InitLogger()

	engine := gin.Default()

	apiV1 := engine.Group("/api/v1")
	{
		apiV1.Use(middleware.AuthLoginToken())
		noticeApi := apiV1.Group("/notice")
		{
			noticeApi.GET("/list", controller.NoticeList)
			noticeApi.POST("/add", controller.NoticeAdd)
			noticeApi.POST("/task/add", controller.NoticeTaskAdd)
		}
		loginApi := apiV1.Group("/login")
		{
			loginApi.GET("/code/:code", controller.LoginByCode)
			loginApi.GET("/auth/:token", controller.AuthToken)
		}
	}

	engine.GET("/", controller.TestRoot)

	engine.GET("/config", controller.TestConfig)

	engine.GET("/user/:id", controller.TestUserGetId)

	//ticker := time.NewTicker(time.Minute)
	ticker := time.NewTicker(time.Second * 10)
	go service.DoNoticeTimeService(*ticker)

	err := engine.Run(":8001")
	if err != nil {
		logger.Logger.Error("start engine error ", zap.String("error", err.Error()))
	}

	shutDownServer(*ticker)
}

func shutDownServer(ticker time.Ticker) {
	quit := make(chan os.Signal)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Logger.Info("Shutdown Server")

	ticker.Stop()
}
