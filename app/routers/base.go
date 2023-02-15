package routers

import (
	"github.com/gin-gonic/gin"
	"sse/app/demo"
	sse_server "sse/app/sse"
)

func RegRouters(eng *gin.Engine) {

	healthController := new(demo.HealthController)
	eng.GET("/health", healthController.Health)

	v2api := eng.Group("/v2/api")
	{
		sse_server.AddApiRouter(v2api)
	}
	v2InnerApi := eng.Group("/v2/inner_api")
	{
		sse_server.AddInnerApiRouter(v2InnerApi)
	}
}
