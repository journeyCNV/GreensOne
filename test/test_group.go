package test

import (
	"GreensOne/gsweb"
)

func Register(g *gsweb.GreensCore) {
	group := g.Group("/hhhh")
	group1 := group.Group("/okk1")
	group1.Get("/ohhhh", TestHandler)
}

func TestHandler(c *gsweb.Context) error {
	return nil
}
