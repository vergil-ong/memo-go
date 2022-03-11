package config

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"runtime"
)

var redisDBInstance *redis.Client

func GetRedisPoolClient() *redis.Client {
	once.Do(func() {
		redisDb := redis.NewClient(&redis.Options{
			Addr:     viper.GetString("redis.host"),
			Password: viper.GetString("redis.password"),
			DB:       0,
			PoolSize: runtime.NumCPU() * 4,
		})
		redisDBInstance = redisDb
	})

	return redisDBInstance
}
