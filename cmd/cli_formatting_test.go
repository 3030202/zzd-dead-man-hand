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
	// Restore outputWriter after test
	originalOutputWriter := outputWriter
	defer func() { outputWriter = originalOutputWriter }()

	// Mock outputWriter
	var buf bytes.Buffer
	outputWriter = &buf

	// Mock Server
	mockResponse := `[{"kind":"test","data":"some-data"}]`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer server.Close()

	// Mock getClient
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(cmd *cli.Command) *http.Client {
		return server.Client()
	}

	// Run Command
	cmd := createCLI()
	params := []string{"dmh-cli", "action", "list", "--server", server.URL}
	err := cmd.Run(context.Background(), params)

	require.Nil(t, err)
	// Expect pretty-printed output
	expectedOutput := `[
  {
    "kind": "test",
    "data": "some-data"
  }
]
`
	require.Equal(t, expectedOutput, buf.String())
}
