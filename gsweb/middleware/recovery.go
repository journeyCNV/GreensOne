package middleware

import "GreensOne/gsweb"

func Recovery() gsweb.ControllerHandler {
	return func(c *gsweb.Context) error {
		// 捕获c.Next的出现的panic
		defer func() {
			if err := recover(); err != nil {
				c.Json(500, err)
			}
		}()
		c.Next()
		return nil
	}
}
