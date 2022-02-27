package middleware

import (
	"MemoProjects/src/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"strings"
)

func AuthLoginToken() gin.HandlerFunc {
	return func(context *gin.Context) {

		if checkUrlNeedLogin(context.Request.URL.Path) {
			context.Next()
			return
		}

		authToken := context.GetHeader("authentication_token")
		localTest := "local_test" == authToken
		result := service.AuthLoginToken(authToken)

		if result || localTest {
			context.Next()
		} else {
			context.Abort()
			context.JSON(http.StatusNonAuthoritativeInfo, gin.H{"message": "身份验证失败"})
		}
	}
}

func checkUrlNeedLogin(url string) bool {
	noLoginUrls := viper.GetStringSlice("auth.no-login")
	for _, noLoginUrl := range noLoginUrls {
		if strings.Contains(url, noLoginUrl) {
			return true
		}
	}
	return false
}
