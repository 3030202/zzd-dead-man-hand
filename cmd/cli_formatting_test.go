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
	// Mock server returning JSON
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"id":"1","kind":"test"}]`))
	}))
	defer ts.Close()

	// Mock output writer
	var buf bytes.Buffer
	originalOutputWriter := outputWriter
	outputWriter = &buf
	defer func() { outputWriter = originalOutputWriter }()

	// Mock getClient
	originalGetClient := getClient
	getClient = func(*cli.Command) *http.Client {
		return ts.Client()
	}
	defer func() { getClient = originalGetClient }()

	cmd := createCLI()
	err := cmd.Run(context.Background(), []string{"dmh-cli", "action", "list", "--server", ts.URL})
	require.Nil(t, err)

	expected := "[\n  {\n    \"id\": \"1\",\n    \"kind\": \"test\"\n  }\n]\n"
	require.Equal(t, expected, buf.String())
}

func TestListActionsFallbackFormatting(t *testing.T) {
	// Mock server returning non-JSON
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`not json`))
	}))
	defer ts.Close()

	// Mock output writer
	var buf bytes.Buffer
	originalOutputWriter := outputWriter
	outputWriter = &buf
	defer func() { outputWriter = originalOutputWriter }()

	// Mock getClient
	originalGetClient := getClient
	getClient = func(*cli.Command) *http.Client {
		return ts.Client()
	}
	defer func() { getClient = originalGetClient }()

	cmd := createCLI()
	err := cmd.Run(context.Background(), []string{"dmh-cli", "action", "list", "--server", ts.URL})
	require.Nil(t, err)

	require.Equal(t, "not json", buf.String())
}

func TestUpdateAliveSuccessMessage(t *testing.T) {
	// Mock server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Mock output writer
	var buf bytes.Buffer
	originalOutputWriter := outputWriter
	outputWriter = &buf
	defer func() { outputWriter = originalOutputWriter }()

	// Mock getClient
	originalGetClient := getClient
	getClient = func(*cli.Command) *http.Client {
		return ts.Client()
	}
	defer func() { getClient = originalGetClient }()

	cmd := createCLI()
	err := cmd.Run(context.Background(), []string{"dmh-cli", "alive", "update", "--server", ts.URL})
	require.Nil(t, err)

	require.Equal(t, "Successfully updated last seen status\n", buf.String())
}
