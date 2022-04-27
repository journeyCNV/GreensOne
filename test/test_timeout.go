package test

import (
	"GreensOne/gsweb"
	"fmt"
)

func TestH1() gsweb.HandlerFunc {
	return func(c *gsweb.Context) error {
		fmt.Println("mid 1")
		c.Next()
		fmt.Println("mid 1 ")
		return nil
	}
}

func TestH2() gsweb.HandlerFunc {
	return func(c *gsweb.Context) error {
		fmt.Println("mid 2")
		c.Next()
		fmt.Println("mid 2 ")
		return nil
	}
}
