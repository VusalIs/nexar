package nexar

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"net"
	"os"
	"strings"
)

func errInternalServer() *response {
	return &response{
		protocol: "HTTP/1.1",
		code:     "500",
		status:   "Internal Server Error",
		headers:  make(map[string]string),
	}
}

func errNotFound() *response {
	return &response{
		protocol: "HTTP/1.1",
		code:     "404",
		status:   "Not Found",
		headers:  make(map[string]string),
	}
}
type Request struct {
	Method   string
	Target   string
	Protocol string
	Headers  map[string]string
	Body     []byte
}

type response struct {
	status   string
	protocol string
	code     string
	headers  map[string]string
	body     []byte
}

type Config struct {
	Directory *string
	AcceptedEncoding string
}

type Nexar struct{
	port string
	tree *Tree
	config *Config
}

func Default(config *Config) *Nexar {
	return &Nexar{
		tree: New(),
		port: "8080",
		config: config,
	}
}

func (n *Nexar) Get(route string, fn func(cntx *Context) *Context) {
	n.tree.AddNode(append([]string{"GET"}, strings.Split(route, "/")...), fn)
}

func (n *Nexar) Post(route string, fn func(cntx *Context) *Context) {
	n.tree.AddNode(append([]string{"POST"}, strings.Split(route, "/")...), fn)
}

func (n *Nexar) Run(port string) error {
	l, err := net.Listen("tcp", "0.0.0.0:" + port)
	if err != nil {
		return fmt.Errorf("failed to bind to port %s: %w", port, err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		
		go engine(n, conn)
	}
}

func engine(nexar *Nexar, conn net.Conn) {
		reader := bufio.NewReader(conn)

		for {
			request, err := parseRequest(reader)
			if request == nil {
				conn.Close()
				break
			}
			if err != nil {
				fmt.Println("Error while parsing the request: ", err.Error())
				conn.Write(parseResponse(errInternalServer()))
	
				continue
			}
		
			treeNode, params := nexar.tree.FindNodeByRoute(request.Method + "/" + request.Target)
		
			if treeNode == nil {
				conn.Write(parseResponse(errNotFound()))
				continue
			}
		
			cntx := newContext(params, request, nexar.config)
			
			treeNode.handler(cntx)
		
			if err := cntx.applyEncoding(nexar.config.AcceptedEncoding); err != nil {
				fmt.Println("Error while encoding the response body")
				cntx.Response = errInternalServer()
			}
			cntx.finalize()
			
			conn.Write(parseResponse(cntx.Response))
	
			if shouldClose(cntx.Request) {
				conn.Close()
				break
			}
		}
}

func shouldClose(req *Request) bool {
    conn, ok := req.Headers["connection"]
    return ok && conn == "close"
}

func gzipCompress(dt []byte) ([]byte, error) {
	var buf bytes.Buffer

	gz := gzip.NewWriter(&buf)
	if _, err := gz.Write(dt); err != nil {
		gz.Close()

		return nil, err
	}
	// TODO: Learn more about why close can return an error
	if err := gz.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}