package gee

import (
	"log"
	"net/http"
)

type HandlerFunc func(c *Context)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // 中间件
	//parent      *RouterGroup  // 支持嵌套
	engine *Engine // 所有的group共享同一个engine实例
}

type Engine struct {
	http.Handler         // implement interface
	*RouterGroup         // 组合 composition
	router       *router // record route rules
	//groups       []*RouterGroup
}

func New() *Engine {
	engine := &Engine{}
	engine.RouterGroup = &RouterGroup{
		engine: engine,
	}
	engine.router = newRouter()
	//engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group is defined to create a new group
// remember all groups share the same engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		//parent: group,
		engine: engine,
	}
	//engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// addRoute is a private method and defines the method to add router.
// comp: Composition
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %4s\n", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add get router
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST defines the method to add post router
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
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
