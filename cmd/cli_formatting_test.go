package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"os"
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
			serverResponse: `[{"id":"1","name":"test"}]`,
			expectedOutput: "[\n  {\n    \"id\": \"1\",\n    \"name\": \"test\"\n  }\n]",
		},
		{
			name:           "Invalid JSON",
			serverResponse: `not json`,
			expectedOutput: `not json`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mock server
			fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(test.serverResponse))
			}))
			defer fakeServer.Close()

			// Mock getClient
			originalGetClient := getClient
			defer func() { getClient = originalGetClient }()
			getClient = func(*cli.Command) *http.Client {
				return fakeServer.Client()
			}

			// Mock outputWriter
			var buf bytes.Buffer
			outputWriter = &buf
			defer func() { outputWriter = os.Stdout }()

			// Run command
			cmd := createCLI()
			params := []string{"dmh-cli", "action", "list", "--server", fakeServer.URL}
			err := cmd.Run(context.Background(), params)

			require.Nil(t, err)
			require.Equal(t, test.expectedOutput, buf.String())
		})
	}
}
