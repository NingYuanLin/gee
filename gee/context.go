package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Json map[string]any

type Context struct {
	// origin object
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int // default:-1
}

func NewContext(writer http.ResponseWriter, Req *http.Request) *Context {
	return &Context{
		Writer: writer,
		Req:    Req,
		Path:   Req.URL.Path,
		Method: Req.Method,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	for ; c.index < len(c.handlers); c.index++ {
		c.handlers[c.index](c)
	}
}

// PostForm defines the method to get form message. Be careful that only POST method have form message.
func (c *Context) PostForm(key string) string {
	//return c.Req.Form.Get(key) // 不能调用这个
	return c.Req.FormValue(key)
}

// Query defines the method to get query message
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key] // string类型不可能为nil，当为空时，是""
	return value
}

// Status defines the method to set status code
func (c *Context) Status(code int) {
	// 避免重复调用
	// TODO: 需要优化逻辑
	if c.StatusCode == 0 {
		c.StatusCode = code
		c.Writer.WriteHeader(code)
	}
}

// SetHeader defines the method to set response header
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) SetContentType(value string) {
	c.SetHeader("Content-Type", value)
}

// String defines the method to set response string message
func (c *Context) String(code int, format string, values ...any) {
	c.SetContentType("text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON defines the method to set response json message
func (c *Context) JSON(code int, obj Json) {
	c.SetContentType("application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// Data defines the method to set response data message
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

// HTML defines the method to set response html message
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}
