# Nexar

A lightweight HTTP framework for Go, built from scratch without using `net/http`.

## Why I Built This

I wanted to understand how web frameworks work under the hood. Instead of relying on Go's standard HTTP library, I implemented the core components myself:

- **TCP connection handling** — Raw socket management with keep-alive support
- **HTTP parsing** — Request/response serialization following the HTTP/1.1 spec
- **Radix tree router** — Fast route matching with dynamic path parameters
- **Gzip compression** — Content encoding based on client preferences

## Features

- Clean, Express-style API
- Path parameters (`:id`, `:name`)
- JSON response helpers
- Request body parsing with size limits
- Connection keep-alive
- Gzip response compression

## Example

```go
package main

import "github.com/yourusername/nexar"

func main() {
    app := nexar.New(&nexar.Config{
        AcceptedEncoding: "gzip",
    })

    app.Get("/hello", func(ctx *nexar.Context) *nexar.Context {
        return ctx.JSON(200, map[string]string{
            "message": "Hello, World!",
        })
    })

    app.Get("/users/:id", func(ctx *nexar.Context) *nexar.Context {
        userID := ctx.Param["id"]
        return ctx.JSON(200, map[string]string{
            "user": userID,
        })
    })

    app.Run("8080")
}
```
