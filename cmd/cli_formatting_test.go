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
	// Original state restoration
	originalGetClient := getClient
	originalOutputWriter := outputWriter
	defer func() {
		getClient = originalGetClient
		outputWriter = originalOutputWriter
	}()

	// Mock server
	jsonResponse := `[{"id":"123","name":"test action"}]`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(jsonResponse))
	}))
	defer ts.Close()

	// Mock client
	getClient = func(cmd *cli.Command) *http.Client {
		return ts.Client()
	}

	// Mock output
	var buf bytes.Buffer
	outputWriter = &buf

	// Run command
	cmd := createCLI()
	err := cmd.Run(context.Background(), []string{"dmh-cli", "action", "list", "--server", ts.URL})
	require.Nil(t, err)

	// Verify output
	expectedOutput := `[
  {
    "id": "123",
    "name": "test action"
  }
]
`
	require.Equal(t, expectedOutput, buf.String())
}

func TestListActionsFallbackFormatting(t *testing.T) {
	// Original state restoration
	originalGetClient := getClient
	originalOutputWriter := outputWriter
	defer func() {
		getClient = originalGetClient
		outputWriter = originalOutputWriter
	}()

	// Mock server returning non-JSON
	rawResponse := `Not JSON`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(rawResponse))
	}))
	defer ts.Close()

	// Mock client
	getClient = func(cmd *cli.Command) *http.Client {
		return ts.Client()
	}

	// Mock output
	var buf bytes.Buffer
	outputWriter = &buf

	// Run command
	cmd := createCLI()
	err := cmd.Run(context.Background(), []string{"dmh-cli", "action", "list", "--server", ts.URL})
	require.Nil(t, err)

	// Verify output matches raw input
	require.Equal(t, rawResponse, buf.String())
}
