package gee

import (
	"net/http"
	"strings"
)

type HandleFunc func(*Context)

type (
	Engine struct {
		Route *Route
		*routeGrounp
		routeGrounps []*routeGrounp
	}

	routeGrounp struct {
		prefix string
		parent *routeGrounp
		engine *Engine
		middle []HandleFunc
	}
)

func New() *Engine {
	engine := &Engine{Route: NewRoute()}
	engine.routeGrounp = &routeGrounp{engine: engine}
	engine.routeGrounps = []*routeGrounp{engine.routeGrounp}
	engine.Use(Recoverier())
	return engine
}

//添加中间件
func (g *routeGrounp) Use(handleFunc ...HandleFunc) {
	g.middle = append(g.middle, handleFunc...)
}

func (g *routeGrounp) Group(prefix string) *routeGrounp {
	engine := g.engine
	newGroup := &routeGrounp{
		prefix: g.prefix + prefix,
		parent: g,
		engine: engine,
		middle: make([]HandleFunc, 0),
	}
	engine.routeGrounps = append(engine.routeGrounps, newGroup)

	return newGroup
}

func (g *routeGrounp) GET(pattern string, handleFunc HandleFunc) {
	g.addRoute("GET", pattern, handleFunc)
}

func (g *routeGrounp) POST(pattern string, handleFunc HandleFunc) {

	g.addRoute("POST", pattern, handleFunc)
}

func (g *routeGrounp) addRoute(method string, pattern string, handleFunc HandleFunc) {
	pattern = g.prefix + pattern
	g.engine.Route.addRoute(method, pattern, handleFunc)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

//处理http请求
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandleFunc
	for _, group := range engine.routeGrounps {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middle...)
		}
	}

	context := NewContext(w, req)
	context.middle = middlewares
	engine.Route.handle(context)
}
