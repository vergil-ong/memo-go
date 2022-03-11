package service

import (
	"MemoProjects/src/config"
	"MemoProjects/src/logger"
	"context"
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
