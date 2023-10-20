package gee

import (
	"fmt"
	"net/http"
)

type handlerFunc func(w http.ResponseWriter, req *http.Request)

type Engine struct {
	router map[string]handlerFunc //key: Method-pattern value:handlerFunc
}

func New() *Engine {
	return &Engine{router: map[string]handlerFunc{}}
}
func (engine *Engine) addRoute(method string, pattern string, handler handlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}
func (engine *Engine) GET(pattern string, handler handlerFunc) {
	engine.addRoute("GET", pattern, handler)
}
func (engine *Engine) POST(pattern string, handler handlerFunc) {
	engine.addRoute("POST", pattern, handler)
}
func (engine *Engine) Run(addr string) {
	http.ListenAndServe(addr, engine)
}
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	hanlder, ok := engine.router[key]
	if ok {
		hanlder(w, req)
	} else {
		fmt.Fprint(w, "404 page not found: ", req.URL.Path)
	}
}

//func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
//	switch req.URL.Path {
//	case "/":
//		fmt.Fprint(w, "welcome", req.URL.Path)
//	case "/img":
//		for k, v := range req.Header {
//			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
//		}
//	default:
//		fmt.Fprint(w, "404")
//	}
//
//}
