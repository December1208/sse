package demo

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HealthController struct{}

func (h HealthController) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": 0,
	})
}

//func (h HealthController) AddApiRouter(router *gin.RouterGroup) {
//	router.GET("/healthy", h.Health)
//}
