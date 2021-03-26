package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	W   http.ResponseWriter
	Req *http.Request

	Path   string
	Method string
	Params map[string]string

	StatusCode int
}

func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		W:      w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

//获取post表单字段
func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

//获取get字段
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

//设置状态码,应该在设置完头部后调用
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.W.WriteHeader(code)
}

//设置头部
func (c *Context) SetHeader(key string, value string) {
	c.W.Header().Set(key, value)
}

//返回json格式
func (c *Context) Json(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.W)
	if err := encoder.Encode(obj); err != nil {
		panic(err)
	}

}

func (c *Context) GetParams(key string) string {
	value := c.Params[key]
	return value
}

//返回html
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.W.Write([]byte(html))
}

//格式化输出文本
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.W.Write([]byte(fmt.Sprintf(format, values...)))
}

//直接返回数据
