package nexar

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
    MB = 1 << 20
)

func serializeRequest(reader *bufio.Reader) (*Request, error) {
    line, err := reader.ReadString('\n')
    if err != nil {
        return nil, fmt.Errorf("reading request line: %w", err)
    }

    parts := strings.SplitN(strings.TrimSpace(line), " ", 3)
    if len(parts) != 3 {
        return nil, fmt.Errorf("malformed request line")
    }

    target := parts[1]
    if len(target) > 0 && target[0] == '/' {
        target = target[1:]
    }

    req := &Request{
        Method:   parts[0],
        Target:   target,
        Protocol: parts[2],
        Headers:  make(map[string]string),
    }

    // Parse headers
    for {
        header, err := reader.ReadString('\n')
        if err != nil {
            return nil, fmt.Errorf("reading header: %w", err)
        }
        if header == "\r\n" {
            break
        }

        key, value, found := strings.Cut(header, ":")
        if !found {
            continue
        }
        req.Headers[strings.ToLower(strings.TrimSpace(key))] = strings.TrimSpace(value)
    }

    if contentLengthStr, ok := req.Headers["content-length"]; ok {
        contentLength, err := strconv.Atoi(contentLengthStr)
        if err != nil {
            return nil, fmt.Errorf("invalid content-length: %w", err)
        }
		// 10 MB limit
        if contentLength > 10 * MB { 
            return nil, fmt.Errorf("content-length too large")
        }

        req.Body = make([]byte, contentLength)
        if _, err = io.ReadFull(reader, req.Body); err != nil {
            return nil, fmt.Errorf("reading body: %w", err)
        }
    }

    return req, nil
}

func formatResponse(res *response) []byte {
    var buf bytes.Buffer
    
    // Status line
    buf.WriteString(res.protocol)
    buf.WriteByte(' ')
    buf.WriteString(res.code)
    buf.WriteByte(' ')
    buf.WriteString(res.status)
    buf.WriteString("\r\n")
    
    // Headers
    for key, value := range res.headers {
        buf.WriteString(key)
        buf.WriteString(": ")
        buf.WriteString(value)
        buf.WriteString("\r\n")
    }
    
    // Blank line + body
    buf.WriteString("\r\n")
    buf.Write(res.body)
    
    return buf.Bytes()  // no conversion needed
}
