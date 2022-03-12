package service

import (
	"MemoProjects/src/config"
	"MemoProjects/src/logger"
	"MemoProjects/src/model"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const AuthCodePrefix = "AUTH_"

func GetAuthCode(code string) string {
	return AuthCodePrefix + code
}

func AuthLoginToken(token string) bool {

	client := config.GetRedisPoolClient()
	result, err := client.Exists(context.Background(), token).Result()
	if err != nil {
		logger.Logger.Error("cannot get check token " + token)
		return false
	}

	if result > 0 {
		timeOutMinutes := viper.GetDuration("session.timeout")
		client.Expire(context.Background(), token, timeOutMinutes)
		return true
	} else {
		return false
	}
}

func GetAuthInfo(ginContext *gin.Context) (AuthCode2SessionVo, model.User) {
	authToken := ginContext.GetHeader("authentication_token")
	client := config.GetRedisPoolClient()
	result, err := client.Get(context.Background(), authToken).Result()
	var sessionVo AuthCode2SessionVo
	var user model.User
	if err != nil {
		logger.Logger.Error("cannot get check token " + authToken)
		return sessionVo, user
	}

	err = json.Unmarshal([]byte(result), &sessionVo)
	if err != nil {
		logger.Logger.Error("json.Unmarshal error " + authToken)
		return sessionVo, user
	}

	dbCon := config.GetConn()
	dbCon.
		Table(config.TableUser).
		Where("open_id = ?", sessionVo.Openid).
		First(&user)

	return sessionVo, user
}
