package gowebserver

import (
	"net/http"
	"strings"
)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	webServer   *WebServer
}

type WebServer struct {
	router *router
	*RouterGroup
	groups []*RouterGroup
}

func New() *WebServer {
	server := &WebServer{router: newRouter()}
	server.RouterGroup = &RouterGroup{webServer: server}
	server.groups = []*RouterGroup{server.RouterGroup}

	return server
}

func (s *WebServer) Run(add string) error {

	return http.ListenAndServe(add, s)

}

func (g *RouterGroup) Group(prefix string) *RouterGroup {

	server := g.webServer

	newGroup := &RouterGroup{
		prefix:    g.prefix + prefix,
		parent:    g,
		webServer: server,
	}

	server.groups = append(server.groups, newGroup)

	return newGroup

}

func (g *RouterGroup) addRouter(method, pattern string, handler HandlerFunc) {

	p := g.prefix + pattern

	g.webServer.router.addRouter(method, p, handler)

}

func (g *RouterGroup) GET(pattern string, handler HandlerFunc) {
	g.addRouter(http.MethodGet, pattern, handler)
}

func (g *RouterGroup) POST(pattern string, handler HandlerFunc) {
	g.addRouter(http.MethodPost, pattern, handler)
}

func (g *RouterGroup) Use(middleware ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middleware...)
}
func (s *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var middlewares []HandlerFunc

	for _, group := range s.groups {
		if strings.HasPrefix(r.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := NewContext(w, r)
	c.handlers = middlewares
	s.router.handle(c)
}
