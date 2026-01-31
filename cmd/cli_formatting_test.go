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
			serverResponse: `[{"kind":"test","data":"data"}]`,
			expectedOutput: `[
  {
    "kind": "test",
    "data": "data"
  }
]`,
		},
		{
			name:           "Invalid JSON",
			serverResponse: `Not a JSON string`,
			expectedOutput: `Not a JSON string`,
		},
		{
			name:           "Empty Array",
			serverResponse: `[]`,
			expectedOutput: `[]`,
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
			var outputBuffer bytes.Buffer
			originalOutputWriter := outputWriter
			defer func() { outputWriter = originalOutputWriter }()
			outputWriter = &outputBuffer

			// Run command
			cmd := createCLI()
			err := cmd.Run(context.Background(), []string{"dmh-client", "action", "list", "--server", fakeServer.URL})

			require.Nil(t, err)
			require.Equal(t, test.expectedOutput, outputBuffer.String())
		})
	}
}
