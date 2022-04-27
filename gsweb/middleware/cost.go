package middleware

import (
	"GreensOne/gsweb"
	"log"
	"time"
)

/**
请求时长统计中间件
*/

func Cost() gsweb.HandlerFunc {
	return func(c *gsweb.Context) error {
		start := time.Now()

		log.Printf("api uri start: %v", c.GetRequest().RequestURI)
		c.Next()

		end := time.Now()
		cost := end.Sub(start)

		log.Printf("api uri end: %v, cost: %v", c.GetRequest().RequestURI, cost.Seconds())
		return nil
	}
}
