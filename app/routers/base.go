package routers

import (
	"github.com/gin-gonic/gin"
	"sse_demo/app/demo"
	"sse_demo/app/sse"
)

func RegRouters(eng *gin.Engine) {

	healthController := new(demo.HealthController)
	eng.GET("/health", healthController.Health)

	v2api := eng.Group("/v2/api")
	{
		sse.AddApiRouter(v2api)
	}
	v2InnerApi := eng.Group("/v2/inner_api")
	{
		sse.AddInnerApiRouter(v2InnerApi)
	}
}
