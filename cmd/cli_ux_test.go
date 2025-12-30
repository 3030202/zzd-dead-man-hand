package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrettyPrintJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid JSON",
			input:    `[{"id":"123","name":"test"}]`,
			expected: "[\n  {\n    \"id\": \"123\",\n    \"name\": \"test\"\n  }\n]\n",
		},
		{
			name:     "Invalid JSON",
			input:    `not json`,
			expected: "not json",
		},
		{
			name:     "Empty",
			input:    ``,
			expected: "",
		},
		{
			name:     "Number preservation",
			input:    `{"id": 1234567890123456789}`,
			expected: "{\n  \"id\": 1234567890123456789\n}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := prettyPrintJSON(&buf, []byte(tt.input))
			require.NoError(t, err)
			require.Equal(t, tt.expected, buf.String())
		})
	}
}
