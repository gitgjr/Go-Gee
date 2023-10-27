package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{} //JSON map,interface means any data struct

type Context struct {
	// origin objects
	Writer http.ResponseWriter
	Req    *http.Request

	//Req info
	Path   string
	Method string
	Params map[string]string
	//response info
	StatusCode int
}

func (content *Context) Param(key string) string {
	value, _ := content.Params[key]
	return value
}

func NewContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

func (context *Context) PostForm(key string) string {
	return context.Req.FormValue(key)
}

func (context *Context) Query(key string) string {
	return context.Req.URL.Query().Get(key)
}

func (context *Context) Status(statusCode int) {
	context.StatusCode = statusCode
	context.Writer.WriteHeader(statusCode)
}

func (context *Context) SetHeader(key string, value string) {
	context.Writer.Header().Set(key, value)
}

func (context *Context) String(statusCode int, format string, values ...interface{}) {
	context.SetHeader("Content-Type", "text/plain")
	context.Status(statusCode)
	context.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (context *Context) JSON(statusCode int, obj interface{}) {
	context.SetHeader("Content-Type", "application/json")
	context.Status(statusCode)
	encoder := json.NewEncoder(context.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(context.Writer, err.Error(), 500)
	}
}

func (context *Context) Data(statusCode int, data []byte) {
	context.Status(statusCode)
	context.Writer.Write(data)
}

func (context *Context) HTML(statusCode int, html string) {
	context.SetHeader("Content-Type", "text/html")
	context.Status(statusCode)
	context.Writer.Write([]byte(html))
}
