package middleware

import (
	"GreensOne/gsweb"
	"context"
	"fmt"
	"log"
	"time"
)

/**
超时控制中间件
fun : 业务逻辑 handler
d: 超时时间
*/
func TimeoutHandler(d time.Duration) gsweb.HandlerFunc {
	return func(c *gsweb.Context) error {
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)

		// 初始化超时context
		timeCtx, cancel := context.WithTimeout(c.BaseContext(), d)
		defer cancel()

		c.GetRequest().WithContext(timeCtx)

		go func() {
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			c.Next()
			finish <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			log.Println(p)
			c.GetResponse().WriteHeader(500)
		case <-finish:
			fmt.Println("finish")
		case <-timeCtx.Done():
			c.SetHasTimeout()
			c.GetResponse().Write([]byte("time out"))
		}

		return nil
	}
}
