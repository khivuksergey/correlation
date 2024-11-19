package correlation

import (
	"strings"

	"github.com/google/uuid"
)

const defaultCorrelationIdKey = "X-Correlation-ID"

var defaultGenerateFunc = func() string { return strings.ToUpper(uuid.New().String()) }

var (
	key      = defaultCorrelationIdKey
	generate = defaultGenerateFunc
)

// SetKey sets the correlation ID header key to the provided value if it is non-empty.
// If an empty string is provided, the key remains unchanged. This function allows customization
// of the header key used for correlation IDs across the application.
func SetKey(value string) {
	if len(value) == 0 {
		return
	}
	key = value
}

// SetGenerateFunc sets the function used to generate correlation IDs.
// If the provided function is nil, the generator function remains unchanged.
// This function allows for customization of the correlation ID generation logic.
func SetGenerateFunc(f func() string) {
	if f == nil {
		return
	}
	generate = f
}

// Default resets the correlation ID header key and the correlation ID generator function
// to their default values. The default header key is "X-Correlation-ID" and the default
// generator function generates a new UUID in uppercase format.
func Default() {
	key = defaultCorrelationIdKey
	generate = defaultGenerateFunc
}
