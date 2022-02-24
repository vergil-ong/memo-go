package controller

import (
	"MemoProjects/src/logger"
	"MemoProjects/src/model"
	"MemoProjects/src/service"
	"github.com/gin-gonic/gin"
)

func LoginByCode(context *gin.Context) {
	code := context.Param("code")

	logger.Logger.Info("code is " + code)

	authCode := service.AuthCode2Session(code)

	success := model.Success(authCode)
	context.JSON(model.HttpSuccess, success)
}

func AuthToken(context *gin.Context) {
	token := context.Param("token")

	logger.Logger.Info("token is " + token)

	success := model.Success(true)
	context.JSON(model.HttpSuccess, success)
}
