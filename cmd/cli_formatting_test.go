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
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`[{"kind":"test","data":"test"}]`))
	}

	fakeServer := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer fakeServer.Close()

	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return fakeServer.Client()
	}

	var buf bytes.Buffer
	originalOutputWriter := outputWriter
	outputWriter = &buf
	defer func() { outputWriter = originalOutputWriter }()

	cmd := createCLI()
	params := []string{"dmh-cli", "action", "list", "--server", fakeServer.URL}

	err := cmd.Run(context.Background(), params)
	require.Nil(t, err)

	expected := "[\n  {\n    \"kind\": \"test\",\n    \"data\": \"test\"\n  }\n]\n"
	require.Equal(t, expected, buf.String())
}

func TestUpdateAliveOutput(t *testing.T) {
	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	fakeServer := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer fakeServer.Close()

	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(*cli.Command) *http.Client {
		return fakeServer.Client()
	}

	var buf bytes.Buffer
	originalOutputWriter := outputWriter
	outputWriter = &buf
	defer func() { outputWriter = originalOutputWriter }()

	cmd := createCLI()
	params := []string{"dmh-cli", "alive", "update", "--server", fakeServer.URL}

	err := cmd.Run(context.Background(), params)
	require.Nil(t, err)

	expected := "Last seen updated successfully\n"
	require.Equal(t, expected, buf.String())
}
