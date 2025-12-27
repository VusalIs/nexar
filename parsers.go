package nexar

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)


type Parsers struct {}

func (p Parsers) parseRequest(reader *bufio.Reader) (*request, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error while parsing request line: ", err.Error())

		return nil, err
	}

	line = strings.TrimSpace(line)
	requestLineArr := strings.SplitN(line, " ", 3)
	if len(requestLineArr) != 3 {
		return nil, fmt.Errorf("malformed request line")
	}
	req := &request{
		method: requestLineArr[0],
		target: requestLineArr[1][1:],
		protocol: requestLineArr[2],
		Headers: make(map[string]string),
	}

	for {
		header, err := reader.ReadString('\n') 
		if err != nil {
			fmt.Println("Skipping problematic header")
		}
		if header == "\r\n" {
			break
		}
		
		headerKey, headerValue, found := strings.Cut(header, ":")
		headerKey = strings.TrimSpace(strings.ToLower(headerKey))
		if !found {
			fmt.Println("Header wasn't constructed properly, so skipping: ", headerKey)
			continue
		} else {
			req.Headers[headerKey] = strings.TrimSpace(headerValue)
		}
	}
	
	if contentLengthSt, ok := req.Headers["content-length"]; ok {
		if !ok {
			fmt.Println("Missing Content-Length so there is no content body")
	
			return req, nil
		} 

		contentLength, err := strconv.Atoi(contentLengthSt)
		if err != nil {
			return nil, fmt.Errorf("invalid content-length")
		}

		content := make([]byte, contentLength)
	
		if _, err = io.ReadFull(reader, content); err != nil {
			fmt.Println("Error while reding the content body: ", err.Error())
	
			return req, nil
		}

		req.Body = content
	}

	return req, nil
}

func (p Parsers) parseResponse(res *response) []byte {
	var stringB strings.Builder

	stringB.Write([]byte(fmt.Sprintf("%s %s %s\r\n", res.protocol, res.code, res.status)))

	for key, value := range res.headers {
		stringB.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
	}
	stringB.Write([]byte("\r\n"))
	stringB.Write(res.body)

	return []byte(stringB.String())
}
