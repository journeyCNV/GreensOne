package middleware

import "GreensOne/gsweb"

/**
请求时长统计中间件
*/

func TimeStatistic() gsweb.ControllerHandler {
	return func(c *gsweb.Context) error {
		return nil
	}
}
