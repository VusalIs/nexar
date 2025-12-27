package nexar

import (
	"encoding/json"
)

type Context struct {
	Request *request
	directory *string
	Param map[string]string
	Response *response
	Config *Config
}

func(c *Context)Init(params map[string]string, directory *string, request *request) {
	c.Response = &response{
		protocol: "HTTP/1.1",
		headers: make(map[string]string),
	}
	c.Param = params
	c.directory = directory
	c.Request = request
}

func(c *Context) Header(key string, vl string) {
	c.Response.headers[key] = vl
}

func(c *Context) Data(status int, dt []byte) {
	c.Status(200)

	c.Response.body = dt
}

func (c *Context) JSON(status int, bd any) {
	c.Status(status)

	body, err := json.Marshal(bd)
	if err != nil {
		panic("bd can't not be marshaled")
	}

	c.Response.body = body
}

func(c *Context) Status(status int) {
	switch status {
		case 200:	
			c.Response.code = "200"
			c.Response.status = "OK"
		case 201:
			c.Response.code = "201"
			c.Response.status = "Created"
		case 404:
			c.Response.code = "404"
			c.Response.status = "Not Found"
		default:
			panic("Status code is invalid")

	}
}


