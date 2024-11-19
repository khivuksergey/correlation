# correlation

This Go library provides utilities for managing correlation IDs, which are used for tracking and tracing requests across different services in distributed systems. The library allows you to define and manage the correlation ID header key, the generation logic for correlation IDs, and provides middleware for popular web frameworks like net/http, Gin, and Echo.

## Features
- **Middleware Support**: Provides middleware for `net/http`, `Gin`, and `Echo` to handle correlation ID propagation.
- **Configurable Key**: Allows customization of the correlation ID header key.
- **Customizable Generation Function**: Allows customization of the function used to generate correlation IDs.
- **Context-Based ID Retrieval**: Provides a way to retrieve the correlation ID from the request context.

## Installation
To install the library, use the go get command:

```bash
go get github.com/yourusername/correlation
```

## Usage

### Middleware for HTTP (`net/http`)
You can use the middleware with the standard `net/http` package to automatically set and propagate a correlation ID.

```go
package main

import (
    "github.com/khivuksergey/correlation"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    correlationID := r.Header.Get(correlation.Key())
    //or use helper function to retrieve id from context
    correlationID = correlation.Id(r.Context())
    w.Write([]byte(correlationID))
}

func main() {
    http.Handle("/", correlation.Middleware(http.HandlerFunc(handler)))
    http.ListenAndServe(":8080", nil)
}
```

### Middleware for [Gin](https://github.com/gin-gonic/gin)
This library provides a middleware for the `Gin` web framework.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/khivuksergey/correlation"
)

func main() {
    r := gin.Default()
    r.Use(correlation.GinMiddleware)
    r.GET("/", func(c *gin.Context) {
        correlationID := c.GetString(correlation.Key()) 
        //or use helper function to retrieve id from context
        correlationID = correlation.Id(c)
        c.String(http.StatusOK, correlationID)
    })
    r.Run(":8081")
}
```

### Middleware for [Echo](https://github.com/labstack/echo)
You can use the middleware with the `Echo` web framework.

```go
package main

import (
    "github.com/khivuksergey/correlation"
    "github.com/labstack/echo/v4"
)

func main() {
    e := echo.New()
    e.Use(correlation.EchoMiddleware)
    e.GET("/", func(c echo.Context) error {
        correlationID := c.Get(correlation.Key()).(string)
        //or use helper function to retrieve id from context
        correlationID = correlation.Id(c)
        return c.String(http.StatusOK, correlationID)
    })
    e.Start(":8082")
}
```

## Functions
### `SetKey(value string)`
Sets the key used for the correlation ID header. If the key is an empty string, it does not change the key.

```go
correlation.SetKey("Custom-Correlation-Key")
```

### `SetGenerateFunc(f func() string)`
Sets the function used to generate the correlation ID. If nil is passed, the default generation function is used.

```go
correlation.SetGenerateFunc(func() string {
    return "Custom-ID"
})
```

### `Default()`
Resets the correlation ID header key and generation function to their default values.

```go
correlation.Default()
```

### `Key() string`
Returns the current correlation ID header key.

```go
key := correlation.Key()
```

### `Id(c context.Context) string`
Retrieves the correlation ID from the provided context. Returns an empty string if the ID is not found or is not a string.

```go
c := context.WithValue(context.Background(), correlation.Key(), "correlation-id")
id := correlation.Id(c)
```