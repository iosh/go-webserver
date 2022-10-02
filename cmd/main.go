package main

import (
	"net/http"

	gowebserver "github.com/iosh/go-webserver"
)

func main() {
	r := gowebserver.New()

	r.GET("/index", func(c *gowebserver.Context) {
		c.String(http.StatusOK, "ok")
	})

	v1 := r.Group("/v1")

	v1.GET("/", func(c *gowebserver.Context) {
		c.HTML(http.StatusOK, "<h1>hello world</h1>")
	})
	r.Run(":9900")
}
