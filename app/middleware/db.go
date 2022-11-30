package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"sse_demo/database"
)

func NewRedis() *redis.Client {
	useRedis := viper.GetBool("setting.use_redis")
	if !useRedis {
		return nil
	}
	redisAddr := viper.GetString("redis.engine")
	pwd := viper.GetString("redis.pwd")
	redisDB := viper.GetInt("redis.db")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: pwd,     // no password set
		DB:       redisDB, // use default DB
	})

	// pong, err := rdb.Ping(ctx).Result()
	// fmt.Println(pong, err)
	// Output: PONG <nil>
	return rdb
}

func DBMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		db := database.GetDatabase().Begin()
		c.Set("db", db)
		c.Next()
		db.RollbackUnlessCommitted()
	}
}

func RedisMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		newRedis := NewRedis()
		c.Set("redis", newRedis)
		c.Next()
	}
}
