package gee

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node       // eg: roots["GET"] root["POST"]
	handlers map[string]HandlerFunc // eg: handlers["GET-/p/:lang/doc"], handlers["POST-/p/book"]
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	patternItems := strings.Split(pattern, "/")

	parts := make([]string, 0)

	for _, patternItem := range patternItems {
		if patternItem != "" {
			parts = append(parts, patternItem)
			// only one '*' is permitted
			if patternItem[0] == '*' {
				break
			}
		}
	}
	return parts
}

// addRouter is a private method and defines the method to add handler.
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	log.Printf("Route %4s - %4s will be added\n", method, pattern)
	parts := parsePattern(pattern)
	if _, ok := r.roots[method]; ok == false {
		// 不是线程安全的
		// new(node) 与 &node{} 等价
		r.roots[method] = new(node)
	}
	r.roots[method].insert(pattern, parts, 0)
	key := method + "-" + pattern

	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	pathItems := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if ok == false {
		return nil, nil
	}

	n := root.search(pathItems, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				// /p/go/doc匹配到/p/:lang/doc
				// params为{lang: "go"}
				params[part[1:]] = pathItems[index]
			} else if part[0] == '*' && len(part) > 1 {
				// /static/css/geektutu.css匹配到/static/*filepath
				// param为{filepath: "css/geektutu.css"}
				params[part[1:]] = strings.Join(pathItems[index:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		r.handlers[key](c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUNT: %s\n", c.Path)
	}
}
