package test

import (
	"GreensOne/gsweb"
	"time"
)

func Register(g *gsweb.GreensCore) {
	g.Get("/welcome", TestH1(), TestCloseController)
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
	ctx.Json("welcome!").SetStatus(200)
	return nil
}

func TestHandler(ctx *gsweb.Context) error {
	ctx.Json("open the door.").SetStatus(200)
	return nil
}

func TestCloseController(c *gsweb.Context) error {
	hhh, _ := c.QueryString("hhh", "okk")
	time.Sleep(10 * time.Second)
	c.SetOkStatus().Json("test close hhhhhhhhh " + hhh)
	return nil
}
