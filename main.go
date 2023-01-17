package main

import (
	"gee"
	"log"
	"net/http"
)

func main() {
	g := gee.New()

	g.AddMiddleware(gee.Recovery(func(c *gee.Context) {
		c.JSON(http.StatusInternalServerError, gee.Json{
			"message": "Internal Server Error",
		})
	}))

	g.GET("/", func(c *gee.Context) {
		c.String(http.StatusOK, "hello")
	})

	g.GET("/panic", func(c *gee.Context) {
		var temp []string
		c.String(http.StatusOK, temp[1])
	})

	log.Fatal(g.Run("0.0.0.0:9999"))
}
