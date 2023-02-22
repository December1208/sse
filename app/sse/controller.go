package sse_server

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io"
	"sse/app/api_error"
	"sse/app/webutil"
	"sse/util"
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

	util.MyLogger.Info("connection open", zap.String("client_id", client.ClientId))
	_ = c.Stream(func(w io.Writer) bool {
		select {
		case msg := <-client.MessageChan:
			c.SSEvent("message", msg)
		default:
			// ignore
		}
		return true
	})
	channel.removeClient(client)
	util.MyLogger.Info("connection closed", zap.String("client_id", client.ClientId))
	//if err != nil {
	//	channel.removeClient(client)
	//	return
	//}
}

func (s Controller) Publish(c *gin.Context) {
	ctx := webutil.NewContext(c)
	type reqJson struct {
		Channel string `json:"channel"`
		Message string `json:"message"`
	}
	var vData reqJson
	if result := ctx.V2BindJsonAndValidate(&vData); !result {
		return
	}
	if vData.Channel == "" || vData.Message == "" {
		ctx.V2AbortWithError(api_error.InvalidParam)
		return
	}
	channel, ok := PortalInstance.Channels[vData.Channel]
	if !ok {
		ctx.V2AbortWithError(api_error.InvalidParam)
		return
	}
	channel.PublishMsg(vData.Message)
	ctx.V2JSON("ok")
}

func (s Controller) Connect(c *gin.Context) {
	//升级get请求为webSocket协议
	upgrade := &websocket.Upgrader{
		ReadBufferSize:    0,
		EnableCompression: false,
	}
	ws, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	for {
		//读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		if string(message) == "ping" {
			message = []byte("pong")
		}
		//写入ws数据
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func AddInnerApiRouter(router *gin.RouterGroup) {

}

func AddApiRouter(router *gin.RouterGroup) {
	controller := new(Controller)
	router.GET("/stream", controller.Subscribe)
	router.POST("/publish", controller.Publish)
	router.GET("/connection", controller.Connect)
}
