package main

import (
	"fmt"
	"gee"
	"html/template"
	"log"
	"net/http"
	"time"
)

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	g := gee.New()
	g.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	g.LoadHTMLGlob("templates/*")
	//g.Static("GET", "/assets", "./statics")
	// or
	g.Group("/assets").Static("GET", "", "./statics")

	g2 := g.Group("/html")

	g2.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	g2.GET("/date", func(c *gee.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gee.Json{
			"title": "ning",
			"now":   time.Date(2023, 1, 17, 00, 18, 00, 00, time.UTC),
		})
	})

	log.Fatal(g.Run("0.0.0.0:9999"))
}
