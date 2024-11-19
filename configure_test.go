package correlation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetKey(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedKey string
	}{
		{"Set valid key", "New-Key", "New-Key"},
		{"Set empty key", "", defaultCorrelationIdKey},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key = defaultCorrelationIdKey
			SetKey(tt.input)
			assert.Equal(t, tt.expectedKey, key)
		})
	}
}

func TestSetGenerateFunc(t *testing.T) {
	testDefaultGenerateFunc := func() string { return "default" }
	testCustomGenerateFunc := func() string { return "Custom-Generate" }

	tests := []struct {
		name             string
		inputGenerate    func() string
		expectedGenerate func() string
	}{
		{"Set valid generate function", testCustomGenerateFunc, testCustomGenerateFunc},
		{"Set nil generate function", nil, testDefaultGenerateFunc},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			generate = testDefaultGenerateFunc
			SetGenerateFunc(tt.inputGenerate)
			assert.Equal(t, tt.expectedGenerate(), generate())
		})
	}
}

func TestDefault(t *testing.T) {
	key = "Modified-Key"
	generate = func() string { return "Modified-Generate" }

	Default()

	assert.Equal(t, defaultCorrelationIdKey, key)
	assert.NotEqual(t, "Modified-Generate", generate())
}
