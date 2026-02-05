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
	// Original outputWriter
	oldWriter := outputWriter
	defer func() { outputWriter = oldWriter }()

	// Capture output
	var buf bytes.Buffer
	outputWriter = &buf

	// Mock server returning JSON
	mockResponse := `[{"id":"1","name":"test"}]`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	// Mock client
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return ts.Client()
	}

	// Run command
	cmd := createCLI()
	err := cmd.Run(context.Background(), []string{"dmh-client", "action", "list", "--server", ts.URL})
	require.Nil(t, err)

	// Verify output
	expected := "[\n  {\n    \"id\": \"1\",\n    \"name\": \"test\"\n  }\n]\n"
	require.Equal(t, expected, buf.String())
}

func TestListActionsFallback(t *testing.T) {
	// Original outputWriter
	oldWriter := outputWriter
	defer func() { outputWriter = oldWriter }()

	// Capture output
	var buf bytes.Buffer
	outputWriter = &buf

	// Mock server returning non-JSON
	mockResponse := `Not JSON`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	// Mock client
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return ts.Client()
	}

	// Run command
	cmd := createCLI()
	err := cmd.Run(context.Background(), []string{"dmh-client", "action", "list", "--server", ts.URL})
	require.Nil(t, err)

	// Verify output
	require.Equal(t, mockResponse, buf.String())
}

func TestAddActionOutput(t *testing.T) {
	// Original outputWriter
	oldWriter := outputWriter
	defer func() { outputWriter = oldWriter }()

	// Capture output
	var buf bytes.Buffer
	outputWriter = &buf

	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	// Mock client
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return ts.Client()
	}

	// Run command
	cmd := createCLI()
	err := cmd.Run(context.Background(), []string{"dmh-client", "action", "add", "--data", "{}", "--kind", "test", "--process-after", "10", "--server", ts.URL})
	require.Nil(t, err)

	// Verify output
	require.Equal(t, "Action added successfully\n", buf.String())
}
