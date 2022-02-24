package controller

import (
	"MemoProjects/src/config"
	"MemoProjects/src/logger"
	"MemoProjects/src/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func TestRoot(context *gin.Context) {
	context.JSON(200, gin.H{"who": "HelloWorld"})
}

func TestConfig(context *gin.Context) {
	context.JSON(200,
		gin.H{"appName": viper.GetString("app.name"),
			"host": viper.GetString("datasource.host")})
}

func TestUserGetId(context *gin.Context) {

	conn := config.GetConn()

	id := context.Param("id")
	table := model.TestTable{}
	conn.Table("test_table1").First(&table, id)

	bytes, _ := json.Marshal(table)

	logger.Logger.Info("table is", zap.String("table", string(bytes)))

	context.JSON(200,
		table)
}
