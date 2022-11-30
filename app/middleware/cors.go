package middleware

import (
	"fmt"
	"sse_demo/app/webutil"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

var globalAllowHosts map[string]bool

func getConfigAllowHosts() map[string]bool {
	if nil != globalAllowHosts {
		return globalAllowHosts
	}
	globalAllowHosts = make(map[string]bool)
	allowHosts := viper.GetStringSlice("cors.allowhosts")
	for _, allowHost := range allowHosts {
		globalAllowHosts[allowHost] = true
	}
	return globalAllowHosts
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		myCtx := webutil.NewContext(c)

		headers := viper.GetString("cors.headers")

		allowHosts := getConfigAllowHosts()
		origin := c.Request.Header.Get("Origin")
		index := strings.Index(origin, ".")
		curHost := origin[index+1:]
		myCtx.Log.Info(fmt.Sprintf("origin: %s, curHost: %s", origin, curHost))
		if strings.Contains(origin, "localhost:") {
			if _, ok := allowHosts["localhost"]; ok {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			}
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		if _, ok := allowHosts[curHost]; ok {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", headers)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}
