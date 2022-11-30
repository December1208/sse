package middleware

import "github.com/gin-gonic/gin"

type MyContext struct {
	*gin.Context
}

func RegMiddleware(eng *gin.Engine) {
	eng.Use(CORSMiddleware())
	//eng.Use(DBMiddleware())
	eng.Use(LoadUserMiddleware())
	//eng.Use(RedisMiddleware())
}
