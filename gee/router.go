package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc //key: Method-URL value handler function
}

func NewRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	patternArray := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range patternArray {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}

	}

	return parts
}

func (router *router) addRoute(method string, pattern string, handler HandlerFunc) {

	parts := parsePattern(pattern)
	//log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	_, ok := router.roots[method]
	if !ok {
		router.roots[method] = &node{}
	}
	router.roots[method].insert(pattern, parts, 0)
	router.handlers[key] = handler
}

func (router *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := router.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return n, params
	}

	return nil, nil
}

func (router *router) Handle(context *Context) {
	node, params := router.getRoute(context.Method, context.Path)

	if node != nil {
		context.Params = params
		key := context.Method + "-" + node.pattern
		handler, _ := router.handlers[key]
		handler(context)
	} else {
		context.String(http.StatusNotFound, "404 NOT FOUND: %s\n", context.Path)
	}

}
