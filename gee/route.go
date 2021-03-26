package gee

import (
	"fmt"
	"net/http"
	"strings"
)

type Route struct {
	root    map[string]*Node
	handler map[string]HandleFunc
}

func NewRoute() *Route {
	return &Route{
		handler: make(map[string]HandleFunc),
	}
}

func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}

	return parts
}

func (route *Route) addRoute(method string, pattern string, handleFunc HandleFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern

	if _, ok := route.root[method]; !ok {
		route.root = make(map[string]*Node)
		route.root[method] = &Node{}
	}

	route.root[method].insert(pattern, parts, 0)
	route.handler[key] = handleFunc
}

func (r *Route) getRoute(method, path string) (*Node, map[string]string) {
	searchParts := parsePattern(path)

	node, ok := r.root[method]

	if !ok {
		return nil, nil
	}

	node = node.search(searchParts, 0)

	params := make(map[string]string)
	if node != nil {
		parts := parsePattern(node.pattern)
		for index, part := range parts {
			fmt.Println(part)
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}

			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
			}
		}

		return node, params
	}

	return nil, nil

}

func (route *Route) handle(c *Context) {
	n, params := route.getRoute(c.Method, c.Path)

	if n != nil {
		key := c.Method + "-" + n.pattern
		c.Params = params
		handle := route.handler[key]
		handle(c)
	} else {
		c.String(http.StatusNotFound, "404 not found %s", c.Path)
	}

}
