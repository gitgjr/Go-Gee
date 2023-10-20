package main

import (
	"fmt"
	"gee"
	"net/http"
)

func main() {
	server := gee.New()
	server.GET("/", indexHandler)
	server.GET("/hello", helloHandler)
	server.Run("localhost:8080")
}

// handler echoes r.URL.Path
func indexHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
}

// handler echoes r.URL.Header
func helloHandler(w http.ResponseWriter, req *http.Request) {
	for k, v := range req.Header {
		fmt.Fprintf(w, "Header[%q] = %q\n", k, v)
	}
}
