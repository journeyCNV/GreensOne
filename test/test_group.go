package test

import (
	"github.com/journeycnv/greensone/gsweb"
	"github.com/journeycnv/greensone/gsweb/middleware"
	"time"
)

func Register(g *gsweb.GreensCore) {
	g.Get("/welcome", TestH1(), TestCloseController)
	group := g.Group("/show")
	{
		group1 := group.Group("/door")
		group1.Use(TestH1()) // 批量使用中间件
		{
			group1.Get("/open", middleware.TimeoutHandler(5), TestHandler)
		}
		group.Get("/happy", WelcomeHandler)
	}
}

func WelcomeHandler(ctx *gsweb.Context) error {
	gsweb.LogInfo("start", nil)
	service := ctx.MustMake(key).(DService)
	gsweb.LogInfo("here!", nil)
	fun := service.MustSmile()
	ctx.SetOkStatus().Json(fun)
	return nil
}

func TestHandler(ctx *gsweb.Context) error {
	gsweb.LogInfo("test 哇", &gsweb.LogField{
		"msg": "open the door",
	})
	ctx.Json("open the door.").SetStatus(200)
	return nil
}

func TestCloseController(c *gsweb.Context) error {
	hhh, _ := c.QueryString("hhh", "okk")
	time.Sleep(10 * time.Second)
	c.SetOkStatus().Json("test close hhhhhhhhh " + hhh)
	return nil
}
