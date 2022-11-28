package middleware

import (
	"fmt"
	"sse_demo/app/api_error"
	"sse_demo/app/webutil"
	"sse_demo/service"

	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const BearerSchema = "Bearer "
const DefaultUUID = ""

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		curUser := &service.AuthCustomClaims{
			UserID: -1,
			UUID:   DefaultUUID,
			//K:    "",
			Key:            "",
			Iv:             "",
			Quality:        "",
			IsLcVideoAdmin: false,
		}
		clientIP := c.ClientIP()

		myCtx := webutil.NewContext(c)
		myCtx.SetStrVal("trace_id", uuid.New().String())
		authHeader := c.GetHeader("Authorization")
		sentrygin.GetHubFromContext(c).Scope().SetContext("userInfo", map[string]interface{}{
			"client_ip": c.ClientIP(),
			"user_id":   curUser.UserID,
			"uuid":      curUser.UUID,
			//"k":         curUser.K,
			"key":               curUser.Key,
			"iv":                curUser.Iv,
			"quality":           curUser.Quality,
			"is_lc_video_admin": curUser.IsLcVideoAdmin,
		})
		if len(authHeader) <= len(BearerSchema) {
			// myCtx.L.Warn("give an empty auth header", zap.String("clientip", clientIP))
			// sentrygin.GetHubFromContext(c).CaptureMessage(fmt.Sprintf("no auth header try to invoke: header: %s", authHeader))
			// c.AbortWithStatus(http.StatusUnauthorized)

		} else {
			tokenString := authHeader[len(BearerSchema):]
			token, err := service.JWTAuthService("jwt.secret").ValidateToken(tokenString)
			if err == nil && token.Valid {
				curUser = token.Claims.(*service.AuthCustomClaims)
			} else {
				myCtx.L.Warn(fmt.Sprintf("token %s, is invalid", tokenString), zap.String("clientip", clientIP), zap.Error(err))
			}
		}
		myCtx.SetUserInfo(curUser)
	}
}


func JwtRequiredMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		myCtx := webutil.NewContext(c)
		userInfo := myCtx.GetUserInfo()
		if userInfo.UserID == -1 {
			myCtx.V2AbortWithError(api_error.NotAuthenticated)
			c.Abort()
			return
		}
		if !userInfo.IsLcVideoAdmin {
			myCtx.V2AbortWithError(api_error.PermissionDenied)
			c.Abort()
			return
		}
	}
}
