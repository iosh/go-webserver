package gowebserver

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	Write      http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
	Params     map[string]string
	handlers   []HandlerFunc
	index      int
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Write:  w,
		Req:    r,
		Path:   r.URL.Path,
		Method: r.Method,
		index:  -1,
	}
}
func (c *Context) Next() {
	c.index++

	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}

}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) Query(key string) string {

	return c.Req.URL.Query().Get(key)

}

func (c *Context) PostForm(key string) string {

	return c.Req.FormValue(key)

}

func (c *Context) Status(code int) {

	c.StatusCode = code
	c.Write.WriteHeader(code)

}

func (c *Context) SetHeader(key, value string) {

	c.Write.Header().Set(key, value)

}

func (c *Context) String(code int, s string) {

	c.SetHeader("Content-type", "text/plain")
	c.Status(code)
	c.Write.Write([]byte(s))

}

func (c *Context) JSON(code int, m map[string]any) {

	c.SetHeader("Content-type", "application/json")
	c.Status(code)

	if err := json.NewEncoder(c.Write).Encode(m); err != nil {
		http.Error(c.Write, err.Error(), 500)
	}

}

func (c *Context) Data(code int, data []byte) {

	c.Status(code)
	c.Write.Write(data)

}

func (c *Context) HTML(code int, html string) {

	c.SetHeader("Content-type", "text/html")
	c.Status(code)
	c.Write.Write([]byte(html))

}
