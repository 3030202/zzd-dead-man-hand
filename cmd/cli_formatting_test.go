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
	// 1. Mock Server
	mockResponse := `[{"id":"1","data":"test"}]`
	// json.Indent result for the above:
	// [
	//   {
	//     "id": "1",
	//     "data": "test"
	//   }
	// ]
	// Plus Fprintln newline.
	expectedOutput := "[\n  {\n    \"data\": \"test\",\n    \"id\": \"1\"\n  }\n]\n"

	// Note: JSON key order is not guaranteed by map iteration, but `[{"id":"1","data":"test"}]` is a list of objects.
	// When `json.Indent` processes it, it preserves key order if unmarshaled into `interface{}`/`map`?
	// No, `json.Indent` works on the *byte slice* directly, it effectively re-formats it.
	// It doesn't re-sort keys unless it unmarshals.
	// `json.Indent` just adds whitespace. So order should be preserved as in input string.
	// So if input is `id` then `data`, output should be `id` then `data`.

	mockResponse = `[{"id":"1","data":"test"}]`
	expectedOutput = "[\n  {\n    \"id\": \"1\",\n    \"data\": \"test\"\n  }\n]\n"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	// 2. Mock getClient
	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(cmd *cli.Command) *http.Client {
		return ts.Client()
	}

	// 3. Mock outputWriter
	var buf bytes.Buffer
	originalOutputWriter := outputWriter
	defer func() { outputWriter = originalOutputWriter }()
	outputWriter = &buf

	// 4. Run Command
	cmd := createCLI()
	args := []string{"dmh-cli", "action", "list", "--server", ts.URL}
	err := cmd.Run(context.Background(), args)

	// 5. Verify
	require.NoError(t, err)
	require.Equal(t, expectedOutput, buf.String())
}

func TestListActionsFallback(t *testing.T) {
	// Test that invalid JSON is printed as-is
	mockResponse := `not json`
	expectedOutput := "not json\n"

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockResponse))
	}))
	defer ts.Close()

	originalGetClient := getClient
	defer func() { getClient = originalGetClient }()
	getClient = func(cmd *cli.Command) *http.Client {
		return ts.Client()
	}

	var buf bytes.Buffer
	originalOutputWriter := outputWriter
	defer func() { outputWriter = originalOutputWriter }()
	outputWriter = &buf

	cmd := createCLI()
	args := []string{"dmh-cli", "action", "list", "--server", ts.URL}
	err := cmd.Run(context.Background(), args)

	require.NoError(t, err)
	require.Equal(t, expectedOutput, buf.String())
}
