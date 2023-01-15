package main

import (
	"gee"
	"log"
	"net/http"
)

func main() {
	g := gee.New()

	g.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>index page</h1>")
	})

	v1 := g.Group("/v1")
	{
		v1.GET("/hello", func(c *gee.Context) {
			c.HTML(http.StatusOK, "hello@v1\n")
		})
		v1.GET("/hello/:username", func(c *gee.Context) {
			// expect /hello/ning
			c.String(http.StatusOK, "hello %s, you're at %s @v1\n", c.Param("username"), c.Path)
		})
	}
	v2 := g.Group("/v2")
	{
		v2.GET("/hello", func(c *gee.Context) {
			c.HTML(http.StatusOK, "hello@v2\n")
		})
		v2.GET("/hello/:username", func(c *gee.Context) {
			// expect /hello/ning
			c.String(http.StatusOK, "hello %s, you're at %s @v2\n", c.Param("username"), c.Path)
		})
	}

	log.Fatal(g.Run("0.0.0.0:9999"))
}
