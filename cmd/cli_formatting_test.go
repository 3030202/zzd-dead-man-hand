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
	// Setup
	originalGetClient := getClient
	originalOutputWriter := outputWriter
	defer func() {
		getClient = originalGetClient
		outputWriter = originalOutputWriter
	}()

	// Mock server returning minified JSON
	mockResponse := `[{"id":1,"name":"test","data":{"foo":"bar"}}]`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Configure client to use mock server
	getClient = func(*cli.Command) *http.Client {
		return server.Client()
	}

	// Capture output
	var buf bytes.Buffer
	outputWriter = &buf

	// Run command
	cmd := createCLI()
	// We pass the server URL to the command
	err := cmd.Run(context.Background(), []string{"dmh-cli", "action", "list", "--server", server.URL})
	require.Nil(t, err)

	// Assert output is pretty printed
	expectedOutput := `[
  {
    "id": 1,
    "name": "test",
    "data": {
      "foo": "bar"
    }
  }
]
`
	require.Equal(t, expectedOutput, buf.String())
}

func TestListActionsInvalidJSON(t *testing.T) {
	// Setup
	originalGetClient := getClient
	originalOutputWriter := outputWriter
	defer func() {
		getClient = originalGetClient
		outputWriter = originalOutputWriter
	}()

	// Mock server returning invalid JSON (e.g. plain text error or just non-JSON)
	mockResponse := `Internal Server Error`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // Assuming 200 OK but bad body for some reason, or just a string response
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Configure client to use mock server
	getClient = func(*cli.Command) *http.Client {
		return server.Client()
	}

	// Capture output
	var buf bytes.Buffer
	outputWriter = &buf

	// Run command
	cmd := createCLI()
	err := cmd.Run(context.Background(), []string{"dmh-cli", "action", "list", "--server", server.URL})
	require.Nil(t, err)

	// Assert output is passed through raw
	require.Equal(t, mockResponse, buf.String())
}
