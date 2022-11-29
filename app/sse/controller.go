package sse

import (
	"github.com/gin-gonic/gin"
)

type HLSVideoController struct{}

func (h HLSVideoController) CreatePlayList(c *gin.Context) {

}

func AddInnerApiRouter(router *gin.RouterGroup) {
	//hlsVideoController := new(HLSVideoController)

}

func AddApiRouter(router *gin.RouterGroup) {
	//hlsVideoController := new(HLSVideoController)
}
