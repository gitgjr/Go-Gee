package gee

import (
	"net/http"
)

type HandlerFunc func(*Context)

type Engine struct {
	router *Router //key: Method-pattern value:HandlerFunc
}

func New() *Engine {
	return &Engine{router: NewRouter()}
}
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}
func (engine *Engine) Run(addr string) {
	http.ListenAndServe(addr, engine)
}
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	context := NewContext(w, req)
	engine.router.Handle(context)
}

//func (engine *Engine) ServeHTTP(w http.ResponseWriter, Req *http.Request) {
//	switch Req.URL.Path {
//	case "/":
//		fmt.Fprint(w, "welcome", Req.URL.Path)
//	case "/img":
//		for k, v := range Req.Header {
//			fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
//		}
//	default:
//		fmt.Fprint(w, "404")
//	}
//
//}
