package test

import (
	"GreensOne/gsweb"
)

func Register(g *gsweb.GreensCore) {
	g.Get("/welcome", TestH1(), WelcomeHandler)
	group := g.Group("/show")
	{
		group1 := group.Group("/door")
		group1.Use(TestH1()) // 批量使用中间件
		{
			group1.Get("/open", TestH2(), TestHandler)
		}
	}
}

func WelcomeHandler(ctx *gsweb.Context) error {
	ctx.Json(200, "welcome!")
	return nil
}

func TestHandler(ctx *gsweb.Context) error {
	ctx.Json(200, "open the door.")
	return nil
}
