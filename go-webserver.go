package gowebserver

import (
	"net/http"
)

type WebServer struct {
	router *router
}

func New() *WebServer {
	return &WebServer{router: newRouter()}
}

func (s *WebServer) Run(add string) error {

	return http.ListenAndServe(add, s)

}

func (s *WebServer) GET(pattern string, handler HandlerFunc) {
	s.router.addRouter(http.MethodGet, pattern, handler)
}

func (s *WebServer) POST(pattern string, handler HandlerFunc) {
	s.router.addRouter(http.MethodPost, pattern, handler)
}

func (s *WebServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	c := NewContext(w, r)

	s.router.handle(c)
}
