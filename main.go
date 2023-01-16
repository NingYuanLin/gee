package main

import (
	"gee"
	"log"
	"net/http"
)

func v2Middleware1() gee.HandlerFunc {
	return func(c *gee.Context) {
		c.String(http.StatusOK, "v2Middleware1 enter\n")
		c.Next()
		c.String(http.StatusOK, "v2Middleware1 exit\n")
	}
}

func v2Middleware2() gee.HandlerFunc {
	return func(c *gee.Context) {
		c.String(http.StatusOK, "v2Middleware2 enter\n")
		c.Next()
		c.String(http.StatusOK, "v2Middleware2 exit\n")
	}
}

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
	v2.AddMiddleware(v2Middleware1())
	v2.AddMiddleware(v2Middleware2())
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
