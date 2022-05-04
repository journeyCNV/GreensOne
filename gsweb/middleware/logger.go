package middleware

import (
	"github.com/journeycnv/greensone/gsweb"
	"time"
)

func LoggerDefault() gsweb.HandlerFunc {
	logger := gsweb.Logger()
	return func(c *gsweb.Context) error {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		costTime := endTime.Sub(startTime)

		reqMethod := c.Method()
		reqUri := c.Uri()
		clientIP := c.ClientIp()

		logger.Infof("|%13v |%15s |%s |%s|",
			costTime, clientIP, reqMethod, reqUri)
		return nil
	}
}
