package main

import (
	"MemoProjects/src/config"
	"MemoProjects/src/controller"
	"MemoProjects/src/logger"
	"MemoProjects/src/middleware"
	"github.com/gin-gonic/gin"
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

	engine.Run(":8001")
}
