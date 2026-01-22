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
	mockResponse := `[{"id":"1","kind":"email","data":"test"}]`
	expectedOutput := `[
	{
		"data": "test",
		"id": "1",
		"kind": "email"
	}
]`

	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}

	fakeServer := httptest.NewServer(http.HandlerFunc(mockHandler))
	defer fakeServer.Close()

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

	// Verify that the JSON is equal (semantically)
	require.JSONEq(t, expectedOutput, buf.String())

	// Verify that the formatting is pretty-printed (contains newlines/indentation as expected)
	require.Equal(t, expectedOutput, buf.String())
}
