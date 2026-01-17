package main

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func captureOutput(f func() error) (string, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	originalStdout := os.Stdout
	os.Stdout = w

	err = f()

	w.Close()
	os.Stdout = originalStdout

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}

func TestListActionsFormatting(t *testing.T) {
	tests := []struct {
		name           string
		mockIsTerminal bool
		responseBody   string
		expectedOutput string
	}{
		{
			name:           "Terminal Pretty Print",
			mockIsTerminal: true,
			responseBody:   `[{"id":1,"name":"test"}]`,
			expectedOutput: "[\n  {\n    \"id\": 1,\n    \"name\": \"test\"\n  }\n]\n",
		},
		{
			name:           "Non-Terminal Raw Output",
			mockIsTerminal: false,
			responseBody:   `[{"id":1,"name":"test"}]`,
			expectedOutput: `[{"id":1,"name":"test"}]`,
		},
		{
			name:           "Invalid JSON Terminal",
			mockIsTerminal: true,
			responseBody:   `invalid-json`,
			expectedOutput: `invalid-json`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock Server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.responseBody))
			}))
			defer server.Close()

			// Mock getClient
			originalGetClient := getClient
			getClient = func(cmd *cli.Command) *http.Client {
				return server.Client()
			}
			defer func() { getClient = originalGetClient }()

			// Mock isTerminal
			originalIsTerminal := isTerminal
			isTerminal = func() bool {
				return tt.mockIsTerminal
			}
			defer func() { isTerminal = originalIsTerminal }()

			// Run command and capture output
			cmd := createCLI()

			output, err := captureOutput(func() error {
				return cmd.Run(context.Background(), []string{"dmh-client", "action", "list", "--server", server.URL})
			})

			require.NoError(t, err)
			require.Equal(t, tt.expectedOutput, output)
		})
	}
}
