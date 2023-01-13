package main

import (
	"gee"
	"net/http"
)

func main() {
	g := gee.NewEngine()
	g.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>root path</h1>")
	})
	g.GET("/hello", func(c *gee.Context) {
		// expect /hello?username=ning
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("username"), c.Path)
	})
	g.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.Json{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	g.Run("0.0.0.0:9999")
}
