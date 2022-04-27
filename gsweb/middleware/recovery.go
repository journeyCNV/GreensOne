package middleware

import "GreensOne/gsweb"

/**
全局捕获panic中间件
*/
func Recovery() gsweb.HandlerFunc {
	return func(c *gsweb.Context) error {
		// 捕获c.Next的出现的panic
		defer func() {
			if err := recover(); err != nil {
				c.SetStatus(500).Json(err)
			}
		}()
		c.Next()
		return nil
	}
}
