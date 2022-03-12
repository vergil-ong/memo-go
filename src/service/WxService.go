package service

import (
	"MemoProjects/src/config"
	"MemoProjects/src/constant"
	"MemoProjects/src/logger"
	"MemoProjects/src/model"
	"bytes"
	"context"
	"encoding/json"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type AuthCode2SessionVo struct {
	Openid     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type AccessTokenVo struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

type WxTemplateValue struct {
	Value string `json:"value"`
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

	var sessionVo AuthCode2SessionVo
	err = json.Unmarshal(body, &sessionVo)
	if err != nil {
		return ""
	}

	client := config.GetRedisPoolClient()
	//s, err := client.Ping(context.Background()).Result()
	//logger.Logger.Info("ping result is " + s)
	timeOutMinutes := viper.GetDuration("session.timeout")
	authCode := GetAuthCode(code)
	client.Set(context.Background(), authCode, logger.GetJson(sessionVo), timeOutMinutes)

	checkAndSaveUser(sessionVo.Openid)

	return authCode
}

func checkAndSaveUser(openId string) {
	dbCon := config.GetConn()
	var user model.User
	dbCon.
		Table(config.TableUser).
		Where("open_id = ?", openId).
		First(&user)

	if user == (model.User{}) {
		user = model.User{
			OpenId: openId,
		}
		dbCon.
			Table(config.TableUser).
			Create(&user)
	}
}

func GetAccessToken() string {

	redisClient := config.GetRedisPoolClient()
	accessTokenStr, _ := redisClient.Get(context.Background(), constant.AccessTokenReidsKey).Result()

	if accessTokenStr != "" {
		return accessTokenStr
	}

	request, err := http.NewRequest(http.MethodGet, constant.AccessToken, nil)
	if err != nil {
		panic(err)
	}

	params := make(url.Values)
	params.Add("appid", constant.AppId)
	params.Add("secret", constant.AppSecret)
	params.Add("grant_type", "client_credential")

	request.URL.RawQuery = params.Encode()

	r, err := http.DefaultClient.Do(request)

	if err != nil {
		panic(err)
	}
	defer func() { _ = r.Body.Close() }()
	body, err := ioutil.ReadAll(r.Body)
	logger.Logger.Info("response is " + string(body))

	var accessTokenVo AccessTokenVo
	err = json.Unmarshal(body, &accessTokenVo)
	if err != nil {
		return ""
	}

	redisClient.SetEX(context.Background(), constant.AccessTokenReidsKey, accessTokenVo.AccessToken, time.Hour)

	return accessTokenVo.AccessToken
}

func SendSubscribeMsg(
	openId string,
	templateId string,
	noticeTimeStr string,
	title string,
	desc string,
) {
	accessToken := GetAccessToken()
	url := constant.MsgSubscribeSend + "?access_token=" + accessToken

	dataMap := map[string]interface{}{
		"time1": WxTemplateValue{Value: noticeTimeStr},
		//"time1":WxTemplateValue{Value: "2019年10月1日 15:01"},
		"thing3": WxTemplateValue{Value: title},
		"thing4": WxTemplateValue{Value: desc},
	}

	param := map[string]interface{}{
		"access_token":      accessToken,
		"touser":            openId,
		"template_id":       templateId,
		"data":              dataMap,
		"miniprogram_state": "developer",
	}
	paramBytes, _ := json.Marshal(param)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(paramBytes))
	request.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}

	r, err := http.DefaultClient.Do(request)

	if err != nil {
		panic(err)
	}
	defer func() { _ = r.Body.Close() }()
	body, err := ioutil.ReadAll(r.Body)
	logger.Logger.Info("response is " + string(body))

}
