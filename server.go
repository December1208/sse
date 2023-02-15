package main

import (
	"fmt"
	"runtime/debug"
	"sse_demo/app/middleware"
	"sse_demo/app/routers"
	"sse_demo/app/webutil"
	"sse_demo/util"
	"time"

	"github.com/gin-gonic/gin"
)

func ServerInit() {
	gin.SetMode(gin.ReleaseMode)
	r := createEng()
	util.MyLogger.Info("try to start listen on 8855")
	err := r.Run(":8855")
	if err != nil {
		panic(err.Error())
	}

}

func createEng() *gin.Engine {
	eng := gin.New()

	eng.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		//定制日志格式
		return fmt.Sprintf("[%s] - [%s] [%s] [%s] %d %s %s\"\n",
			param.TimeStamp.Format(time.RFC3339),
			param.ClientIP,
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ErrorMessage,
		)
	}))
	middleware.RegMiddleware(eng)
	routers.RegRouters(eng)
	webutil.InitValidator()
	return eng
}
func main() {

	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	defer func() {
		err := recover()
		if err != nil {
			util.MyLogger.Info(string(debug.Stack()))
		}
	}()
	ServerInit()
}
