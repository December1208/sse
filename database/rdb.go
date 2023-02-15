package database

import (
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"sse/util"
)

var rdb *redis.Client

func init() {
	redisAddr := viper.GetString("redis.engine")
	pwd := viper.GetString("redis.pwd")
	redisDB := viper.GetInt("redis.db")
	rdb = redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: pwd,     // no password set
		DB:       redisDB, // use default DB
	})
	util.MyLogger.Info("redis init.")
}

func GetRdbInstance() *redis.Client {
	return rdb
}
