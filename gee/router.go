package gee

import (
	"fmt"
	"strings"
)

type Router struct {
	roots    map[string]*Node
	handlers map[string]HandlerFunc //key: Method-URL value handler function
}

func NewRouter() *Router {
	return &Router{
		roots:    make(map[string]*Node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	patternArray := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range patternArray {
		if item != "" {
			parts = append(parts, item)
			if item == "*" {
				break
			}
		}

	}
	fmt.Println(parts)
	return parts
}

func (router *Router) addRoute(method string, pattern string, handler HandlerFunc) {

	parts := parsePattern(pattern)
	//log.Printf("Route %4s - %s", method, pattern)
	key := method + "-" + pattern
	_, ok := router.roots[method]
	if !ok {
		router.roots[method] = &Node{}
	}
	router.roots[method].insert(pattern, parts, 0)
	router.handlers[key] = handler
}

func (router *Router) getRoute(method string, pattern string) (*Node, map[string]string) {
	parts := parsePattern(pattern)
	params := make(map[string]string)
	root, ok := router.roots[method]

	if !ok {
		return nil, nil
	}
	node := root.search(parts, 0)
	if node != nil {
		parts := parsePattern(node.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = parts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(parts[index:], "/")
				break
			}
		}
		return node, params
	}
	return nil, nil
} //TODO can not understand

func (router *Router) Handle(context *Context) {
	node, params := router.getRoute(context.Method, context.Path)

	if node != nil {
		context.Params = params
		key := context.Method + "-" + node.pattern
		handler, _ := router.handlers[key]
		handler(context)
	} else {
		fmt.Fprint(context.Writer, "404 page not found: ", context.Req.URL.Path)
	}

}
