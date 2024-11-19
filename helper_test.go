package correlation

import (
	"context"
	"testing"

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
