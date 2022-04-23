package test

import (
	"GreensOne/gsweb"
)

func Register(g *gsweb.GreensCore) {
	g.Get("/welcome", WelcomeHandler)
	group := g.Group("/show")
	{
		group1 := group.Group("/door")
		{
			group1.Get("/open", TestHandler)
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
