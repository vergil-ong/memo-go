package service

import (
	"MemoProjects/src/config"
	"MemoProjects/src/constant"
	"MemoProjects/src/logger"
	"context"
	"encoding/json"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/url"
)

type AuthCode2SessionVo struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

func AuthCode2Session(code string) string {
	request, err := http.NewRequest(http.MethodGet, constant.AuthCode2session, nil)
	if err != nil {
		panic(err)
	}

	params := make(url.Values)
	params.Add("appid", constant.AppId)
	params.Add("secret", constant.AppSecret)
	params.Add("js_code", code)
	params.Add("grant_type", "authorization_code")

	request.URL.RawQuery = params.Encode()

	r, err := http.DefaultClient.Do(request)

	if err != nil {
		panic(err)
	}
	defer func() { _ = r.Body.Close() }()
	body, err := ioutil.ReadAll(r.Body)
	logger.Logger.Info("response is " + string(body))

	var result AuthCode2SessionVo
	json.Unmarshal(body, &result)

	client := config.GetRedisPoolClient()
	s, err := client.Ping(context.Background()).Result()
	logger.Logger.Info("ping result is " + s)
	timeOutMinutes := viper.GetDuration("session.timeout")
	authCode := GetAuthCode(code)
	client.Set(context.Background(), authCode, logger.GetJson(result), timeOutMinutes)
	return authCode
}
