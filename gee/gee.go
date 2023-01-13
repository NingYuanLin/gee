package gee

import (
	"net/http"
)

type HandlerFunc func(c *Context)

type Engine struct {
	http.Handler         // implement interface
	router       *router // record route rules
}

func NewEngine() *Engine {
	return &Engine{router: NewRouter()}
}

// addRouter is a private method and defines the method to add router.
func (engine *Engine) addRouter(method string, pattern string, handler HandlerFunc) {
	engine.router.addRouter(method, pattern, handler)
}

// GET defines the method to add get router
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRouter("GET", pattern, handler)
}

// POST defines the method to add post router
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRouter("POST", pattern, handler)
}

// Run defines the method to run the engine
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// implement http.Handler interface
func (engine *Engine) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	context := NewContext(res, req)
	engine.router.handle(context)
}
