package gee

import (
	"log"
	"net/http"
)

type HandlerFunc func(*Context)

type Engine struct {
	*RouterGroup         //routergroup 是父类
	router       *router //key: Method-pattern value:HandlerFunc
	groups       []*RouterGroup
}

type RouterGroup struct {
	prefix      string
	middlewares HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

// New is the constructor of gee.Engine
func New() *Engine {
	engine := &Engine{router: NewRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// Group is defined to create a new RouterGroup
// remember all groups share the same Engine instance
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
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
