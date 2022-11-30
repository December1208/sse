package sse_server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"sse_demo/app/webutil"
	"sse_demo/util"
)

type Controller struct{}

func (s Controller) Subscribe(c *gin.Context) {
	ctx := webutil.NewContext(c)
	type reqJson struct {
		Channel string `form:"channel"`
	}
	var vData reqJson
	if result := ctx.V2BindQueryAndValidate(&vData); !result {
		return
	}

	channel := GetOrCreateChannel(vData.Channel)
	util.MyLogger.Info("channel: " + vData.Channel)
	client := &Client{
		ClientId:    uuid.New().String(),
		MessageChan: make(chan string),
		channel:     channel,
	}
	channel.Subscribe(client)

	//msg := <-client.MessageChan
	suc := c.Stream(func(w io.Writer) bool {
		if msg, ok := <-client.MessageChan; ok {
			c.SSEvent("message", msg)
			return true
		}
		return false
	})
	if suc {
		channel.removeClient(client)
	}
	//if err != nil {
	//	channel.removeClient(client)
	//	return
	//}
}

func AddInnerApiRouter(router *gin.RouterGroup) {

}

func AddApiRouter(router *gin.RouterGroup) {
	controller := new(Controller)
	router.GET("/stream", controller.Subscribe)
}
