package nexar

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)



type request struct {
	method string
	target string
	protocol string
	headers map[string]string
	body []byte
}

type response struct {
	status string
	protocol string
	code string
	headers map[string]string
	body []byte
}

type Config struct {
	directory *string
}

type Nexar struct{
	port string
	tree *Tree
	config *Config
}

func (n *Nexar) Default() *Nexar {
	return &Nexar{
		tree: New(),
		port: "8080",
	}
}

func (n *Nexar) Get(route string, fn func(cntx *Context) *Context) {
	n.tree.AddNode(append([]string{"GET"}, strings.Split(route, "/")...), fn)
}

func (n *Nexar) Post(route string, fn func(cntx *Context) *Context) {
	n.tree.AddNode(append([]string{"POST"}, strings.Split(route, "/")...), fn)
}

func (n *Nexar) Run(port string) {
	l, err := net.Listen("tcp", "0.0.0.0:" + port)
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
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
	parsers := Parsers{}

	request, err := parsers.parseRequest(reader)
	if err != nil {
		fmt.Println("Error while parsing the request: ", err.Error())

		conn.Write(parsers.parseResponse(&response{
			protocol: "HTTP/1.1",
			status: "Internal Problem",
			code: "500",
		}))
	}

	treeNode, params := nexar.tree.FindNodeByRoute(request.method + request.target)

	cntx := &Context{
		Config: nexar.config,
	}
	cntx.Init(params, cntx.Config.directory, cntx.Request)

	if treeNode == nil {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))

		return
	}
	
	treeNode.handler(cntx)

	conn.Write(parsers.parseResponse(cntx.Response))
}