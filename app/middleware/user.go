package middleware

import (
	"github.com/gin-gonic/gin"
	"sse/app/webutil"
	"sse/service"
)

func LoadUserMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		ctx := webutil.NewContext(context)
		userData := context.GetHeader("X-Forwarded-User")
		if userData == "" {
			ctx.SetUserInfo(&service.User{ID: -1})
		}
		// TODO SetUserInfo
	}
}
