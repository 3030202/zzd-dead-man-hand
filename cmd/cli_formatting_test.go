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
	tests := []struct {
		name           string
		serverResponse string
		expectedOutput string
	}{
		{
			name:           "Valid JSON",
			serverResponse: `[{"id":"1","kind":"test"}]`,
			expectedOutput: "[\n  {\n    \"id\": \"1\",\n    \"kind\": \"test\"\n  }\n]\n",
		},
		{
			name:           "Invalid JSON",
			serverResponse: `Not JSON`,
			expectedOutput: "Not JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.serverResponse))
			}))
			defer server.Close()

			// Mock client to use server
			originalGetClient := getClient
			getClient = func(*cli.Command) *http.Client {
				return server.Client()
			}
			defer func() { getClient = originalGetClient }()

			// Mock outputWriter
			var buf bytes.Buffer
			originalWriter := outputWriter
			outputWriter = &buf
			defer func() { outputWriter = originalWriter }()

			// Run command
			cmd := createCLI()
			err := cmd.Run(context.Background(), []string{"dmh-client", "action", "list", "--server", server.URL})
			require.Nil(t, err)

			// Verify output
			require.Equal(t, tt.expectedOutput, buf.String())
		})
	}
}
