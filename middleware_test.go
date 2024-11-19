package correlation

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

const (
	testKey = "correlation-id"
	testId  = "mock-correlation-id"
)

var tests = []struct {
	name           string
	existingHeader string
	expectedHeader string
	expectedBody   string
}{
	{"With existing correlation ID", "existing-id", "existing-id", "existing-id"},
	{"Without correlation ID", "", testId, testId},
}

func setup() {
	SetKey(testKey)
	SetGenerateFunc(func() string { return testId })
}

func TestMiddleware(t *testing.T) {
	setup()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		correlationID, _ := r.Context().Value(Key()).(string)
		_, _ = w.Write([]byte(correlationID))
	})

	server := httptest.NewServer(Middleware(handler))
	defer server.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", server.URL, nil)
			if tt.existingHeader != "" {
				req.Header.Set(testKey, tt.existingHeader)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("Request failed: %v", err)
			}
			defer func() {
				_ = resp.Body.Close()
			}()

			if got := resp.Header.Get(testKey); got != tt.expectedHeader {
				t.Errorf("Expected header %v, got %v", tt.expectedHeader, got)
			}

			body := make([]byte, 1024)
			n, _ := resp.Body.Read(body)
			if got := string(body[:n]); got != tt.expectedBody {
				t.Errorf("Expected body %v, got %v", tt.expectedBody, got)
			}
		})
	}

	Default()
}

func TestEchoMiddleware(t *testing.T) {
	setup()

	e := echo.New()

	handler := func(c echo.Context) error {
		correlationID := c.Get(testKey).(string)
		return c.String(http.StatusOK, correlationID)
	}

	e.GET("/", handler, EchoMiddleware)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()

			if tt.existingHeader != "" {
				req.Header.Set(testKey, tt.existingHeader)
			}

			e.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, tt.expectedHeader, rec.Header().Get(testKey))
			assert.Equal(t, tt.expectedBody, rec.Body.String())
		})
	}

	Default()
}

func TestGinMiddleware(t *testing.T) {
	setup()

	gin.SetMode(gin.TestMode)
	r := gin.New()

	r.Use(GinMiddleware)

	r.GET("/", func(c *gin.Context) {
		correlationID := c.GetString(key)
		c.String(http.StatusOK, correlationID)
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.existingHeader != "" {
				req.Header.Set(key, tt.existingHeader)
			}
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, tt.expectedHeader, rec.Header().Get(testKey))
			assert.Equal(t, tt.expectedBody, rec.Body.String())
		})
	}

	Default()
}
