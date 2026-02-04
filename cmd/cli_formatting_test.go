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

func TestListActions_PrettyPrint(t *testing.T) {
	// Mock JSON response
	mockJSON := `[{"id":"1","name":"test"}]`
	expectedOutput := "[\n  {\n    \"id\": \"1\",\n    \"name\": \"test\"\n  }\n]\n"

	// Setup mock server
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockJSON))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	// Mock getClient
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return server.Client()
	}

	// Capture output
	var buf bytes.Buffer
	originalOutputWriter := outputWriter
	outputWriter = &buf
	defer func() { outputWriter = originalOutputWriter }()

	// Execute command
	cmd := createCLI()
	args := []string{"dmh-client", "action", "list", "--server", server.URL}
	err := cmd.Run(context.Background(), args)

	require.NoError(t, err)
	require.Equal(t, expectedOutput, buf.String())
}

func TestListActions_RawOutput(t *testing.T) {
	// Mock non-JSON response
	mockResponse := "Just some plain text"

	// Setup mock server
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}
	server := httptest.NewServer(http.HandlerFunc(handler))
	defer server.Close()

	// Mock getClient
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return server.Client()
	}

	// Capture output
	var buf bytes.Buffer
	originalOutputWriter := outputWriter
	outputWriter = &buf
	defer func() { outputWriter = originalOutputWriter }()

	// Execute command
	cmd := createCLI()
	args := []string{"dmh-client", "action", "list", "--server", server.URL}
	err := cmd.Run(context.Background(), args)

	require.NoError(t, err)
	require.Equal(t, mockResponse, buf.String())
}
