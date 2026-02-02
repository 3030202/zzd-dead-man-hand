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

func TestListActionsPrettyPrint(t *testing.T) {
	// Setup mock server
	mockResponse := `[{"id":"1","kind":"test","data":"{}"},{"id":"2","kind":"email","data":"{}"}]`
	expectedOutput := `[
  {
    "id": "1",
    "kind": "test",
    "data": "{}"
  },
  {
    "id": "2",
    "kind": "email",
    "data": "{}"
  }
]`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	// Mock getClient
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return ts.Client()
	}

	// Mock outputWriter
	var buf bytes.Buffer
	originalOutputWriter := outputWriter
	outputWriter = &buf
	defer func() { outputWriter = originalOutputWriter }()

	// Execute command
	cmd := createCLI()
	params := []string{"dmh-cli", "action", "list", "--server", ts.URL}
	err := cmd.Run(context.Background(), params)

	require.Nil(t, err)
	require.Equal(t, expectedOutput, buf.String())
}

func TestUpdateAliveOutput(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Mock getClient
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return ts.Client()
	}

	// Mock outputWriter
	var buf bytes.Buffer
	originalOutputWriter := outputWriter
	outputWriter = &buf
	defer func() { outputWriter = originalOutputWriter }()

	// Execute command
	cmd := createCLI()
	params := []string{"dmh-cli", "alive", "update", "--server", ts.URL}
	err := cmd.Run(context.Background(), params)

	require.Nil(t, err)
	require.Contains(t, buf.String(), "Last seen timestamp updated successfully")
}
