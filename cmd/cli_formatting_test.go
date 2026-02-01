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
			name:           "Pretty print JSON",
			serverResponse: `[{"id":"1","kind":"test"}]`,
			expectedOutput: "[\n  {\n    \"id\": \"1\",\n    \"kind\": \"test\"\n  }\n]\n",
		},
		{
			name:           "Raw output for non-JSON",
			serverResponse: "Internal Server Error", // Though listActions checks status code, let's assume status OK but bad body for this test path or simply generic non-json text if valid json check fails.
			// Actually listActions tries to Unmarshal. If it fails, it prints raw body.
			// However, json.Unmarshal("Internal Server Error", ...) will fail.
			expectedOutput: "Internal Server Error\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mock Server
			fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(test.serverResponse))
			}))
			defer fakeServer.Close()

			// Mock Client
			originalGetClient := getClient
			defer func() { getClient = originalGetClient }()
			getClient = func(*cli.Command) *http.Client {
				return fakeServer.Client()
			}

			// Mock Output
			var buf bytes.Buffer
			originalOutputWriter := outputWriter
			outputWriter = &buf
			defer func() { outputWriter = originalOutputWriter }()

			cmd := createCLI()
			params := []string{"dmh-cli", "action", "list", "--server", fakeServer.URL}

			err := cmd.Run(context.Background(), params)
			require.Nil(t, err)
			require.Equal(t, test.expectedOutput, buf.String())
		})
	}
}

func TestUpdateAliveFormatting(t *testing.T) {
	// Mock Server
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer fakeServer.Close()

	// Mock Client
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return fakeServer.Client()
	}

	// Mock Output
	var buf bytes.Buffer
	originalOutputWriter := outputWriter
	outputWriter = &buf
	defer func() { outputWriter = originalOutputWriter }()

	cmd := createCLI()
	params := []string{"dmh-cli", "alive", "update", "--server", fakeServer.URL}

	err := cmd.Run(context.Background(), params)
	require.Nil(t, err)
	require.Equal(t, "Alive status updated successfully\n", buf.String())
}
