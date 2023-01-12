package gee

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	http.Handler                        // implement interface
	router       map[string]HandlerFunc // record route rules
}

func NewEngine() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// addRouter is a private method.
func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

// GET is called from user to add GET router
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRouter("GET", pattern, handler)
}

// POST is called from user to add POST router
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRouter("POST", pattern, handler)
}

// Run will listen addr and serve service
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// implement http.Handler interface
func (engine *Engine) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(res, req)
	} else {
		res.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(res, "404 NOT FOUND: %q\n", req.URL)
	}
}
