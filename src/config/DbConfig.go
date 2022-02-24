package config

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strconv"
	"sync"
)

const TableMemo = "memo"

type DbProperty struct {
	Username string
	Password string
	Host     string
	Port     uint32
	DataBase string
}

var instance *DbProperty
var once sync.Once

func GetDataSourceProperty() *DbProperty {
	once.Do(func() {
		instance = &DbProperty{
			Username: viper.GetString("datasource.username"),
			Password: viper.GetString("datasource.password"),
			Host:     viper.GetString("datasource.host"),
			Port:     viper.GetUint32("datasource.port"),
			DataBase: viper.GetString("datasource.database"),
		}
	})

	return instance
}

func GetConn() *gorm.DB {
	dbProperty := GetDataSourceProperty()

	dsn := dbProperty.Username + ":" + dbProperty.Password +
		"@tcp(" + dbProperty.Host + ":" + strconv.Itoa(int(dbProperty.Port)) + ")/" +
		dbProperty.DataBase + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("fail to open mysql connection %s \n", err))
	}

	return db
}
