package main

import (
	"gee"
	"log"
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
	g.GET("/hello/:username", func(c *gee.Context) {
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("username"), c.Path)
	})

	g.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.Json{
			"filepath": c.Param("filepath"),
		})
	})

	log.Fatal(g.Run("0.0.0.0:9999"))
}
