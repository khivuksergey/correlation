package correlation

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	Default()
	assert.Equal(t, defaultCorrelationIdKey, Key())
}

func TestId(t *testing.T) {
	Default()
	tests := []struct {
		name       string
		context    context.Context
		expectedID string
	}{
		{
			name:       "With correlation ID in context",
			context:    context.WithValue(context.Background(), key, "test-correlation-id"),
			expectedID: "test-correlation-id",
		},
		{
			name:       "Without correlation ID in context",
			context:    context.Background(),
			expectedID: "",
		},
		{
			name:       "Non-string value in context",
			context:    context.WithValue(context.Background(), key, 12345),
			expectedID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retrievedID := Id(tt.context)
			assert.Equal(t, tt.expectedID, retrievedID)
		})
	}
}

func TestIdFromEcho(t *testing.T) {
	Default()
	tests := []struct {
		name       string
		setupFunc  func(c echo.Context)
		expectedID string
	}{
		{
			name: "With correlation ID in context",
			setupFunc: func(c echo.Context) {
				c.Set(key, "test-correlation-id")
			},
			expectedID: "test-correlation-id",
		},
		{
			name:       "Without correlation ID in context",
			setupFunc:  func(c echo.Context) {},
			expectedID: "",
		},
		{
			name: "Non-string value in context",
			setupFunc: func(c echo.Context) {
				c.Set(key, 12345)
			},
			expectedID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			tt.setupFunc(c)

			retrievedID := IdFromEcho(c)

			assert.Equal(t, tt.expectedID, retrievedID)
		})
	}
}
