package main

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v3"
)

func TestListActionsFormatting(t *testing.T) {
	// Mock JSON response
	jsonResponse := `[{"kind":"test","data":"some data"}]`
	expectedOutput := `[
  {
    "data": "some data",
    "kind": "test"
  }
]
`

	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonResponse))
	}

	fakeServer := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer fakeServer.Close()

	// Mock getClient to return client pointing to fake server
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return fakeServer.Client()
	}

	// Capture output
	var buf bytes.Buffer
	originalOutputWriter := outputWriter
	outputWriter = &buf
	defer func() { outputWriter = originalOutputWriter }()

	cmd := createCLI()
	params := []string{"dmh-cli", "action", "list", "--server", fakeServer.URL}

	err := cmd.Run(context.Background(), params)
	require.Nil(t, err)

	// Verify output is pretty printed
	// We normalize line endings just in case, though json.MarshalIndent uses \n
	output := buf.String()
	// Replace \r\n with \n for Windows compatibility if running there, though sandbox is likely Linux
	output = strings.ReplaceAll(output, "\r\n", "\n")

	require.Equal(t, expectedOutput, output)
}

func TestListActionsFormattingFallback(t *testing.T) {
	// Mock non-JSON response
	rawResponse := `Not JSON`

	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(rawResponse))
	}

	fakeServer := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer fakeServer.Close()

	// Mock getClient
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return fakeServer.Client()
	}

	// Capture output
	var buf bytes.Buffer
	originalOutputWriter := outputWriter
	outputWriter = &buf
	defer func() { outputWriter = originalOutputWriter }()

	cmd := createCLI()
	params := []string{"dmh-cli", "action", "list", "--server", fakeServer.URL}

	err := cmd.Run(context.Background(), params)
	require.Nil(t, err)

	output := buf.String()
	require.Equal(t, rawResponse, output)
}
