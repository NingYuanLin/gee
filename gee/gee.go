package gee

import (
	"html/template"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(c *Context)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // 中间件
	engine      *Engine       // 所有的group共享同一个engine实例
}

type Engine struct {
	http.Handler         // implement interface
	*RouterGroup         // 组合 composition
	router       *router // record route rules
	groups       []*RouterGroup
	// html render
	htmlTemplates *template.Template
	funcMap       template.FuncMap
}

func New() *Engine {
	engine := &Engine{}
	engine.RouterGroup = &RouterGroup{
		engine: engine,
	}
	engine.router = newRouter()
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	// fs may be: http.Dir(rootPath)
	// assume: relativePath: /v1; group.prefix: /assets; rootPath: /root/statics
	// absolutePath: /v1/assets
	absolutePath := path.Join(group.prefix, relativePath)
	// remove prefix /v1/assets from req.URL.Path
	// /v1/assets/file.txt => /file.txt
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		// file.txt
		filepath := c.Params["filepath"]
		// check if file exits, and we have permission to access it
		// fs.Open(filepath) = /root/statics/file.txt
		if _, err := fs.Open(filepath); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Static is defined to serve static files
// relativePath: the URL prefix you want to use. such as /assets/statics
// root: dir path in local systems. such as /root/statics/
func (group *RouterGroup) Static(method string, relativePath string, root string) {
	staticHandler := group.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	group.addRoute(method, urlPattern, staticHandler)
}

func (group *RouterGroup) AddMiddleware(middleware ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middleware...)
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
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// addRoute is a private method and defines the method to add router.
// comp: Composition
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
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

func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

func (engine *Engine) LoadHTMLGlob(pattern string) {
	/*
		pattern may be "templates/*"
		分析：
		1. template.Must() 让template对象的加载，如果加载不到，就产生panic
		2. template.New() 产生template对象
		3. Funcs() 添加模板函数，engine中的funcMap是个map类型里面能够保存多个模板函数
		4. ParseGlob() 将pattern中的文件全都读取出来，按照name，存入engine.htmlTemplates中（每个里面都携带了自定义模板函数）
	*/
	engine.htmlTemplates = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
}

// implement http.Handler interface
func (engine *Engine) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	// TODO: 性能问题
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			// TODO: 这里的顺序没办法很好地去控制, 只能根据AddMiddleware函数调用的顺序来进行
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	context := NewContext(res, req)
	context.handlers = middlewares
	context.engine = engine
	engine.router.handle(context)
}
