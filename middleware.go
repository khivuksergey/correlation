package correlation

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
)

// Middleware is an HTTP middleware that checks the request headers for an existing correlation ID.
// If the correlation ID is not present, it generates a new one. The generated or retrieved
// correlation ID is then set in the request context, request and response headers, allowing downstream
// handlers to access and propagate the correlation ID throughout the request chain.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(key)
		if id == "" {
			id = generate()
		}
		r = r.WithContext(context.WithValue(r.Context(), key, id))
		r.Header.Set(key, id)
		w.Header().Set(key, id)
		next.ServeHTTP(w, r)
	})
}

// EchoMiddleware is an Echo framework middleware that checks the request headers for an existing
// correlation ID. If the correlation ID is missing, it generates a new one. The correlation ID
// is then set both in the Echo context and the request and response headers, ensuring that it is available
// for downstream handlers and can be propagated across services.
func EchoMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Request().Header.Get(key)
		if id == "" {
			id = generate()
		}
		c.Set(key, id)
		c.Request().Header.Set(key, id)
		c.Response().Header().Set(key, id)
		return next(c)
	}
}

// GinMiddleware is a Gin framework middleware that checks the incoming request for an existing
// correlation ID. If the correlation ID is not found in the headers, it generates a new one. The
// correlation ID is then stored in the Gin context and added to the response headers. This ensures
// that each request handled by the Gin framework has a valid correlation ID for tracing and logging.
func GinMiddleware(c *gin.Context) {
	id := c.GetHeader(key)
	if id == "" {
		id = generate()
	}
	c.Set(key, id)
	c.Request.Header.Set(key, id)
	c.Writer.Header().Set(key, id)
	c.Next()
}
