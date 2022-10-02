package gowebserver

import (
	"fmt"
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type router struct {
	handlers map[string]HandlerFunc
	roots    map[string]*node
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc), roots: make(map[string]*node)}
}

func parsePattern(pattern string) []string {
	v := strings.Split(pattern, "/")
	parts := make([]string, 0)

	for _, item := range v {
		if item != "" {
			parts = append(parts, item)

			if item[0] == '*' {
				break
			}
		}
	}

	return parts

}

func (r *router) handle(c *Context) {

	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[fmt.Sprintf("%s-%s", c.Method, n.pattern)])
		r.handlers[fmt.Sprintf("%s-%s", c.Method, n.pattern)](c)
	} else {

		c.String(http.StatusNotFound, "not found")
	}
	c.Next()
}

func (r *router) addRouter(method, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[fmt.Sprintf("%s-%s", method, pattern)] = handler

}

func (r *router) getRoute(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)

	params := make(map[string]string)

	root, ok := r.roots[method]

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

type node struct {
	pattern  string
	part     string
	children []*node
	isWild   bool
}

func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}

	return nil
}

func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)

	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}

	return nodes

}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]

	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}

	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildren(part)

	for _, child := range children {
		if result := child.search(parts, height+1); result != nil {
			return result
		}
	}
	return nil
}
