package correlation

import "context"

// Key returns the current correlation ID header key used in the application.
// The key's default value is "X-Correlation-ID", unless customized using the SetKey function.
func Key() string {
	return key
}

// Id retrieves the correlation ID stored in the provided context.
// It looks for the correlation ID using the key defined in the application context.
// If a correlation ID is found and is of type string, it returns the ID, otherwise it returns an empty string.
func Id(c context.Context) string {
	id, ok := c.Value(key).(string)
	if !ok {
		return ""
	}
	return id
}
