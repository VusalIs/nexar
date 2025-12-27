package nexar

import (
	"encoding/json"
	"strconv"
	"strings"
)

type Context struct {
	Request  *Request
	Param    map[string]string
	Response *response
	Config   *Config
}

func newContext(params map[string]string, request *Request, config *Config) *Context {
	return &Context{
		Request: request,
		Param:   params,
		Config:  config,
		Response: &response{
			protocol: "HTTP/1.1",
			headers:  make(map[string]string),
		},
	}
}

func (c *Context) applyEncoding(acceptedEncoding string) error {
    encodingHeader, ok := c.Request.Headers["accept-encoding"]
    if !ok {
        return nil
    }

    for _, enc := range strings.Split(encodingHeader, ",") {
        if strings.TrimSpace(enc) == acceptedEncoding {
            encoded, err := encodeString(c.Response.body)
            if err != nil {
                return err
            }
            c.Response.body = encoded
            c.Response.headers["Content-Encoding"] = acceptedEncoding
            return nil
        }
    }
    return nil
}

func (c *Context) finalize() {
    c.Header("Content-Length", strconv.Itoa(len(c.Response.body)))
    
    if conn, ok := c.Request.Headers["connection"]; ok && conn == "close" {
        c.Response.headers["Connection"] = "close"
    }
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


