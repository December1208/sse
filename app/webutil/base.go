package webutil

import (
	"github.com/go-redis/redis/v8"
	"net/http"
	"sse_demo/app/api_error"
	"sse_demo/service"
	"sse_demo/util"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
)

type BaseContext struct {
	db       *gorm.DB
	redis    *redis.Client
	strDatas map[string]string
	Log      *zap.Logger
}

func NewBaseContext() *BaseContext {
	strData := make(map[string]string)

	logger, _ := zap.NewProduction()
	return &BaseContext{nil, nil, strData, logger}
}

func (ctx BaseContext) SetStrVal(key, val string) *BaseContext {
	ctx.strDatas[key] = val
	return &ctx
}

func (ctx BaseContext) GetStrBy(key string) (string, bool) {
	val, existed := ctx.strDatas["key"]
	if !existed {
		ctx.Log.Warn("try to get by failed: ", zap.String("key", key))
	}

	return val, existed
}

func (ctx *BaseContext) SetDB(val *gorm.DB) *BaseContext {
	ctx.db = val
	return ctx
}

func (ctx BaseContext) GetDB() *gorm.DB {
	if ctx.db == nil {
		panic("db is nil")
	}
	return ctx.db
}

func (ctx *BaseContext) SetRedis(val *redis.Client) *BaseContext {
	ctx.redis = val
	return ctx
}

func (ctx BaseContext) GetRedis() *redis.Client {

	return ctx.redis
}

type MyWebContext struct {
	*BaseContext
	*gin.Context
	*util.SentryLog
}

func NewContext(c *gin.Context) *MyWebContext {
	baseContext := NewBaseContext()
	db, existed := c.Get("db")
	if existed && db != nil {
		baseContext.SetDB(db.(*gorm.DB))
	}
	ctxRedis, existed := c.Get("redis")
	if existed && ctxRedis != nil {
		baseContext.SetRedis(ctxRedis.(*redis.Client))
	}
	return &MyWebContext{baseContext, c, &util.SentryLog{L: baseContext.Log, Ctx: c}}
}

func (ctx *MyWebContext) GetUserInfo() *service.AuthCustomClaims {
	userInfo, exists := ctx.Get("userInfo")
	if exists && nil != userInfo {
		return userInfo.(*service.AuthCustomClaims)
	}
	nonameUserInfo := service.AuthCustomClaims{
		UUID: "",
	}
	return &nonameUserInfo
}

func (ctx *MyWebContext) SetUserInfo(claims *service.AuthCustomClaims) {
	ctx.Set("userInfo", claims)
}

func (ctx *MyWebContext) V2JSON(obj interface{}) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    obj,
		"success": true,
		"detail":  "",
	})
}

func (ctx *MyWebContext) V2AbortWithError(apiError api_error.APIError) {
	ctx.JSON(http.StatusOK, gin.H{
		"code":    apiError.ErrorCode,
		"detail":  apiError.ErrorMsg,
		"data":    nil,
		"success": false,
	})
}

func (ctx *MyWebContext) V2AbortWithErrorAndMsg(apiError api_error.APIError, msg string) {
	if msg == "" {
		msg = apiError.ErrorMsg
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    apiError.ErrorCode,
		"detail":  msg,
		"data":    nil,
		"success": false,
	})
}

func (ctx *MyWebContext) V2AbortWithOldErrorAndMsg(apiError api_error.APIOldError, msg string) {
	if msg == "" {
		msg = apiError.ErrorMsg
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":  apiError.ErrorCode,
		"error": map[string]interface{}{"message": msg},
		"desc":  msg,
		"data":  nil,
	})
}

func (ctx *MyWebContext) V2BindJsonAndValidate(vData interface{}) bool {
	if err := ctx.ShouldBindJSON(vData); err != nil {
		util.GetLogger().Warn("get param failed", zap.Error(err))
		ctx.V2AbortWithErrorAndMsg(api_error.InvalidParam, "")
		return false
	}
	util.GetLogger().Info("get param success", zap.Any("data", vData))
	if err := Validate.Struct(vData); err != nil {
		util.GetLogger().Warn("validate failed", zap.Error(err))
		ctx.V2AbortWithErrorAndMsg(api_error.InvalidParam, "")
		return false
	}
	return true
}

func (ctx *MyWebContext) V2BindAndValidate(vData interface{}) bool {
	if err := ctx.ShouldBind(vData); err != nil {
		util.GetLogger().Warn("get param failed", zap.Error(err))
		ctx.V2AbortWithErrorAndMsg(api_error.InvalidParam, "")
		return false
	}
	util.GetLogger().Info("get param success", zap.Any("data", vData))
	if err := Validate.Struct(vData); err != nil {
		util.GetLogger().Warn("validate failed", zap.Error(err))
		ctx.V2AbortWithErrorAndMsg(api_error.InvalidParam, "")
		return false
	}
	return true
}

func (ctx *MyWebContext) V2BindQueryAndValidate(vData interface{}) bool {
	if err := ctx.ShouldBindQuery(vData); err != nil {
		util.GetLogger().Warn("get param failed", zap.Error(err))
		ctx.V2AbortWithErrorAndMsg(api_error.InvalidParam, "")
		return false
	}
	if err := Validate.Struct(vData); err != nil {
		util.GetLogger().Warn("validate failed", zap.Error(err))
		ctx.V2AbortWithErrorAndMsg(api_error.InvalidParam, "")
		return false
	}
	return true
}

func (ctx *MyWebContext) DataResponse(code int, data []byte) {
	ctx.Data(code, "text/html; charset=utf-8", data)
}
