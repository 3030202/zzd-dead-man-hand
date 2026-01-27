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
	// Mock server response with unformatted JSON
	rawJSON := `[{"kind":"test","data":"some data","comment":"comment"}]`
	expectedOutput := "[\n  {\n    \"kind\": \"test\",\n    \"data\": \"some data\",\n    \"comment\": \"comment\"\n  }\n]"

	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(rawJSON))
	}

	fakeServer := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer fakeServer.Close()

	// Mock getClient
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return fakeServer.Client()
	}

	// Mock outputWriter
	var buf bytes.Buffer
	originalOutputWriter := outputWriter
	outputWriter = &buf
	defer func() { outputWriter = originalOutputWriter }()

	// Run command
	cmd := createCLI()
	params := []string{"dmh-cli", "action", "list", "--server", fakeServer.URL}

	err := cmd.Run(context.Background(), params)
	require.Nil(t, err)

	// Verify output
	require.Equal(t, expectedOutput, buf.String())
}

func TestListActionsFallback(t *testing.T) {
	// Mock server response with non-JSON
	rawText := "Internal Server Error"

	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // Even if status is OK, body might be non-JSON
		w.Write([]byte(rawText))
	}

	fakeServer := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer fakeServer.Close()

	// Mock getClient
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return fakeServer.Client()
	}

	// Mock outputWriter
	var buf bytes.Buffer
	originalOutputWriter := outputWriter
	outputWriter = &buf
	defer func() { outputWriter = originalOutputWriter }()

	// Run command
	cmd := createCLI()
	params := []string{"dmh-cli", "action", "list", "--server", fakeServer.URL}

	err := cmd.Run(context.Background(), params)
	require.Nil(t, err)

	// Verify output
	require.Equal(t, rawText, buf.String())
}
