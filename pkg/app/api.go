package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

func APIConfiguration() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	formatter := func(param gin.LogFormatterParams) string {
		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}

	var skipped = []string{""}
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: skipped, Formatter: formatter}))

	router.Use(gin.Recovery())

	return router
}
