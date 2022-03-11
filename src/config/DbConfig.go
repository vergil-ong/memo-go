package config

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

const TableMemo = "memo"
const TableMemoTask = "memo_task"

type DbProperty struct {
	Username string
	Password string
	Host     string
	Port     uint32
	DataBase string
}

var instance *DbProperty
var once sync.Once
var DBInstance *gorm.DB

func buildDataSourceProperty() *DbProperty {
	instance = &DbProperty{
		Username: viper.GetString("datasource.username"),
		Password: viper.GetString("datasource.password"),
		Host:     viper.GetString("datasource.host"),
		Port:     viper.GetUint32("datasource.port"),
		DataBase: viper.GetString("datasource.database"),
	}

	return instance
}

func GetDataSourceProperty() *DbProperty {
	once.Do(func() {
		instance = buildDataSourceProperty()
	})

	return instance
}

func GetConn() *gorm.DB {
	once.Do(func() {
		dbProperty := buildDataSourceProperty()
		dsn := dbProperty.Username + ":" + dbProperty.Password +
			"@tcp(" + dbProperty.Host + ":" + strconv.Itoa(int(dbProperty.Port)) + ")/" +
			dbProperty.DataBase + "?charset=utf8mb4&parseTime=True&loc=Local"

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Errorf("fail to open mysql connection %s \n", err))
		}

		sqlDB, err := db.DB()
		if err != nil {
			panic(fmt.Errorf("fail to open mysql connection %s \n", err))
		}

		sqlDB.SetMaxIdleConns(8)
		sqlDB.SetMaxOpenConns(64)
		sqlDB.SetConnMaxLifetime(time.Minute)

		DBInstance = db
	})

	return DBInstance
}
