package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestListActionsFormatting(t *testing.T) {
	// Setup mock server
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"foo":"bar"}`))
	}
	server := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer server.Close()

	// Mock getClient
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return server.Client()
	}

	tests := []struct {
		name           string
		isTerminal     bool
		expectedOutput string
	}{
		{
			name:       "Terminal",
			isTerminal: true,
			expectedOutput: `{
  "foo": "bar"
}
`,
		},
		{
			name:           "Non-Terminal",
			isTerminal:     false,
			expectedOutput: `{"foo":"bar"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock isTerminal
			originalIsTerminal := isTerminal
			defer func() { isTerminal = originalIsTerminal }()
			isTerminal = func() bool { return tt.isTerminal }

			// Capture output
			var buf bytes.Buffer
			originalOutputWriter := outputWriter
			defer func() { outputWriter = originalOutputWriter }()
			outputWriter = &buf

			cmd := createCLI()
			err := cmd.Run(context.Background(), []string{"dmh-cli", "action", "list", "--server", server.URL})
			require.Nil(t, err)

			require.Equal(t, tt.expectedOutput, buf.String())
		})
	}
}
