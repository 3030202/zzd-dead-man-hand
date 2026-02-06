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
		responseBody   string
		responseStatus int
		expectedOutput string
	}{
		{
			name:           "valid json",
			responseBody:   `[{"id": "1", "name": "test"}]`,
			responseStatus: http.StatusOK,
			expectedOutput: "[\n  {\n    \"id\": \"1\",\n    \"name\": \"test\"\n  }\n]\n",
		},
		{
			name:           "invalid json",
			responseBody:   "not json",
			responseStatus: http.StatusOK,
			expectedOutput: "not json",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Mock server
			fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.responseStatus)
				w.Write([]byte(test.responseBody))
			}))
			defer fakeServer.Close()

			// Mock getClient
			originalGetClient := getClient
			defer func() { getClient = originalGetClient }()
			getClient = func(*cli.Command) *http.Client {
				return fakeServer.Client()
			}

			// Mock outputWriter
			originalOutputWriter := outputWriter
			defer func() { outputWriter = originalOutputWriter }()
			var buf bytes.Buffer
			outputWriter = &buf

			// Run command
			cmd := createCLI()
			err := cmd.Run(context.Background(), []string{"dmh-cli", "action", "list", "--server", fakeServer.URL})

			// Verify
			require.Nil(t, err)
			require.Equal(t, test.expectedOutput, buf.String())
		})
	}
}

func TestAddActionOutput(t *testing.T) {
	// Mock server
	fakeServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer fakeServer.Close()

	// Mock getClient
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return fakeServer.Client()
	}

	// Mock outputWriter
	originalOutputWriter := outputWriter
	defer func() { outputWriter = originalOutputWriter }()
	var buf bytes.Buffer
	outputWriter = &buf

	// Run command
	cmd := createCLI()
	params := []string{"dmh-cli", "action", "add", "--server", fakeServer.URL, "--data", "{}", "--kind", "test", "--process-after", "10"}
	err := cmd.Run(context.Background(), params)

	// Verify
	require.Nil(t, err)
	require.Equal(t, "Action added successfully\n", buf.String())
}
