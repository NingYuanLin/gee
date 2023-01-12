package main

import (
	"fmt"
	"gee"
	"net/http"
)

func main() {
	g := gee.NewEngine()
	g.GET("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "URL.Path = %q\n", request.URL.Path)
	})
	g.GET("/hello", func(writer http.ResponseWriter, request *http.Request) {
		for k, v := range request.Header {
			fmt.Fprintf(writer, "Header[%q] = %s\n", k, v)
		}
	})
	g.Run("0.0.0.0:9999")
}
