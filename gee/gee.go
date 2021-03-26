package gee

import (
	"net/http"
)

type HandleFunc func(*Context)

type Engine struct {
	Route *Route
}

func New() *Engine {
	return &Engine{Route: NewRoute()}
}

func (engine *Engine) GET(pattern string, handleFunc HandleFunc) {
	engine.Route.addRoute("GET", pattern, handleFunc)
}

func (engine *Engine) POST(pattern string, handleFunc HandleFunc) {
	engine.Route.addRoute("POST", pattern, handleFunc)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	context := NewContext(w, req)
	engine.Route.handle(context)
}
