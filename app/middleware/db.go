package middleware

import (
	"github.com/gin-gonic/gin"
	"sse_demo/database"
)

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
		newRedis := database.GetRdbInstance()
		c.Set("redis", newRedis)
		c.Next()
	}
}
