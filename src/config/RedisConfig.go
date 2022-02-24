package config

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"runtime"
)

func GetRedisPoolClient() *redis.Client {
	redisDb := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.host"),
		Password: viper.GetString("redis.password"),
		DB:       0,
		PoolSize: runtime.NumCPU() * 4,
	})

	return redisDb
}
