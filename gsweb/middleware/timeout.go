package middleware

import (
	"GreensOne/gsweb"
	"context"
	"fmt"
	"log"
	"time"
)

/**
fun : 业务逻辑 handler
d: 超时时间
*/
func TimeoutHandler(fun gsweb.ControllerHandler, d time.Duration) gsweb.ControllerHandler {
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
			fun(c)
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

		c.Next()
		return nil
	}
}
